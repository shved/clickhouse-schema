package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"
)

func main() {
	var helpPtr = flag.Bool("help", false, " Print usage")
	var clickhouseUrlPtr = flag.String("url", "", " ClickHouse url with port, user and password if needed (clickhouse://your.host:9000?username=default&password=&x-multi-statement=true)")
	var filePtr = flag.String("file", "schema.sql", " Output file with path")

	flag.Parse()

	if *helpPtr || len(*clickhouseUrlPtr) < 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	err := writeSchema(clickhouseUrlPtr, filePtr)
	if err != nil {
		log.Fatal(err)
	}
}

func writeSchema(url *string, path *string) error {
	db := connectToClickHouse(url)
	defer db.Close()

	var databases []string
	rows, err := db.Query("SHOW DATABASES FORMAT TabSeparated;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		databases = append(databases, name)
	}

	if !rows.NextResultSet() {
		fmt.Println("No databases in given cluster")
		os.Exit(0)
	}

	if rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)

	return nil
}

func connectToClickHouse(url *string) *sql.DB {
	db, err := sql.Open("clickhouse", *url)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec("SELECT 1"); err != nil {
		log.Fatal("Fail to ping clickhouse. ", err)
	}

	return db
}
