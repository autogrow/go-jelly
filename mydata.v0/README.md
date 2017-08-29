# MyData API SDK

This SDK can be used to get programmatic access to the data contained in the Autogrow MyData API.

## Installation

Use `go get` as normal:

    $ go get github.com/AutogrowSystems/go-jelly/mydata.v0

A simple command line client can be build by running:

    $ go install github.com/AutogrowSystems/go-jelly/mydata.v0/cmd/mydata

## Usage

First import it:

```go
import (
  "github.com/AutogrowSystems/go-jelly/mydata.v0"
)
```

Create a new connection, giving your device ID as the argument:

```go
api := mydata.Dial("000000111111222222aaaaaabbbbbbccccccdd")
```

You'll want to know the names of the sections that you can query from the API.  You can do this by checking the labels that are available:

```go
labels, err := api.Labels()

for _, lbl := range labels {
  fmt.Println(lbl.Type, lbl.Name)
}
```

That returns a slice of label structures containing the type of section (e.g. Compartment, Irrigator, Grow Room) and the label name for that
section (e.g. Lettuces, Comp. 1, Tunnel House):

```
Compartment Lettuces
Compartment Comp. 1
Irrigation Tunnel House
```

### Singular Metrics

Then you can start to query for metrics using the chainable interface.  First we call the method for the environmental section
that you wish to query, giving it the name you have set for it.  Then you can set the metric to get data for, along with the
number of days back in time to go.  The API is not actually called until the call to the `Points()` method which returns the data points.

```go
q := api.Compartment("Lettuces")
q.Metric("temp")
q.Days(10)
points, err := q.Points()
```

Or the same thing in one line:

```go
points, err := api.Compartment("Lettuces").Metric("temp").Days(10).Points()
```

Points is returned as a slice of `Point` structs, each with a `Timestamp` (unix epoch) and `Value` field, both of type `float64`.

```go
fmt.Println("Lettuce Compartment Temps")

for _, pt := range points {
  t := time.Unix(int64(pt.Timestamp), 0)
  fmt.Printf("%s\t%0.f\n", t, pt.Value)
}
```

This will produce output similar to:

```
Lettuce Compartment Temps
2017-07-16 10:00:00 +0000 UTC    22.45
2017-07-16 11:00:00 +0000 UTC    22.80
2017-07-16 12:00:00 +0000 UTC    22.90
2017-07-16 13:00:00 +0000 UTC    22.72
2017-07-16 14:00:00 +0000 UTC    22.70
2017-07-16 15:00:00 +0000 UTC    22.65
2017-07-16 16:00:00 +0000 UTC    22.70
2017-07-16 17:00:00 +0000 UTC    22.52
```

### Multiple Metrics

Getting multiple metrics is a similar operation to the singular method, but the `All()` method is called instead.  Instead of a simple data point
structure being returned, a structure containing all the metrics paired with the timestamp is returned.

```go
records := mydata.Records.Compartment()
err := api.Compartment("Lettuces").Days(10).All(&records)

rec := records[0]
fmt.Println("First point:")
fmt.Println("Timestamp:", time.Unix(int64(rec.Timestamp), 0))
fmt.Println("Temp:", rec.Temp)
fmt.Println("RH:", rec.Rh)
fmt.Println("VPD:", rec.Vpd)
```

## CLI Usage

If you have built the command line client as above, the weather stations temps for one day can be queried like so:

```
mydata -site -m temp -id 000000111111222222aaaaaabbbbbbccccccdd -d 1
```

Or the lettuce temps for 10 days:

```
mydata -comp -n Lettuces -m temp -id 000000111111222222aaaaaabbbbbbccccccdd -d 10
```

Or the labels you can query:

```
mydata -labels -id 000000111111222222aaaaaabbbbbbccccccdd
TYPE                 NAME
Compartment          Lettuces
Compartment          Comp. 1
Irrigation           Tunnel House
```

You can also see a summary of the given dataset:

```
mydata -summary -site -m temp -id 000000111111222222aaaaaabbbbbbccccccdd -d 1
Environment:   Site
Name:          
Metric:        temp
Time Period:   1 days
Data points:   25
First point:   2017-07-13 10:00:00 +1200 NZST
Last point:    2017-07-14 10:00:00 +1200 NZST
Min value:     0
Max value:     6.6755555555555555
```

The usage can be obtained by using the `-h` argument:

```
-d int         Get data from this many days ago (default 10)
-id string     Device ID to use
-m string      Metric to pull
-n string      Name of the env section to use
-site          Query site
-comp          Query compartments
-irrig         Query irrigation systems
-room          Query grow rooms
-mon           Query monitors
-field         Query fields
-labels        Query what sections/labels are available
```
