package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"os"
	"os/exec"
)

type Scope struct {
	Name        string              `json:"Scope"`
	Collections []map[string]string `json:"Collections"`
}

type Res struct {
	results []string
}

func pgExport(scope string, table string, collection string, pool *pgxpool.Pool) {

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
		os.Exit(1)
	}
	defer conn.Release()

	file_name := table + ".csv"
	curr_path, _ := os.Getwd()
	file_path := curr_path + "/" + file_name
	//creates csv files from the data in CB
	query := "copy " + table + " to '" + file_path + "' delimiter ',' csv header;"

	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error in query:", err)
		os.Exit(1)
	}

	cbImport(file_name, scope, table, collection)

}

func cbImport(filename string, scope string, table string, collection string) {

	csvfile, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", (err))
		os.Exit(1)
	}

	r := csv.NewReader(csvfile)

	records, err := r.ReadAll()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", (err))
	}

	var header []string
	header = records[0]

	//create a JSON object to insert
	for i := 1; i < len(records); i++ {
		json_obj := make(map[string]string)
		key := records[i][0]
		for j := 0; j < len(header); j++ {
			json_obj[header[j]] = records[i][j]
		}

		final_json, _ := json.Marshal(json_obj)

		_, _ = exec.Command("/bin/bash", "./insert.sh", scope, collection, string(final_json), key).CombinedOutput()

	}
}

func main() {

	byteValue, _ := ioutil.ReadFile("public.json")
	var scopes []Scope

	json.Unmarshal(byteValue, &scopes)

	var pool *pgxpool.Pool

	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < len(scopes); i++ {
		scope_name := scopes[i].Name + "_scope"
		_, _ = exec.Command("/bin/bash", "scope.sh", scope_name).CombinedOutput()

		for coll, table := range scopes[i].Collections[0] {
			_, _ = exec.Command("/bin/bash", "collection.sh", scope_name, coll).CombinedOutput()
			pgExport(scope_name, table, coll, pool)
		}
	}

}
