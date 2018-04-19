# go-jelly

[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/AutogrowSystems/go-jelly)

SDK for Autogrow devices and associated integrations.

Traditionally agricultural providers utilise proprietary products and services which leave little flexibility for the growers to create their own solutions or connect them with existing infrastructure. This SDK is a work in progress and aims to provide customers with programmatic access to Autogrow devices.

Currently the IntelliDose is supported when used in combination with the [go-intelli](https://github.com/AutogrowSystems/go-intelli) gateway.

Future products will follow suit and we will try to add other products (such as the Multigrow) as well.

## Installation

Install as you would any go library:

    go get github.com/AutogrowSystems/go-jelly

## Usage

Basic usage would be like so:

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
- **cmd/mydata/main.go**: interface with the MyData API

## TODO

- [x] add Autogrow MyData client
- [ ] fix CloudFront breakage of MyData client
- [x] add IntelliDose support through IntelliLink
- [ ] add IntelliClimate support through IntelliLink
- [ ] add IntelliGrow online client so IntelliLink is not needed
- [ ] add Heliospectre light support
- [ ] port to Python (jellypy)
