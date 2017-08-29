package mydata

import (
	"encoding/json"
	"fmt"
	"time"
)

// Point represents a single data point
type Point struct {
	Timestamp float64 `json:"time"`
	Value     float64 `json:"value"`
}

type RecordStringer interface {
	json.Unmarshaler
	Lines() []string
	MetricString() string
}

// Points represents a collection of points
type Points []Point

// CompRecords represents a collection of CompRecord objects
type CompRecords []*CompRecord

// SiteRecords represents a collection of SiteRecord objects
type SiteRecords []*SiteRecord

// FieldRecords represents a collection of FieldRecord objects
type FieldRecords []*FieldRecord

// IrrigRecords represents a collection of IrrigRecord objects
type IrrigRecords []*IrrigRecord

// RoomRecords represents a collection of RoomRecord objects
type RoomRecords []*RoomRecord

// MonitorRecords represents a collection of MonitorRecord objects
type MonitorRecords []*MonitorRecord

// CompRecord represents a record containing the metrics of a Comp
type CompRecord struct {
	// time in unixtimestamp format
	Time int64 `json:"time"`

	// temperature
	Temp float64 `json:"temp,omitempty"`

	// dew point temperature
	Dewpoint float64 `json:"dewpoint,omitempty"`

	// relative humidity
	Rh float64 `json:"rh,omitempty"`

	// vvapour-pressure deficit
	Vpd float64 `json:"vpd,omitempty"`

	// co2 density
	Co2 float64 `json:"co2,omitempty"`

	// solar radiation
	SolarRad float64 `json:"solar_rad,omitempty"`

	// solar par
	SolarPar float64 `json:"solar_par,omitempty"`
}

func (c CompRecord) String() string {
	time := time.Unix(c.Time, 0)
	return fmt.Sprintf("%-20s %-10.2f %-10.2f %-10.f %-10.2f %-10.f %-10.2f %-10.2f", time, c.Temp, c.Dewpoint, c.Rh, c.Vpd, c.Co2, c.SolarRad, c.SolarPar)
}

// MetricString returns a formatted string of the Metric fields on one line
func (c CompRecords) MetricString() string {
	return fmt.Sprintf("%-20s %-10s %-10s %-10s %-10s %-10s %-10s %-10s", "Timestamp", "Temp", "Dewpoint", "Rh", "Vpd", "Co2", "SolarRad", "SolarPar")
}

func (c CompRecords) Lines() []string {
	lines := []string{}

	for _, r := range c {
		lines = append(lines, r.String())
	}

	return lines
}

func (c CompRecords) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &c) }

// SiteRecord represents a record containing the metrics of a Site
type SiteRecord struct {
	// time in unixtimestamp format
	Time int64 `json:"time"`

	// temperature
	Temp float64 `json:"temp,omitempty"`

	// temperature of plate
	PlateTemp float64 `json:"plate_temp,omitempty"`

	// dew point temperature
	Dewpoint float64 `json:"dewpoint,omitempty"`

	// relative humidity
	Rh float64 `json:"rh,omitempty"`

	// wind speed
	Windspeed float64 `json:"windspeed,omitempty"`

	// vvapour-pressure deficit
	Vpd float64 `json:"vpd,omitempty"`

	// baro pressure
	Pressure float64 `json:"pressure,omitempty"`

	// co2 density
	Co2 float64 `json:"co2,omitempty"`

	// solar radiation
	SolarRad float64 `json:"solar_rad,omitempty"`

	// solar par
	SolarPar float64 `json:"solar_par,omitempty"`
}

func (c SiteRecord) String() string {
	time := time.Unix(c.Time, 0)
	return fmt.Sprintf("%-20s %-10.2f %-10.2f %-10.2f %-10.f %-10.2f %-10.2f %-10.f %-10.f %-10.2f %-10.2f", time, c.Temp, c.PlateTemp, c.Dewpoint, c.Rh, c.Windspeed, c.Vpd, c.Pressure, c.Co2, c.SolarRad, c.SolarPar)
}

// MetricString returns a formatted string of the Metric fields on one line
func (c SiteRecords) MetricString() string {
	return fmt.Sprintf("%-20s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s %-10s", "Timestamp", "Temp", "PlateTemp", "Dewpoint", "Rh", "Windspeed", "Vpd", "Pressure", "Co2", "SolarRad", "SolarPar")
}

func (c SiteRecords) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c)
}

func (c SiteRecords) Lines() []string {
	lines := []string{}

	for _, r := range c {
		lines = append(lines, r.String())
	}

	return lines
}

// FieldRecord represents a record containing the metrics of a Field
type FieldRecord struct {
	// time in unixtimestamp format
	Time int64 `json:"time"`

	// temperature
	AirTemp float64 `json:"air_temp,omitempty"`

	// dew point temperature
	Dewpoint float64 `json:"dewpoint,omitempty"`

	// plate temperature
	PlateTemp float64 `json:"plate_temp,omitempty"`

	// soil temperature
	SoilTemp float64 `json:"soil_temp,omitempty"`
}

func (c FieldRecord) String() string {
	time := time.Unix(c.Time, 0)
	return fmt.Sprintf("%-20s %-10.2f %-10.2f %-10.2f %-10.2f", time, c.AirTemp, c.PlateTemp, c.Dewpoint, c.SoilTemp)
}

// MetricString returns a formatted string of the Metric fields on one line
func (c FieldRecords) MetricString() string {
	return fmt.Sprintf("%-20s %-10s %-10s %-10s %-10s", "Timestamp", "Temp", "PlateTemp", "Dewpoint", "SoilTemp")
}

func (c FieldRecords) Lines() []string {
	lines := []string{}

	for _, r := range c {
		lines = append(lines, r.String())
	}

	return lines
}

func (c FieldRecords) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &c) }

// IrrigRecord represents a record containing the metrics of a Irrig
type IrrigRecord struct {
	// time in unixtimestamp format
	Time int64 `json:"time"`

	// electrictiy conductivity
	Ec float64 `json:"ec,omitempty"`

	// ec temperature
	EcTemp float64 `json:"ec_temp,omitempty"`

	// water flow rate
	Flowrate float64 `json:"flowrate,omitempty"`

	// PH
	Ph float64 `json:"ph,omitempty"`
}

func (c IrrigRecord) String() string {
	time := time.Unix(c.Time, 0)
	return fmt.Sprintf("%-20s %-10.2f %-10.2f %-10.2f %-10.2f", time, c.Ec, c.EcTemp, c.Flowrate, c.Ph)
}

// MetricString returns a formatted string of the Metric fields on one line
func (c IrrigRecords) MetricString() string {
	return fmt.Sprintf("%-20s %-10s %-10s %-10s %-10s", "Timestamp", "EC", "ECTemp", "Flowrate", "pH")
}

func (c IrrigRecords) Lines() []string {
	lines := []string{}

	for _, r := range c {
		lines = append(lines, r.String())
	}

	return lines
}

func (c IrrigRecords) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &c) }

// RoomRecord represents a record containing the metrics of a Room
type RoomRecord struct {
	// time in unixtimestamp format
	Time int64 `json:"time"`

	// temperature
	Temp float64 `json:"temp,omitempty"`

	// relative humidity
	Rh float64 `json:"rh,omitempty"`

	// co2 density
	Co2 float64 `json:"co2,omitempty"`
}

func (c RoomRecord) String() string {
	time := time.Unix(c.Time, 0)
	return fmt.Sprintf("%-20s %-10.2f %-10.f %-10.f", time, c.Temp, c.Rh, c.Co2)
}

// MetricString returns a formatted string of the Metric fields on one line
func (c RoomRecords) MetricString() string {
	return fmt.Sprintf("%-20s %-10s %-10s %-10s", "Timestamp", "Temp", "RH", "CO2")
}

func (c RoomRecords) Lines() []string {
	lines := []string{}

	for _, r := range c {
		lines = append(lines, r.String())
	}

	return lines
}

func (c RoomRecords) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &c) }

// MonitorRecord represents a record containing the metrics of a Monitor
type MonitorRecord struct {
	// time in unixtimestamp format
	Time int64 `json:"time"`

	// temperature
	Temp float64 `json:"temp,omitempty"`

	// dew point temperature
	Dewpoint float64 `json:"dewpoint,omitempty"`

	// relative humidity
	Rh float64 `json:"rh,omitempty"`

	// vvapour-pressure deficit
	Vpd float64 `json:"vpd,omitempty"`

	// co2 density
	Co2 float64 `json:"co2,omitempty"`

	// solar radiation
	SolarRad float64 `json:"solar_rad,omitempty"`

	// solar par
	SolarPar float64 `json:"solar_par,omitempty"`
}

func (c MonitorRecord) String() string {
	time := time.Unix(c.Time, 0)
	return fmt.Sprintf("%-20s %-10.2f %-10.2f %-10.f %-10.2f %-10.f %-10.2f %-10.2f", time, c.Temp, c.Dewpoint, c.Rh, c.Vpd, c.Co2, c.SolarRad, c.SolarPar)
}

// MetricString returns a formatted string of the Metric fields on one line
func (c MonitorRecords) MetricString() string {
	return fmt.Sprintf("%-20s %-10s %-10s %-10s %-10s %-10s %-10s %-10s", "Timestamp", "Temp", "Dewpoint", "RH", "VPD", "CO2", "SolarRAD", "SolarPAR")
}

func (c MonitorRecords) Lines() []string {
	lines := []string{}

	for _, r := range c {
		lines = append(lines, r.String())
	}

	return lines
}

func (c MonitorRecords) UnmarshalJSON(data []byte) error { return json.Unmarshal(data, &c) }

// RecordsFactory provides syntactical sugar for getting new record collection objects
type RecordsFactory struct{}

// Site returns an empty collection of Site records
func (fact RecordsFactory) Site() SiteRecords { return SiteRecords{} }

// Compartment returns an empty collection of Compartment records
func (fact RecordsFactory) Compartment() CompRecords { return CompRecords{} }

// Field returns an empty collection of Field records
func (fact RecordsFactory) Field() FieldRecords { return FieldRecords{} }

// GrowRoom returns an empty collection of GrowRoom records
func (fact RecordsFactory) GrowRoom() RoomRecords { return RoomRecords{} }

// Irrigator returns an empty collection of Irrigator records
func (fact RecordsFactory) Irrigator() IrrigRecords { return IrrigRecords{} }

// Monitor returns an empty collection of Monitor records
func (fact RecordsFactory) Monitor() MonitorRecords { return MonitorRecords{} }

// Records provides a syntax sugar to getting new record collections
var Records = RecordsFactory{}
