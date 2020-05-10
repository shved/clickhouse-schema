package schema

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/shved/clickhouse-schema/db"
)

func TestWrite(t *testing.T) {
	testCH := "tcp://127.0.0.1:9000?debug=true"
	testOutput := "test/test_schema.sql"
	testSchema := "test/fixture_schema.sql"
	testDB := "testdb"

	conn := db.NewCHConn(&testCH)
	defer conn.Close()

	_, err := conn.Exec("DROP DATABASE IF EXISTS testdb")
	if err != nil {
		log.Fatal("cleaning test database: ", err)
	}

	schema, err := ioutil.ReadFile(testSchema)
	if err != nil {
		log.Fatal("reading fixture file: ", err)
	}

	stmts := strings.Split(string(schema), "\n\n")
	for _, stmt := range stmts {
		if stmt == "" {
			continue // skip
		}
		_, err = conn.Exec(string(stmt))
		if err != nil {
			log.Fatal("filling the test db: ", err)
		}
	}

	options := Options{
		DB:          conn,
		Path:        testOutput,
		SpecifiedDB: testDB,
		Raw:         true,
	}

	Write(&options)

	expected := schema
	given, err := ioutil.ReadFile(testOutput)
	if err != nil {
		log.Fatal("reading output schema file: ", err)
	}

	if string(expected) != string(given) {
		t.Fatalf("Input and output schema doesn't match.\n\tExpected:\n\t%v\n\tGot:\n\t%v\n", string(expected), string(given))
	}
}
