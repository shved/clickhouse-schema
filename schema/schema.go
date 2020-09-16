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
	var fd *os.File
	var err error
	if len(opts.Path) > 0 {
		fd, err = os.OpenFile(opts.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		defer fd.Close()
		if err != nil {
			return fmt.Errorf("opening file: %v", err)
		}
	} else {
		fd = os.Stdout
	}

	allDatabases, err := getDatabases(opts.DB)
	if err != nil {
		return fmt.Errorf("getting databases: %v", err)
	}

	var databases []string
	switch opts.SpecifiedDB {
	case "":
		databases = allDatabases
	case "system":
		return fmt.Errorf("'%s' is a special internal ClickHouse database and can't be specified", opts.SpecifiedDB)
	default:
		if includes(allDatabases, opts.SpecifiedDB) {
			databases = []string{opts.SpecifiedDB}
		} else {
			return fmt.Errorf("specified database '%s' doesnt exist", opts.SpecifiedDB)
		}
	}

	for _, dbName := range databases {
		if dbName == "system" {
			continue
		}
		dbCreateStmt, err := dbCreateStmt(opts.DB, dbName)
		if err != nil {
			return err
		}
		_, err = fd.Write([]byte(dbCreateStmt + "\n\n"))
		if err != nil {
			return fmt.Errorf("writing database '%s' create statement: %v", dbName, err)
		}
		tables, err := getTables(opts.DB, dbName)
		if err != nil {
			return err
		}
		for _, tableName := range tables {
			tableCreateStmt, err := fetchTableCreateStmt(opts.DB, dbName, tableName)
			if err != nil {
				return err
			}
			_, err = fd.Write([]byte(tableCreateStmt + "\n\n"))
			if err != nil {
				return fmt.Errorf("writing table '%s' create statement: %v", tableName, err)
			}
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
		return []string{}, fmt.Errorf("getting tables for '%s': %v", dbName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return []string{}, fmt.Errorf("getting tables for '%s': %v", dbName, err)
		}
		if !(len(name) > 6 && name[:7] == ".inner.") {
			tables = append(tables, name)
		}
	}

	if rows.Err(); err != nil {
		return []string{}, fmt.Errorf("getting tables for '%s': %v", dbName, err)
	}

	return tables, nil
}

func dbCreateStmt(db *sql.DB, dbName string) (string, error) {
	var createStmt string
	queryStmt := fmt.Sprintf("SHOW CREATE DATABASE %s FORMAT PrettySpaceNoEscapes;", dbName)
	err := db.QueryRow(queryStmt).Scan(&createStmt)
	if err != nil {
		return "", fmt.Errorf("getting database %s create statement: %v", dbName, err)
	}

	return createStmt, nil
}

func fetchTableCreateStmt(db *sql.DB, dbName string, tableName string) (string, error) {
	var createStmt string
	queryStmt := fmt.Sprintf("SHOW CREATE TABLE %s.%s FORMAT PrettySpaceNoEscapes;", dbName, tableName)
	err := db.QueryRow(queryStmt).Scan(&createStmt)
	if err != nil {
		return "", fmt.Errorf("getting table '%s.%s' statement: %v", dbName, tableName, err)
	}

	return createStmt, nil
}

func includes(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}

	return false
}
