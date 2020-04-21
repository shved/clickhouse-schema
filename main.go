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

	databases := getDatabases(db)

	for _, dbName := range databases {
		fmt.Println("DB:", dbName)
		tables := getTables(db, dbName)
		fmt.Println(tables)
	}

	return nil
}

func getTables(db *sql.DB, dbName string) []string {
	var tables []string
	rows, err := db.Query("SELECT name FROM system.tables WHERE database = ?;", dbName)
	if err != nil {
		log.Fatalf("getting tables for %s: %v", dbName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("getting tables for %s: %v", dbName, err)
		}
		tables = append(tables, name)
	}

	if rows.Err(); err != nil {
		log.Fatalf("getting tables for %s: %v", dbName, err)
	}

	return tables
}

func getDatabases(db *sql.DB) []string {
	var databases []string
	rows, err := db.Query("SHOW DATABASES FORMAT TabSeparated;")
	if err != nil {
		log.Fatalf("getting databases: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatalf("getting databases: %v", err)
		}
		databases = append(databases, name)
	}

	if rows.Err(); err != nil {
		log.Fatalf("getting databases: %v", err)
	}

	return databases
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
