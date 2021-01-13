/*
find all schemas in a DB
create scopes for each in a CB bucket
go through each table in every schema
create a collection and import data into the right scope.
make a json document showing schema-table and scope-collection relationship
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"os"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "aditi"
	password = ""
	dbname   = "public"
)

var psqlInfo string

type Schema struct {
	Scope       string
	Collections []map[string]string
}

func main() {

	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	conn, err := pgx.Connect(context.Background(), psqlInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var schema, table string
	query := "SELECT schemaname,tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';"
	//query := "select schemaname,tablename from pg_catalog.pg_tables;"
	rows, err := conn.Query(context.Background(), query)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	p := map[string][]string{}

	for rows.Next() {

		err := rows.Scan(&schema, &table)
		if err != nil {
			fmt.Println(err)
		}

		p[schema] = append(p[schema], table)

	}

	var final []Schema

	for k, v := range p {
		x := Schema{
			Scope: k,
		}
		temp_map := map[string]string{}
		for _, c := range v {
			temp_map[c] = k + "." + c
		}
		x.Collections = append(x.Collections, temp_map)
		final = append(final, x)
	}

	temp, _ := json.MarshalIndent(final, "", " ")

	_ = ioutil.WriteFile("public.json", temp, 0644)

}
