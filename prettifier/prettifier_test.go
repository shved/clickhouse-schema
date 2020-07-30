package prettifier

import (
	"io/ioutil"
	"testing"
)

func TestPrettify(t *testing.T) {
	inputCreateTable := "CREATE TABLE testdb.watch_progress (`customer_id` String, `resource_id` String, `resource_type" +
		"` String, `timestamp` Int32, `duration` Int32, `event_time` DateTime, `session_id` String, `action` String) ENGI" +
		"NE = ReplacingMergeTree(event_time) PARTITION BY toYYYYMM(event_time) ORDER BY (customer_id, resource_id) SETTIN" +
		"GS index_granularity = 8192"
	outputCreateTable, err := ioutil.ReadFile("test/table.sql")

	inputCreateView := "CREATE MATERIALIZED VIEW testdb.watch_progress_handler TO testdb.watch_progress (`customer_id` " +
		"String, `resource_id` String, `resource_type` String, `timestamp` Int32, `duration` Int32, `event_time` DateTime" +
		", `session_id` DateTime, `action` String) AS SELECT toString(customer_id) AS customer_id, toString(stream_id) AS" +
		" resource_id, resource_type, timestamp, duration, toDateTime(event_time) AS event_time, toString(session_id) AS s" +
		"ession_id, action FROM testdb.ping_events_stream"
	outputCreateView, err := ioutil.ReadFile("test/mat_view.sql")

	inputCreateDatabase := "CREATE DATABASE testdb ENGINE = Ordinary"
	outputCreateDatabase, err := ioutil.ReadFile("test/database.sql")

	if err != nil {
		panic(err)
	}

	examples := []struct {
		input  string
		output string
	}{
		{
			inputCreateTable,
			string(outputCreateTable),
		},
		{
			inputCreateView,
			string(outputCreateView),
		},
		{
			inputCreateDatabase,
			string(outputCreateDatabase),
		},
	}

	for _, example := range examples {
		result := Prettify(example.input)
		if result != example.output {
			t.Fatalf("Prettified text doesnt match.\n\tExpected:\n\t%v\n\tGot:\n\t%v\n", example.output, result)
		}
	}
}
