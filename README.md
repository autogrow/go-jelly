# go-jelly

[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/AutogrowSystems/go-jelly)

SDK for Autogrow devices and associated integrations.

Traditionally agricultural providers utilise proprietary products and services which leave little flexibility for the growers to create their own solutions or connect them with existing infrastructure. This SDK is a work in progress and aims to provide customers with programmatic access to Autogrow devices.

## Installation

Install as you would any go library:

    go get github.com/AutogrowSystems/go-jelly

## Usage

### IntelliGrow

[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/AutogrowSystems/go-jelly/ig)

IntelliDoses and IntelliClimates are supported via the IntelliGrow API when used in conjuction with an IntelliLink.

```go
import "github.com/AutogrowSystems/go-jelly/ig"

user := "me"
pass := "secret"

client := ig.NewClient(user, pass)

devices, err := client.Devices()
if err != nil {
    panic(err)
}

// print the known devices
for _, d := range devices {
    fmt.Printf("%-12s %-18s %-20s %s", d.Type, d.ID, d.DeviceName, d.Growroom)
}

doser, err := client.GetIntelliDoseByName("ASLID17081149")
if err != nil {
    panic(err)
}

if doser.ForceIrrigation() {
    err := doser.UpdateState()
    if err != nil {
        panic(err)
    }
}
```

### Multigrow

Coming soon!

### go-intelli

Currently the IntelliDose is supported when used in combination with the [go-intelli](https://github.com/AutogrowSystems/go-intelli) gateway for event driven readings.

```go
import "github.com/AutogrowSystems/go-jelly/sfc"

// connect to the IntelliDose on the IntelliLink
idose := sfc.NewIntelliDose(sn)
err := sfc.ConnectToNATS(idose, "nats://localhost:4222", 10)
if err != nil {
    panic(err)
}

for {
    // block until we get an update
    idose.WaitForUpdate()

    // print out the readings
    r := idose.Readings()
    log.Printf("EC: %02.f pH: %02.f Temp: %0.2f", r.Ec, r.PH, r.NutTemp)
}
```

You can see some usage examples in this repo:

- **examples/daynightonoff.go**: send a push notification when an IntelliDose transitions from day to night (or vice versa)
- **examples/printreadings.go**: print readings to the terminal every time they change
