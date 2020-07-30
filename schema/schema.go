package schema

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/shved/clickhouse-schema/prettifier"
)

type Options struct {
	DB          *sql.DB
	Path        string
	SpecifiedDB string
	Raw         bool
}

func Write(opts *Options) {
	fd, err := os.OpenFile(opts.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer fd.Close()
	if err != nil {
		log.Fatalf("opening file: %v", err)
	}

	var databases []string
	if opts.SpecifiedDB == "" {
		databases = getDatabases(opts.DB)
	} else {
		databases = []string{opts.SpecifiedDB}
	}

	for _, dbName := range databases {
		if dbName == "system" {
			continue // skip system database
		}
		dbCreateStmt := dbCreateStmt(opts.DB, dbName)
		fd.Write([]byte(dbCreateStmt + "\n\n"))
		tables := getTables(opts.DB, dbName)
		for _, tableName := range tables {
			tableCreateStmt := fetchTableCreateStmt(opts.DB, dbName, tableName)
			if !opts.Raw {
				tableCreateStmt = prettifier.Prettify(tableCreateStmt)
			}
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

func fetchTableCreateStmt(db *sql.DB, dbName string, tableName string) string {
	var createStmt string
	queryStmt := fmt.Sprintf("SHOW CREATE TABLE %s.%s FORMAT PrettySpaceNoEscapes;", dbName, tableName)
	err := db.QueryRow(queryStmt).Scan(&createStmt)
	if err != nil {
		log.Fatalf("getting table %s.%s statement: %v", dbName, tableName, err)
	}

	return createStmt
}
