package db

import (
	"fmt"
	"time"

	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"
)

func NewCHConn(url *string) (*sql.DB, error) {
	db, err := sql.Open("clickhouse", *url)
	if err != nil {
		return nil, fmt.Errorf("getting clickhouse connection: %v", err)
	}

	retry(3, time.Second, func() error {
		if _, err := db.Exec("SELECT 1"); err != nil {
			return fmt.Errorf("trying to ping clickhouse: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
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
