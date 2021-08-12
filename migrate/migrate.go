package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var psqlInfo string
var pool *pgxpool.Pool
var mode *string

type ConnInfo struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	User     string `json:"Username"`
	Password string `json:"Password"`
	Database string `json:"Database"`
}

type Scope struct {
	Name        string              `json:"Scope"`
	Collections []map[string]string `json:"Collections"`
}

type Res struct {
	results []string
}

func export(scope string, table string, collection string, pool *pgxpool.Pool) {

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
		os.Exit(1)
	}
	defer conn.Release()

	file_name := table + ".csv"
	curr_path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	file_path := curr_path + "/" + file_name
	//creates csv files from the data in CB
	query := "copy " + table + " to '" + file_path + "' delimiter ',' csv header;"

	_, err = conn.Exec(context.Background(), query)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error in query:", err)
		os.Exit(1)
	}

	index_query := "select t.relname,i.relname,a.attname from pg_class t,pg_class i,pg_attribute a,pg_index ix where t.oid = ix.indrelid and i.oid = ix.indexrelid and a.attrelid = t.oid and a.attnum = ANY(ix.indkey) and t.relkind = 'r' and t.relname = '"

	temp_table := strings.Split(table, ".")[1]
	index_res, err := conn.Query(context.Background(), index_query+temp_table+"';")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error in query:", err)
		os.Exit(1)
	}

	var table_name, index_name, col string
	var columns []string
	for index_res.Next() {
		err := index_res.Scan(&table_name, &index_name, &col)
		if err != nil {
			fmt.Println(err)
		}
		columns = append(columns, col)

	}

	cbImport(file_name, scope, table, collection)
    if len(columns) >= 1 {
		createIndex(scope, collection, index_name, columns)
	}

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
	fmt.Println(*mode)

	if *mode == "docker" {
		op, err := exec.Command("/bin/bash", "./cbimport.sh", scope, collection, filename, header[0]).CombinedOutput()
		fmt.Println(string(op))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		_, _ = exec.Command("/bin/bash", "./app_cbimport.sh", scope, collection, filename, header[0], "csv").CombinedOutput()

	}
}

//separate function for index creation since indices need to be created on collections to avoid errors
func createIndex(scope string, collection string, index_name string, columns []string) {

	col_str := ""
	for _, c := range columns {
		col_str = col_str + c + ","
	}
	col_str = col_str[:len(col_str)-1]

	op, err := exec.Command("/bin/bash", "index.sh", scope, collection, index_name, col_str).CombinedOutput()
	fmt.Println("index output",string(op))
    if err != nil {
		fmt.Println(err)
	}

}

func main() {

	mode = flag.String("mode", "docker", "options: docker/app")
	flag.Parse()

	byteValue, _ := ioutil.ReadFile(".couchgres")
	var creds ConnInfo
	json.Unmarshal(byteValue, &creds)

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", creds.Host, creds.Port, creds.User, creds.Password, creds.Database)

	pool, err := pgxpool.Connect(context.Background(), psqlInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	byteValue, _ = ioutil.ReadFile("public.json")
	var scopes []Scope

	json.Unmarshal(byteValue, &scopes)

	for i := 0; i < len(scopes); i++ {
		scope_name := scopes[i].Name + "_scope"
		_, err = exec.Command("/bin/bash", "scope.sh", scope_name).CombinedOutput()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("creating scope ", scope_name)

		for coll, table := range scopes[i].Collections[0] {
			_, err = exec.Command("/bin/bash", "collection.sh", scope_name, coll).CombinedOutput()
			fmt.Println("creating collection ", coll)
			if err != nil {
				fmt.Println(err)
			}
			export(scope_name, table, coll, pool)
		}
	}

}
