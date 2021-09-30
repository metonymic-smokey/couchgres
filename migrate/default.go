package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"sync"
)

var mode *string

type Scope struct {
	Name        string              `json:"Scope"`
	Key         string              `json:"Key"`
	Collections []map[string]string `json:"Collections"`
}

func main() {

	mode = flag.String("mode", "app", "options: docker/app")
	flag.Parse()

	byteValue, _ := ioutil.ReadFile("split.json")
	var Bucket_Org []Scope
	json.Unmarshal(byteValue, &Bucket_Org)

	for _, x := range Bucket_Org {
		_, err := exec.Command("/bin/bash", "scope.sh", x.Name).CombinedOutput()
		if err != nil {
			panic(err)
		}
		c := x.Collections[0]

		//attempting to concurrently create collections
		var wg sync.WaitGroup
		wg.Add(len(c))
		for coll := range c {
			go func(coll string) {
				defer wg.Done()
				_, err := exec.Command("/bin/bash", "collection.sh", x.Name, coll).CombinedOutput()
				if err != nil {
					panic(err)
				}
			}(coll)
		}
		wg.Wait()

	}

	for _, x := range Bucket_Org {
		c := x.Collections[0]

		var wg sync.WaitGroup
		wg.Add(len(c))

		for coll, val := range c {
			go func(coll string, val string) {
				defer wg.Done()
				op, err := exec.Command("/bin/bash", "divide.sh", x.Name, coll, x.Key, val).CombinedOutput()
				fmt.Println("divide op: ", string(op))
				if err != nil {
					panic(err)
				}
				if *mode == "docker" {
					_, err := exec.Command("/bin/bash", "./cbimport_json.sh", x.Name, coll, "final_res.json", x.Key).CombinedOutput()
					if err != nil {
						panic(err)
					}
				} else {
					_, err := exec.Command("/bin/bash", "./app_cbimport_json.sh", x.Name, coll, "final_res.json", x.Key).CombinedOutput()
					if err != nil {
						panic(err)
					}
				}
			}(coll, val)
		}
		wg.Wait()
	}
}
