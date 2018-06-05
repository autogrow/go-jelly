package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/autogrow/go-jelly/ig"
)

type creds struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

func main() {
	var listDevices, listGrowrooms bool
	var id, gr string
	var printReadings, fmtJSON bool
	flag.BoolVar(&listDevices, "l", false, "list known devices")
	flag.BoolVar(&listGrowrooms, "g", false, "list growrooms")
	flag.StringVar(&id, "id", "", "serial number to work with")
	flag.StringVar(&gr, "gr", "", "growroom name to work with")
	flag.BoolVar(&printReadings, "r", false, "print readings")
	flag.BoolVar(&fmtJSON, "json", false, "format as JSON")
	flag.Parse()

	credsFile := os.Getenv("HOME") + "/.intelligrow/creds"
	creds, err := readCreds(credsFile)
	if err != nil {
		initCreds(credsFile)
		log.Fatalf("you need to enter your credentials in the file %s to continue", credsFile)
	}

	cl, err := ig.NewClient(creds.User, creds.Pass)
	if err != nil {
		log.Fatalf("failed to create IG client: %s", err)
	}

	switch {
	case listGrowrooms:
		fmt.Println("Growrooms:")
		for _, name := range cl.ListGrowrooms() {
			fmt.Printf("- %s\n", name)
		}

	case listDevices:
		fmt.Printf("%-12s %-18s %-10s %s\n", "Type", "ID", "Name", "Growroom")
		for _, d := range cl.Devices() {
			fmt.Printf("%-12s %-18s %-10s %s\n", d.Type, d.ID, d.DeviceName, d.Growroom)
		}

	case id != "":
		if doser, errd := cl.IntelliDose(id); errd == nil {
			switch {
			case printReadings:
				_ = doser.GetMetrics()
				fmt.Printf("%20s: %0.2f mS/cm²\n", "Nutrient", doser.Metrics.Ec)
				fmt.Printf("%20s: %0.2f pH\n", "Acidity", doser.Metrics.PH)
				fmt.Printf("%20s: %0.2f °C\n", "Water", doser.Metrics.NutTemp)
				return
			case fmtJSON:
				if err := doser.GetAll(); err != nil {
					log.Fatalf("failed to get the doser data: %s", err)
				}
				dumpJSON(doser)
				return
			}
		}

		if clim, errc := cl.IntelliClimate(id); errc == nil {
			switch {
			case printReadings:
				_ = clim.GetMetrics()
				fmt.Printf("%20s: %0.2f °C\n", "Air", clim.Metrics.AirTemp)
				fmt.Printf("%20s: %0.2f %%H\n", "RH", clim.Metrics.Rh)
				fmt.Printf("%20s: %0.2f kPa\n", "VPD", clim.Metrics.Vpd)
				fmt.Printf("%20s: %0.2f ppm\n", "CO2", clim.Metrics.Co2)
				return
			case fmtJSON:
				if err := clim.GetAll(); err != nil {
					log.Fatalf("failed to get the climate data: %s", err)
				}
				dumpJSON(clim)
				return
			}
		}

		log.Fatalf("no device found with serial %s", id)

	case gr != "" && printReadings:
		room, ok := cl.Growroom(gr)

		if !ok {
			log.Fatalf("Growroom %s not found", gr)
		}

		ics, _ := room.IntelliClimates()
		if len(ics) > 0 {
			fmt.Printf("%20s: %0.2f °C\n", "Air", room.Climate.AirTemp)
			fmt.Printf("%20s: %0.2f %%H\n", "RH", room.Climate.RH)
			fmt.Printf("%20s: %0.2f kPa\n", "VPD", room.Climate.VPD)
			fmt.Printf("%20s: %0.2f ppm\n", "CO2", room.Climate.CO2)
		}

		ids, _ := room.IntelliDoses()
		if len(ids) > 0 {
			fmt.Printf("%20s: %0.2f mS/cm²\n", "Nutrient", room.Rootzone.EC)
			fmt.Printf("%20s: %0.2f pH\n", "Acidity", room.Rootzone.PH)
			fmt.Printf("%20s: %0.2f °C\n", "Water", room.Rootzone.Temp)
		}
	}
}

func initCreds(credsFile string) {
	data, err := json.Marshal(creds{})
	if err != nil {
		log.Fatalf("failed to encode empty creds file to put at %s", credsFile)
	}

	if err := os.Mkdir(filepath.Dir(credsFile), 0755); err != nil {
		log.Fatalf("failed to create creds file at %s: %s", credsFile, err)
	}

	if err := ioutil.WriteFile(credsFile, data, 0644); err != nil {
		log.Fatalf("failed to create creds file at %s: %s", credsFile, err)
	}
}

func readCreds(credsFile string) (creds, error) {
	creds := creds{}
	data, err := ioutil.ReadFile(credsFile)
	if err != nil {
		return creds, err
	}

	err = json.Unmarshal(data, &creds)
	if err != nil {
		log.Fatalf("failed to read creds file at %s: %s", credsFile, err)
	}

	if creds.Pass == "" || creds.User == "" {
		log.Fatalf("you need to enter your credentials in the file %s to continue", credsFile)
	}

	return creds, nil
}

func dumpJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("failed to format JSON: %s", err)
	}

	fmt.Println(string(data))
}
