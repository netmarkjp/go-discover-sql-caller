# go-discover-sql-caller

what's this

- list sql queries in go source code
- track call hierarchy in go source code

restrictions

- can handle only 1 file

# Usage

```sh
  -checksum
        output checksum of SQL (default:false)
  -file string
        file path
  -format string
        output format. tsv or json (default:tsv) (default "tsv")
```

tsv(default)

```
$ ./go-discover-sql-caller -file /path/to/file.go
Location\tSQL
dispenseID:106\tREPLACE INTO idgen (stub) VALUES (?);
parseViewer:294,retrieveTenantRow:339\tSELECT * FROM tenant WHERE id = ?;
...
```

tsv(default),with checksum

```
$ ./go-discover-sql-caller -file /path/to/file.go -checksum
Location\tSQL
dispenseID:106\tREPLACE INTO idgen (stub) VALUES (?);
parseViewer:294,retrieveTenantRow:339\tC55C59B417205E38BD8968D58C1D3059\tSELECT * FROM tenant WHERE id = ?;
...
```

json

```
$ ./go-discover-sql-caller -file /path/to/file.go -format json
[{"FileName":"main.go","LineNum":106,"ColNum":40,"FuncName":"dispenseID","SQL":"REPLACE INTO idgen (stub) VALUES (?);"},{"FileName":"main.go","LineNum":294,"ColNum":17,"FuncName":"parseViewer","Caller":{"FileName":"main.go","LineNum":339,"ColNum":3,"FuncName":"retrieveTenantRow","SQL":"SELECT * FROM tenant WHERE id = ?"}},...]
```

jsonl

```
$ ./go-discover-sql-caller -file /path/to/file.go -format json | jq '.[]' -c 
{"FileName":"main.go","LineNum":106,"ColNum":40,"FuncName":"dispenseID","SQL":"REPLACE INTO idgen (stub) VALUES (?);"}
{"FileName":"main.go","LineNum":294,"ColNum":17,"FuncName":"parseViewer","Caller":{"FileName":"main.go","LineNum":339,"ColNum":3,"FuncName":"retrieveTenantRow","SQL":"SELECT * FROM tenant WHERE id = ?"}}
```
