# Overview
ClickHouse schema dumper. It could be helpful if you want to generate pretty schema file to read and understand your database state. Can be used as a cli tool and a go library.  

# Usage
```
$ clickhouse-schema --help

$ clickhouse-schema \
  --url "clickhouse://your.host:9000?username=default&password=querty123&x-multi-statement=true" \
  --file "/path/to/schema.sql"

$ clickhouse-schema \
  --url "clickhouse://your.host:9000?username=default&password=querty123&x-multi-statement=true" \
  --database "testdb" \
  --raw
```
or in your code
```
import "github.com/shved/clickhouse-schema/schema"
...
options = schema.Options{
		DB:          db,          // *sql.DB
		Path:        file,        // string
		SpecifiedDB: specifiedDB, // string
		Raw:         rawFlag,     // bool
}
schema.Write(&options)
```

# Testing
```
docker-compose up
go test ./...
```

---  
TODO:
- [x] option to choose only one database
- [x] make it work as a library along with executable
- [x] prettier table create statements
- [x] research output formats
- [x] make prettifier optional
- [x] pass options struct pointer instead of options list
- [x] add success messages
- [x] integration test with docker
- [ ] separate prettifier package
- [ ] throw an error if the database name not found (or it is `system` database)
- [ ] add couple retries while connecting to database
- [ ] set cli default output to stdout
- [ ] add MIT license
- [ ] refactor with just enough architecture (https://blog.carlmjohnson.net/post/2020/go-cli-how-to-and-advice/)
