# go-jelly

[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/AutogrowSystems/go-jelly)

SDK for Autogrow devices and associated integrations.

Traditionally agricultural providers utilise proprietary products and services which leave little flexibility for the growers to create their own solutions or connect them with existing infrastructure. This SDK is a work in progress and aims to provide customers with programmatic access to Autogrow devices.

Currently the IntelliDose is supported when used in combination with the IntelliLink, with IntelliClimate support closely following.  The [MyData API](https://www.autogrow.com/mydata) is also supported.

Future products will follow suit and we will try to add older products (such as the Multigrow) as well.

## Installation

Install as you would any go library:

    go get github.com/AutogrowSystems/go-jelly

## Usage

You can see some usage examples in examples directory:

- **daynightonoff.go**: send a push notification when an IntelliDose transitions from day to night (or vice versa)
- **printreadings.go**: print readings to the terminal every time they change
- **cmd/mydata/main.go**: interface with the MyData API

## TODO

- [x] add Autogrow MyData client
- [ ] fix CloudFront breakage of MyData client
- [x] add IntelliDose support through IntelliLink
- [ ] add IntelliClimate support through IntelliLink
- [ ] add IntelliGrow online client so IntelliLink is not needed
- [ ] add Heliospectre light support
- [ ] port to Python (jellypy)
