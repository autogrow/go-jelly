// +build ignore

package main

import (
	"flag"

	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
	"github.com/AutogrowSystems/go-jelly/sfc"
)

const PB_TOKEN = "YOUR PUSHBULLET TOKEN GOES HERE"

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

	// create a pushbullet client
	pb := pushbullet.New(PB_TOKEN)

	// create night and day notifications that we can reuse
	dayMsg := requests.NewNote()
	dayMsg.Title = "IntelliDose "+sn
	dayMsg.Body = "Switched to day time mode"
	nightMsg := requests.NewNote()
	nightMsg.Title = "IntelliDose "+sn
	nightMsg.Body = "Switched to night time mode"

	// Send the note via Pushbullet.
	if _, err := pb.PostPushesNote(n); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}

	inDayTime = false
	for {
		// block until the IntelliDose sends a new packet
		idose.WaitForUpdate()

		// don't spam
		if inDayTime && idose.IsDayTime() || !inDayTime && !idose.IsDayTime() {
			continue
		}

		// transition to night time
		if inDayTime && !idose.IsDayTime() {
			if _, err := pb.PostPushesNote(nightMsg); err != nil {
				log.Printf("ERROR: %s", err)
			}
		}

		// transition to day time
		if !inDayTime && idose.IsDayTime() {
			if _, err := pb.PostPushesNote(dayMsg); err != nil {
				log.Printf("ERROR: %s", err)
			}
		}
	}
}
