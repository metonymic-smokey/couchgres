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
	
    _,_ = exec.Command("/bin/bash","./cbimport.sh",scope,collection,filename,header[0]).CombinedOutput()
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
