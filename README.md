# go-discover-sql-caller

what's this

- list sql queries in go source code
    - function name
    - line number
    - sql query checksum
    - sql query
- track call hierarchy in go source code

restrictions

- can handle only 1 file

# Usage

When use with one file, use `-file` option.
When use with package, use `-dir` option.

```sh

```sh
  -dir string
        directory path
  -file string
        file path
  -format string
        output format. text or tsv or json (default "text")
  -notrace
        do not trace caller function. only show sql query definitions in the file
```

format: text(default)

```
$ go-discover-sql-caller -file /path/to/file.go
Location                               Checksum                          SQL
dispenseID:106                         8FAC9DB94464380B4EAB33D717A942BE  REPLACE INTO idgen (stub) VALUES (?);
parseViewer:294,retrieveTenantRow:339  C55C59B417205E38BD8968D58C1D3059  SELECT * FROM tenant WHERE id = ?;
...
```

format: text(default), notrace

```
$ go-discover-sql-caller -file /path/to/file.go -notrace
Location                               Checksum                          SQL
dispenseID:106         8FAC9DB94464380B4EAB33D717A942BE  REPLACE INTO idgen (stub) VALUES (?);
retrieveTenantRow:339  C55C59B417205E38BD8968D58C1D3059  SELECT * FROM tenant WHERE id = ?;
...
```

format: tsv

```
$ go-discover-sql-caller -file /path/to/file.go -format tsv
Location\tChecksum\tSQL
dispenseID:106\t8FAC9DB94464380B4EAB33D717A942BE\tREPLACE INTO idgen (stub) VALUES (?);
parseViewer:294,retrieveTenantRow:339\tC55C59B417205E38BD8968D58C1D3059\tSELECT * FROM tenant WHERE id = ?;
...
```

format: json 

```
$ go-discover-sql-caller -file /path/to/file.go -format json
[{"FileName":"main.go","LineNum":106,"ColNum":40,"FuncName":"dispenseID","Checksum":"8FAC9DB94464380B4EAB33D717A942BE","SQL":"REPLACE INTO idgen (stub) VALUES (?);"},{"FileName":"main.go","LineNum":294,"ColNum":17,"FuncName":"parseViewer","Caller":{"FileName":"main.go","LineNum":339,"ColNum":3,"FuncName":"retrieveTenantRow","Checksum":"C55C59B417205E38BD8968D58C1D3059","SQL":"SELECT * FROM tenant WHERE id = ?"}},...]
```

format: jsonl

```
$ go-discover-sql-caller -file /path/to/file.go -format json | jq '.[]' -c
{"FileName":"main.go","LineNum":106,"ColNum":40,"FuncName":"dispenseID","Checksum":"8FAC9DB94464380B4EAB33D717A942BE","SQL":"REPLACE INTO idgen (stub) VALUES (?);"}
{"FileName":"main.go","LineNum":294,"ColNum":17,"FuncName":"parseViewer","Caller":{"FileName":"main.go","LineNum":339,"ColNum":3,"FuncName":"retrieveTenantRow","Checksum":"C55C59B417205E38BD8968D58C1D3059","SQL":"SELECT * FROM tenant WHERE id = ?"}}
...
```

# Install

use `go install` , or download binary from [release page](https://github.com/netmarkjp/go-discover-sql-caller/releases/latest).

```sh
go install github.com/netmarkjp/go-discover-sql-caller@latest
```
