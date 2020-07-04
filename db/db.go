package db

import (
	"fmt"
	"log"
	"time"

	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"
)

func NewCHConn(url *string) *sql.DB {
	db, err := sql.Open("clickhouse", *url)
	if err != nil {
		log.Fatalf("getting clickhouse connection: %v", err)
	}

	retry(3, time.Second, func() error {
		if _, err := db.Exec("SELECT 1"); err != nil {
			return fmt.Errorf("trying to ping clickhouse: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// retry clickhouse connections, taken from https://upgear.io/blog/simple-golang-retry-function/ and slightly modified
func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}
