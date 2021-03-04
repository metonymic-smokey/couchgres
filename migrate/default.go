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
	fmt.Println(*mode)

	byteValue, _ := ioutil.ReadFile("split.json")
	var Bucket_Org []Scope
	json.Unmarshal(byteValue, &Bucket_Org)

	for _, x := range Bucket_Org {
		_, _ = exec.Command("/bin/bash", "scope.sh", x.Name).CombinedOutput()
		c := x.Collections[0]

		//attempting to concurrently create collections
		var wg sync.WaitGroup
		wg.Add(len(c))
		for coll, _ := range c {
			go func(coll string) {
				defer wg.Done()
				_, _ = exec.Command("/bin/bash", "collection.sh", x.Name, coll).CombinedOutput()

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
				_, _ = exec.Command("/bin/bash", "divide.sh", x.Name, coll, x.Key, val).CombinedOutput()
				if *mode == "docker" {
					_, _ = exec.Command("/bin/bash", "./cbimport_json.sh", x.Name, coll, "final_res.json", x.Key).CombinedOutput()
				} else {
					_, _ = exec.Command("/bin/bash", "./app_cbimport_json.sh", x.Name, coll, "final_res.json", x.Key).CombinedOutput()

				}
			}(coll, val)

		}
		wg.Wait()
	}
}
