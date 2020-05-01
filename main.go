package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"
)

func main() {
	var helpPtr = flag.Bool("help", false, " Print usage")
	var clickhouseUrlPtr = flag.String("url", "", " ClickHouse url with port, user and password if needed (clickhouse://your.host:9000?username=default&password=&x-multi-statement=true)")
	var filePtr = flag.String("file", "schema.sql", " Output file with path")
	var specifiedDB = flag.String("database", "", " Specify schema to be dumped. Otherwise dump all the DBs")

	flag.Parse()

	if *helpPtr || len(*clickhouseUrlPtr) < 1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	writeSchema(clickhouseUrlPtr, filePtr, specifiedDB)
}

func writeSchema(url *string, path *string, specifiedDB *string) {
	db := connectToClickHouse(url)
	defer db.Close()

	fd, err := os.OpenFile(*path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer fd.Close()
	if err != nil {
		log.Fatalf("opening file: %v", err)
	}

	var databases []string
	if *specifiedDB == "" {
		databases = getDatabases(db)
	} else {
		databases = []string{*specifiedDB}
	}

	for _, dbName := range databases {
		if dbName == "system" {
			continue // skip system database
		}
		dbCreateStmt := dbCreateStmt(db, dbName)
		fd.Write([]byte(dbCreateStmt + "\n\n"))
		tables := getTables(db, dbName)
		for _, tableName := range tables {
			tableCreateStmt := prettify(tableCreateStmt(db, dbName, tableName))
			fd.Write([]byte(tableCreateStmt + "\n\n"))
		}
	}
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

func dbCreateStmt(db *sql.DB, dbName string) string {
	var createStmt string
	queryStmt := fmt.Sprintf("SHOW CREATE DATABASE %s FORMAT TabSeparated;", dbName)
	err := db.QueryRow(queryStmt).Scan(&createStmt)
	if err != nil {
		log.Fatalf("getting database %s statement: %v", dbName, err)
	}

	return createStmt
}

func tableCreateStmt(db *sql.DB, dbName string, tableName string) string {
	var createStmt string
	queryStmt := fmt.Sprintf("SHOW CREATE TABLE %s.%s FORMAT PrettySpaceNoEscapes;", dbName, tableName)
	err := db.QueryRow(queryStmt).Scan(&createStmt)
	if err != nil {
		log.Fatalf("getting table %s.%s statement: %v", dbName, tableName, err)
	}

	return createStmt
}

func prettify(s string) string {
	var b strings.Builder
	fields := strings.Fields(s)
	var parenthesisStack uint8
	var afterOrderBy bool

	for _, field := range fields {
		switch {

		case strings.HasPrefix(field, "("):
			if parenthesisStack == 0 && !afterOrderBy {
				b.WriteString("\n\t")
				b.WriteString("(")
				b.WriteString("\n\t\t")
				b.WriteString(field[1:])
				b.WriteString(" ")
				parenthesisStack = parenthesisStack + 1
			} else {
				writeDefault(&b, field)
			}

		case strings.HasSuffix(field, ")"):
			if parenthesisStack == 1 {
				b.WriteString(field[:len(field)-1])
				b.WriteString("\n\t)")
				parenthesisStack = parenthesisStack - 1
			} else {
				writeDefault(&b, field)
			}

		case field == "ORDER" || field == "ENGINE" || field == "PARTITION" || field == "SETTINGS" || field == "FROM" || field == "SELECT" || field == "GROUP":
			b.WriteString("\n")
			b.WriteString(field)
			b.WriteString(" ")
			if field == "ORDER" {
				afterOrderBy = true
			}

		case strings.HasSuffix(field, ","):
			if !afterOrderBy {
				b.WriteString(field)
				b.WriteString("\n\t\t")
			} else {
				writeDefault(&b, field)
			}

		default:
			writeDefault(&b, field)
		}
	}

	return b.String()
}

func writeDefault(b *strings.Builder, field string) {
	b.WriteString(field)
	b.WriteString(" ")
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
