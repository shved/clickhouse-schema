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
```
or in your code
```
import "github.com/shved/clickhouse-schema/schema"
...
options = schema.Options{
		DB:          db,          // *sql.DB
		Path:        file,        // string
		SpecifiedDB: specifiedDB, // string
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
- [x] add MIT license
- [x] add couple retries while connecting to a database
- [x] separate prettifier package
- [x] remove prettifier since clickhouse client pretty format suddenly works :(
- [x] remove every log.Fatal
- [ ] set cli default output to stdout
- [ ] throw an error if the database name not found (or it is `system` database)
- [ ] refactor with just enough architecture (https://blog.carlmjohnson.net/post/2020/go-cli-how-to-and-advice/)
