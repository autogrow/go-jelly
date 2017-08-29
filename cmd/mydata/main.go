package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/AutogrowSystems/go-jelly/mydata.v0"
)

func main() {
	var uuid string
	var metric string
	var name string
	var site bool
	var summary bool
	var comp bool
	var room bool
	var irrig bool
	var mon bool
	var field bool
	var labels bool
	var days int

	// setup/selecton stuff
	flag.StringVar(&uuid, "id", "", "Device ID to use")
	flag.StringVar(&metric, "m", "", "Metric to pull")
	flag.StringVar(&name, "n", "", "Name of the env section to use")
	flag.IntVar(&days, "d", 10, "Get data from this many days ago")

	// environmental section selection
	flag.BoolVar(&site, "site", false, "Query site")
	flag.BoolVar(&comp, "comp", false, "Query compartments")
	flag.BoolVar(&room, "room", false, "Query grow rooms")
	flag.BoolVar(&irrig, "irrig", false, "Query irrigation systems")
	flag.BoolVar(&mon, "mon", false, "Query monitors")
	flag.BoolVar(&field, "field", false, "Query fields")

	// actions
	flag.BoolVar(&summary, "summary", false, "print summary of data")
	flag.BoolVar(&labels, "labels", false, "Show the available section labels")

	flag.Parse()

	api := mydata.Dial(uuid)

	if labels {
		printLabels(api)
		os.Exit(0)
	}

	if !site && name == "" {
		panic("must specify name")
	}

	var q *mydata.Query
	var recs mydata.RecordStringer
	var section string

	switch {
	case site:
		q = api.Site()
		section = "Site"
		recs = mydata.Records.Site()
	case comp:
		q = api.Compartment(name)
		section = "Compartment"
		recs = mydata.Records.Compartment()
	case room:
		q = api.GrowRoom(name)
		section = "GrowRoom"
		recs = mydata.Records.GrowRoom()
	case irrig:
		q = api.Irrigator(name)
		section = "Irrigator"
		recs = mydata.Records.Irrigator()
	case mon:
		q = api.Monitor(name)
		section = "Monitor"
		recs = mydata.Records.Monitor()
	case field:
		q = api.Field(name)
		section = "Field"
		recs = mydata.Records.Field()
	default:
		panic("invalid environmental section")
	}

	switch {

	case metric == "":
		err := q.Days(days).All(&recs)
		if err != nil {
			fmt.Println(q.LastURI)
			panic(err)
		}

		printRecords(recs)

	case metric != "":
		points, err := q.Metric(metric).Days(days).Points()
		if err != nil {
			fmt.Println(q.LastURI)
			panic(err)
		}

		if summary {
			printSummary(points, metric, section, name, days)

		} else {
			fmt.Printf("Printing data points for %s in %s %s for the last %d days\n\n", metric, section, name, days)
			printPoints(points)
		}
	}

	os.Exit(0)

}

func printSummary(points mydata.Points, metric, section, name string, days int) {
	first := float64(time.Now().Unix())
	last := 0.0
	min := 99999999.0
	max := 0.0

	for _, pt := range points {
		if pt.Timestamp < first {
			first = pt.Timestamp
		}

		if pt.Timestamp > last {
			last = pt.Timestamp
		}

		if pt.Value < min {
			min = pt.Value
		}

		if pt.Value > max {
			max = pt.Value
		}
	}

	fmt.Println("Environment:  ", section)
	fmt.Println("Name:         ", name)
	fmt.Println("Metric:       ", metric)
	fmt.Println("Time Period:  ", days, "days")
	fmt.Println("Data points:  ", len(points))
	fmt.Println("First point:  ", time.Unix(int64(first), 0))
	fmt.Println("Last point:   ", time.Unix(int64(last), 0))
	fmt.Println("Min value:    ", min)
	fmt.Println("Max value:    ", max)
}

func printPoints(points mydata.Points) {
	for _, pt := range points {
		t := time.Unix(int64(pt.Timestamp), 0)
		fmt.Printf("%s\t%0.2f\n", t, pt.Value)
	}
}

func printLabels(api *mydata.Conn) {
	labels, err := api.Labels()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%-20s %s\n", "TYPE", "NAME")

	for _, l := range labels {
		fmt.Printf("%-20s %s\n", l.Type, l.Name)
	}
}

func printRecords(recs mydata.RecordStringer) {
	fmt.Println(recs.MetricString())

	for _, l := range recs.Lines() {
		fmt.Println(l)
	}
}
