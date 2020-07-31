package schema

import (
	"database/sql"
	"fmt"
	"os"
)

type Options struct {
	DB          *sql.DB
	Path        string
	SpecifiedDB string
}

func Write(opts *Options) error {
	fd, err := os.OpenFile(opts.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer fd.Close()
	if err != nil {
		return fmt.Errorf("opening file: %v", err)
	}

	var databases []string
	if opts.SpecifiedDB == "" {
		databases, err = getDatabases(opts.DB)
		if err != nil {
			return err
		}
	} else {
		databases = []string{opts.SpecifiedDB}
	}

	for _, dbName := range databases {
		if dbName == "system" {
			continue // skip system database
		}
		dbCreateStmt, err := dbCreateStmt(opts.DB, dbName)
		if err != nil {
			return err
		}
		fd.Write([]byte(dbCreateStmt + "\n\n"))
		tables, err := getTables(opts.DB, dbName)
		if err != nil {
			return err
		}
		for _, tableName := range tables {
			tableCreateStmt, err := fetchTableCreateStmt(opts.DB, dbName, tableName)
			if err != nil {
				return err
			}
			fd.Write([]byte(tableCreateStmt + "\n\n"))
		}
	}

	return nil
}

func getDatabases(db *sql.DB) ([]string, error) {
	var databases []string
	rows, err := db.Query("SHOW DATABASES FORMAT TabSeparated;")
	defer rows.Close()
	if err != nil {
		return []string{}, fmt.Errorf("getting databases: %v", err)
	}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return []string{}, fmt.Errorf("getting databases: %v", err)
		}
		databases = append(databases, name)
	}

	if rows.Err(); err != nil {
		return []string{}, fmt.Errorf("getting databases: %v", err)
	}

	return databases, nil
}

func getTables(db *sql.DB, dbName string) ([]string, error) {
	var tables []string
	rows, err := db.Query("SELECT name FROM system.tables WHERE database = ?;", dbName)
	if err != nil {
		return []string{}, fmt.Errorf("getting tables for %s: %v", dbName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return []string{}, fmt.Errorf("getting tables for %s: %v", dbName, err)
		}
		tables = append(tables, name)
	}

	if rows.Err(); err != nil {
		return []string{}, fmt.Errorf("getting tables for %s: %v", dbName, err)
	}

	return tables, nil
}

func dbCreateStmt(db *sql.DB, dbName string) (string, error) {
	var createStmt string
	queryStmt := fmt.Sprintf("SHOW CREATE DATABASE %s FORMAT PrettySpaceNoEscapes;", dbName)
	err := db.QueryRow(queryStmt).Scan(&createStmt)
	if err != nil {
		return "", fmt.Errorf("getting database %s statement: %v", dbName, err)
	}

	return createStmt, nil
}

func fetchTableCreateStmt(db *sql.DB, dbName string, tableName string) (string, error) {
	var createStmt string
	queryStmt := fmt.Sprintf("SHOW CREATE TABLE %s.%s FORMAT PrettySpaceNoEscapes;", dbName, tableName)
	err := db.QueryRow(queryStmt).Scan(&createStmt)
	if err != nil {
		return "", fmt.Errorf("getting table %s.%s statement: %v", dbName, tableName, err)
	}

	return createStmt, nil
}
