ClickHouse schema dumper.  

```
$ clickhouse-schema -help

$ clickhouse-schema \
  --url "clickhouse://your.host:9000?username=default&password=querty123&x-multi-statement=true" \
  --file "/path/to/schema.sql"

$ clickhouse-schema \
  --url "clickhouse://your.host:9000?username=default&password=querty123&x-multi-statement=true" \
  --database "testdb"
```

TODO:
 - [ ] add success messages
 - [x] option to choose only one database
 - [ ] throw an error if the database not found
 - [ ] make it work as a library along with executable
 - [ ] add couple retries while connecting to database
 - [x] research output formats
 - [x] prettier table create statements
 - [ ] make prettifier optional
 - [ ] integration test with docker
