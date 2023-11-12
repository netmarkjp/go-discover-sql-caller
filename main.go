package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"slices"
	"text/tabwriter"

	"flag"
	"strings"
	// tiparser "github.com/pingcap/tidb/parser"
	// _ "github.com/pingcap/tidb/parser/test_driver"
)

func main() {
	filepath := flag.String("file", "", "file path")
	format := flag.String("format", "text", "output format. text or tsv or json")
	flag.Parse()

	if *filepath == "" {
		slog.Error("-file is required")
		return
	}

	*format = strings.ToLower(*format)
	if *format != "text" && *format != "tsv" && *format != "json" {
		slog.Error("-format must be text or tsv or json")
		return
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *filepath, nil, 0)
	if err != nil {
		slog.Error("%+v", err)
		return
	}

	// 深さ優先で探索してSQLが含まれる関数を探す
	sqlCallers := []*SQLCaller{}
	walkFuncNames := []string{}
	ast.Inspect(f, func(n ast.Node) bool {

		fd, ok := n.(*ast.FuncDecl)
		if ok {
			walkFuncNames = append(walkFuncNames, fd.Name.Name)
			return true
		}

		bl, ok := n.(*ast.BasicLit)
		if !ok {
			return true
		}
		if bl.Kind != token.STRING {
			return true
		}
		sql := strip(bl.Value)
		if !isSQL(sql) {
			return true
		}
		pos := fset.Position(bl.Pos())

		sqlCallers = append(sqlCallers, &SQLCaller{
			FileName: pos.Filename,
			LineNum:  pos.Line,
			ColNum:   pos.Column,
			FuncName: walkFuncNames[len(walkFuncNames)-1],
			SQL:      sql,
		})

		return true
	})

	// SQLが含まれる関数を呼び出している関数を探す
	found := 0
	for {
		found, sqlCallers = discoverSQLCallers(sqlCallers, f, fset)
		if found == 0 {
			break
		}
	}

	slices.SortFunc(sqlCallers, func(a, b *SQLCaller) int {
		if a.FileName != b.FileName {
			return strings.Compare(a.FileName, b.FileName)
		}
		if a.LineNum != b.LineNum {
			return a.LineNum - b.LineNum
		}
		if a.ColNum != b.ColNum {
			return a.ColNum - b.ColNum
		}
		return 1
	})

	if *format == "json" {
		// json marshal
		b, err := json.Marshal(sqlCallers)
		if err != nil {
			slog.Error("Error: %+v", err)
			return
		}
		fmt.Printf("%s\n", string(b))
	} else if *format == "tsv" {
		// tsv
		fmt.Println("Location\tChecksum\tSQL")
		for _, c := range sqlCallers {
			c := c
			// fmt.Printf("%s\t%s\t%d\t%d\t%s\t%s\n", c.FileName, c.FuncName, c.LineNum, c.ColNum, c.SQL, c.Caller.Describe())
			fmt.Println(c.Describe())
		}
	} else {
		w := tabwriter.NewWriter(os.Stdout, 2, 0, 1, ' ', 0)
		fmt.Fprintln(w, "Location\tChecksum\tSQL")
		for _, c := range sqlCallers {
			c := c
			fmt.Fprintln(w, c.Describe())
		}
		w.Flush()
	}
}

func discoverSQLCallers(sqlCallers []*SQLCaller, f *ast.File, fset *token.FileSet) (int, []*SQLCaller) {
	found := 0
	sqlCallerFuncNames := []string{}
	for _, sqlCaller := range sqlCallers {
		if !slices.Contains(sqlCallerFuncNames, sqlCaller.FuncName) {
			sqlCallerFuncNames = append(sqlCallerFuncNames, sqlCaller.FuncName)
		}
	}

	walkFuncNames := []string{}
	ast.Inspect(f, func(n ast.Node) bool {
		fd, ok := n.(*ast.FuncDecl)
		if ok {
			walkFuncNames = append(walkFuncNames, fd.Name.Name)
			return true
		}

		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		ident, ok := ce.Fun.(*ast.Ident)
		if !ok {
			return true
		}
		if ident.Obj == nil {
			return true
		}

		callerFuncName := ident.Obj.Name

		if index := slices.Index(sqlCallerFuncNames, callerFuncName); index != -1 {
			foundFuncName := sqlCallerFuncNames[index]

			newCaller := &SQLCaller{}
			pos := fset.Position(ident.NamePos)
			for _, c := range sqlCallers {
				c := c
				if c.FuncName == foundFuncName {
					newCaller = &SQLCaller{
						FileName: pos.Filename,
						LineNum:  pos.Line,
						ColNum:   pos.Column,
						FuncName: walkFuncNames[len(walkFuncNames)-1],
						Caller: &SQLCaller{
							FileName: c.FileName,
							LineNum:  c.LineNum,
							ColNum:   c.ColNum,
							FuncName: c.FuncName,
							SQL:      c.SQL,
							Caller:   c.Caller,
						},
					}
				}
			}

			// if not already exists
			callerFound := false
			for _, c := range sqlCallers {
				c := c
				if SQLCallerEquals(c, newCaller) {
					callerFound = true
					break
				}
			}

			if !callerFound {
				sqlCallers = append(sqlCallers, newCaller)
				found = found + 1
			}
		}

		return true
	})
	return found, sqlCallers
}

func strip(txt string) string {
	txt = strings.TrimSpace(txt)
	txt = strings.Trim(txt, "`")
	txt = strings.Trim(txt, "\"")
	txt = strings.Trim(txt, "'")
	return txt
}

func isSQL(txt string) bool {
	// // tried to use "github.com/pingcap/tidb/parser" but it's too loose
	// p := tiparser.New()
	// _, _, err := p.Parse(txt, "", "")
	// return err == nil

	txt = strings.ToLower(txt)
	if strings.HasPrefix(txt, "select ") {
		return true
	}
	if strings.HasPrefix(txt, "insert ") {
		return true
	}
	if strings.HasPrefix(txt, "update ") {
		return true
	}
	if strings.HasPrefix(txt, "delete ") {
		return true
	}
	if strings.HasPrefix(txt, "replace ") {
		return true
	}
	if strings.HasPrefix(txt, "alter ") {
		return true
	}
	if strings.HasPrefix(txt, "create ") {
		return true
	}
	if strings.HasPrefix(txt, "drop ") {
		return true
	}
	if strings.HasPrefix(txt, "truncate ") {
		return true
	}
	if strings.HasPrefix(txt, "grant ") {
		return true
	}
	if strings.HasPrefix(txt, "revoke ") {
		return true
	}
	if strings.HasPrefix(txt, "begin") {
		return true
	}
	if strings.HasPrefix(txt, "commit") {
		return true
	}
	if strings.HasPrefix(txt, "rollback") {
		return true
	}
	return false
}

// information about the location of the SQL statement
type SQLCaller struct {
	FileName string     `json:"FileName"`
	LineNum  int        `json:"LineNum"`
	ColNum   int        `json:"ColNum"`
	FuncName string     `json:"FuncName"`
	SQL      string     `json:"SQL,omitempty"`
	Checksum string     `json:"Checksum,omitempty"` // empty in usual. dynamically calculated when marshaling/describing
	Caller   *SQLCaller `json:"Caller,omitempty"`
}

func (c *SQLCaller) Describe() string {
	if c.Caller == nil {
		return fmt.Sprintf("%s:%d\t%s\t%s", c.FuncName, c.LineNum, c.SQLChecksum(), c.SQL)
	}
	return fmt.Sprintf("%s:%d,%s", c.FuncName, c.LineNum, c.Caller.Describe())
}

func (c *SQLCaller) MarshalJSON() ([]byte, error) {
	clone := *c
	clone.Checksum = clone.SQLChecksum()

	// avoid infinite loop
	type Alias SQLCaller
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&clone),
	})
}

func (c *SQLCaller) SQLChecksum() string {
	md5sum := md5.Sum([]byte(c.SQL))
	return strings.ToUpper(hex.EncodeToString(md5sum[:]))
}

func SQLCallerEquals(a, b *SQLCaller) bool {
	if a.FileName != b.FileName {
		return false
	}
	if a.LineNum != b.LineNum {
		return false
	}
	if a.ColNum != b.ColNum {
		return false
	}
	if a.FuncName != b.FuncName {
		return false
	}
	if a.SQL != b.SQL {
		return false
	}
	if (a.Caller == nil && b.Caller != nil) || (a.Caller != nil && b.Caller == nil) {
		return false
	}

	if a.Caller == nil && b.Caller == nil {
		return true
	}

	if !SQLCallerEquals(a.Caller, b.Caller) {
		return false
	}

	return true
}
