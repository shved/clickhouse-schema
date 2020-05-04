package main

import (
	"flag"
	"log"
	"os"

	"github.com/shved/clickhouse-schema/schema"

	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"
)

func main() {
	var help = flag.Bool("help", false, " Print usage")
	var clickhouseUrl = flag.String("url", "", " ClickHouse url with port, user and password if needed (clickhouse://your.host:9000?username=default&password=&x-multi-statement=true)")
	var file = flag.String("file", "schema.sql", " Output file with path")
	var specifiedDB = flag.String("database", "", " Specify schema to be dumped. Otherwise dump all the DBs")
	var raw = flag.Bool("raw", false, " Skip pretty sql formatting")

	flag.Parse()

	if *help || len(*clickhouseUrl) < 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	db := connectToClickHouse(clickhouseUrl)
	defer db.Close()

	opts := schema.Options{
		DB:          db,
		Path:        file,
		SpecifiedDB: specifiedDB,
		Raw:         raw,
	}

	schema.Write(&opts)
}

func connectToClickHouse(url *string) *sql.DB {
	db, err := sql.Open("clickhouse", *url)
	if err != nil {
		log.Fatalf("getting clickhouse connection: %v", err)
	}

	if _, err := db.Exec("SELECT 1"); err != nil {
		log.Fatalf("trying to ping clickhouse: %v", err)
	}

	return db
}
