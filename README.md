# go-discover-sql-caller

what's this

- list sql queries in go source code
- track call hierarchy in go source code

restrictions

- can handle only 1 file

# Usage

tsv(default)

```
$ go run main.go -file /path/to/file.go
Location\tSQL
dispenseID:106\tREPLACE INTO idgen (stub) VALUES (?);
parseViewer:294,retrieveTenantRow:339\tSELECT * FROM tenant WHERE id = ?;
...
```

json

```
$ go run main.go -file /path/to/file.go -format json
[{"FileName":"main.go","LineNum":106,"ColNum":40,"FuncName":"dispenseID","SQL":"REPLACE INTO idgen (stub) VALUES (?);"},{"FileName":"main.go","LineNum":294,"ColNum":17,"FuncName":"parseViewer","Caller":{"FileName":"main.go","LineNum":339,"ColNum":3,"FuncName":"retrieveTenantRow","SQL":"SELECT * FROM tenant WHERE id = ?"}},...]
```

jsonl

```
$ go run main.go -f /path/to/file.go -format json | jq '.[]' -c 
{"FileName":"main.go","LineNum":106,"ColNum":40,"FuncName":"dispenseID","SQL":"REPLACE INTO idgen (stub) VALUES (?);"}
{"FileName":"main.go","LineNum":294,"ColNum":17,"FuncName":"parseViewer","Caller":{"FileName":"main.go","LineNum":339,"ColNum":3,"FuncName":"retrieveTenantRow","SQL":"SELECT * FROM tenant WHERE id = ?"}}
```
