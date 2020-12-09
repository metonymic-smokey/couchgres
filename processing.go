package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Collections struct {
	Name string `json:"name"`
}

type Scopes struct {
	Coll []Collections `json:"collections"`
	Name string        `json:"name"`
}

type Result struct {
	Scope []Scopes `json:"scopes"`
}

var org []Scopes

func generatePieItems(scope Scopes, scope_items map[string]int, coll_items map[string]int) []opts.PieData {

	items := make([]opts.PieData, 0)
	for i := 0; i < len(scope.Coll); i++ {
		scope_name := scope.Name
		coll_name := scope.Coll[i].Name
		items = append(items, opts.PieData{Name: coll_name, Value: 100 * coll_items[coll_name] / scope_items[scope_name]})
	}

	return items
}

func pieBase(scope Scopes, scope_items map[string]int, coll_items map[string]int) *charts.Pie {

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: scope.Name + " scope collections"}),
	)

	pie.AddSeries("pie", generatePieItems(scope, scope_items, coll_items)).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
			}),
		)

	return pie
}

func itemDetails() (map[string]int, map[string]int) {
	file, err := os.Open("details")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	line_no := 0
	var scope_lines []int
	var textlines []string

	coll_items := make(map[string]int)
	scope_items := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line_no += 1
		textlines = append(textlines, scanner.Text())
		match, _ := regexp.Match("scope_name", []byte(scanner.Text()))

		if match {
			scope_lines = append(scope_lines, line_no)

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(scope_lines); i++ {
		items := strings.Fields(textlines[scope_lines[i]-7])[1]
		coll := strings.Fields(textlines[scope_lines[i]-5])[1]
		temp, _ := strconv.Atoi(items)

		coll_items[coll] = temp
		scope := strings.Fields(textlines[scope_lines[i]-1])[1]
		scope_items[scope] += temp
	}

	return scope_items, coll_items

}

func main() {

	scope_items, coll_items := itemDetails()
	fmt.Println(scope_items, coll_items)

	_ = exec.Command("listScopes.sh")

	jsonFile, err := os.Open("scopes.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var results Result

	json.Unmarshal(byteValue, &results)

	num_scopes := len(results.Scope)
	scopes := make(map[string][]Collections)

	for i := 0; i < num_scopes; i++ {
		scopes[results.Scope[i].Name] = results.Scope[i].Coll
	}

	for k, v := range scopes {
		fmt.Println(k, v)
	}

	f, _ := os.Create("coll.html")

	page := components.NewPage()

	for i := 0; i < len(results.Scope); i++ {

		page.AddCharts(
			pieBase(results.Scope[i], scope_items, coll_items),
		)
	}

	page.Render(io.MultiWriter(f))

}
