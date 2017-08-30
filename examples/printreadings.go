// +build ignore

package main

import (
	"flag"
	"log"

	"github.com/AutogrowSystems/go-jelly/sfc"
)

func main() {
	var sn string

	flag.StringVar(&sn, "sn", "", "serial number to read metrics")
	flag.Parse()

	if sn == "" {
		panic("must specify serial number")
	}

	// connect to the IntelliDose on the IntelliLink
	idose := sfc.NewIntelliDose(sn)
	err := sfc.ConnectToNATS(idose, "nats://localhost:4222", 10)
	if err != nil {
		panic(err)
	}

	lastR := sfc.MetricsIDose
	for {
		// block until we get an update
		idose.WaitForUpdate()

		// print out the readings
		r := idose.Readings()

		// check if they changed
		if r == lastR {
			continue
		}

		// update the last readings
		lastR = r

		log.Printf("EC: %02.f pH: %02.f Temp: %0.2f", r.Ec, r.PH, r.NutTemp)
	}
}
