package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shved/clickhouse-schema/db"
	"github.com/shved/clickhouse-schema/schema"
)

func main() {
	var help = flag.Bool("help", false, " Print usage")
	var clickhouseUrl = flag.String("url", "", " ClickHouse url with port, user and password if needed (clickhouse://your.host:9000?username=default&password=&x-multi-statement=true)")
	var file = flag.String("file", "schema.sql", " Output file with path")
	var specifiedDB = flag.String("database", "", " Specify schema to be dumped. Otherwise dump all the DBs")

	flag.Parse()

	if *help || len(*clickhouseUrl) < 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	conn := db.NewCHConn(clickhouseUrl)
	defer conn.Close()

	opts := schema.Options{
		DB:          conn,
		Path:        *file,
		SpecifiedDB: *specifiedDB,
	}

	schema.Write(&opts)

	fmt.Printf("Schema successfully saved to %s\n", *file)
}
