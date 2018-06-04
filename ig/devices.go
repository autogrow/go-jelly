package ig

import (
	"fmt"
	"strings"
	"time"
)

const (
	// MetricsEP - use metric endpoint
	MetricsEP = "metrics"
	// ConfigEP - use config endpoint
	ConfigEP = "config"
	// StateEP - use state endpoint
	StateEP = "state"
)

// Devices - object that contains a slice of intellidosers and intelliclimates
type Devices struct {
	IntelliClimates []*IntelliClimate
	IntelliDoses    []*IntelliDose
}

// NewDevices - creates a new devices object
func NewDevices() *Devices {
	return &Devices{}
}

// IsEmpty returns true if there are no intellis in this devices collection
func (ds *Devices) IsEmpty() bool {
	return len(ds.IntelliClimates) == 0 && len(ds.IntelliDoses) == 0
}

// Add - adds a new device to the devices structure it will assign the devcies to the climate or doser slice depending on its type
func (ds *Devices) Add(newDev *Device) {
	if newDev.IsIClimate() {
		for _, dev := range ds.IntelliClimates {
			if dev.GetID() == newDev.GetID() {
				return
			}
		}
		ic := NewIntelliClimate(newDev)
		ds.IntelliClimates = append(ds.IntelliClimates, ic)
	}

	if newDev.IsIDose() {
		for _, dev := range ds.IntelliDoses {
			if dev.GetID() == newDev.GetID() {
				return
			}
		}
		id := NewIntelliDose(newDev)
		ds.IntelliDoses = append(ds.IntelliDoses, id)
	}
}

// Climates - returns the intelliclimate slice
func (ds *Devices) Climates() []*IntelliClimate {
	return ds.IntelliClimates
}

// GetClimateByID - returns a climate with the id that matches the serial number provided
func (ds *Devices) GetClimateByID(id string) (*IntelliClimate, error) {
	for _, dev := range ds.IntelliClimates {
		if dev.ID == id {
			return dev, nil
		}
	}
	return nil, fmt.Errorf("No device with ID %s found", id)
}

// GetClimateByName - returns a climate with the name that matches the one provided
func (ds *Devices) GetClimateByName(name string) (*IntelliClimate, error) {
	for _, dev := range ds.IntelliClimates {
		if dev.DeviceName == name {
			return dev, nil
		}
	}
	return nil, fmt.Errorf("No device with name %s found", name)
}

// UpdateClimateMetrics - updates all intellicliamte metrics
func (ds *Devices) UpdateClimateMetrics() error {
	var errMsg string
	var anErr bool

	for _, ic := range ds.IntelliClimates {
		err := ic.GetMetrics()
		if err != nil {
			anErr = true
			errMsg += fmt.Sprintf("Error Updating: %s: %s ", ic.GetID(), err)
		}
	}

	if anErr {
		return fmt.Errorf(errMsg)
	}

	return nil
}

// Dosers - returns the intellidosers slice
func (ds *Devices) Dosers() []*IntelliDose {
	return ds.IntelliDoses
}

// GetDoserByID - returns a doser with the id that matches the serial number provided
func (ds *Devices) GetDoserByID(id string) (*IntelliDose, error) {
	for _, dev := range ds.IntelliDoses {
		if dev.ID == id {
			return dev, nil
		}
	}
	return nil, fmt.Errorf("No device with ID %s found", id)
}

// GetDoserByName - returns a doser with the name that matches the one provided
func (ds *Devices) GetDoserByName(name string) (*IntelliDose, error) {
	for _, dev := range ds.IntelliDoses {
		if dev.DeviceName == name {
			return dev, nil
		}
	}
	return nil, fmt.Errorf("No device with name %s found", name)
}

// UpdateDoserMetrics - updates all intellidosers metrics
func (ds *Devices) UpdateDoserMetrics() error {
	var errMsg string
	var anErr bool

	for _, id := range ds.IntelliDoses {
		err := id.GetMetrics()
		if err != nil {
			anErr = true
			errMsg += fmt.Sprintf("Error Updating: %s: %s ", id.GetID(), err)
		}
	}

	if anErr {
		return fmt.Errorf(errMsg)
	}

	return nil
}

// Device - general structure that holds devices
type Device struct {
	ID             string  `json:"device_id"`
	Type           string  `json:"device_type"`
	Growroom       string  `json:"growroom"`
	Checked        float64 `json:"schecked"`
	SchedulingMode string  `json:"scheduling_mode"`
	LastUpdated    float64 `json:"last_updated"`
	TimeZoneOffset float64 `json:"time_zone_offset"`
	DeviceName     string  `json:"device_name"`
	client         *Client
	Readings       map[string]interface{}
}

// Average
type devices []*Device

func (ds devices) Average(field string) float64 {
	var sum float64
	var count int

	for _, d := range ds {
		if d.IsValid() {
			sum += d.Readings[field].(float64)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

// IntelliDose returns this device as an IntelliDose, or an error if this device is
// not an IntelliDose
func (d *Device) IntelliDose() (*IntelliDose, error) {
	if d.Type != "idoze" {
		return nil, fmt.Errorf("this device is not an IntelliDose")
	}

	return NewIntelliDose(d), nil
}

// IntelliClimate returns this device as an IntelliClimate, or an error if this device is
// not an IntelliClimate
func (d *Device) IntelliClimate() (*IntelliClimate, error) {
	if d.Type != "iclimate" {
		return nil, fmt.Errorf("this device is not an IntelliClimate")
	}

	return NewIntelliClimate(d), nil
}

// AttachClient - since all information relating to a device has to be got specifically from a device endpoint it makes sense to
// attach a valid API client to the device.
func (d *Device) AttachClient(c *Client) {
	d.client = c
}

// GetID - returns the ID for a device
func (d *Device) GetID() string {
	return d.ID
}

// GetType - return type for a device
func (d *Device) GetType() string {
	return d.Type
}

// IsIClimate - returns a true if the device is a IntelliClimate
func (d *Device) IsIClimate() bool {
	if strings.Contains(d.ID, "IC") {
		return true
	}
	return false
}

// IsIDose - returns a true if the device is a IntelliDose
func (d *Device) IsIDose() bool {
	if strings.Contains(d.ID, "ID") {
		return true
	}
	return false
}

// IsValid - returns a true if the last updated time is with 1 minute of now
func (d *Device) IsValid() bool {
	ld := int64(d.LastUpdated / 1000)
	t := time.Now().Unix()
	if (t - ld) > 6000 {
		return false
	}
	return true
}

// GetGrowroom - returns the growroom name this device is assigned to
func (d *Device) GetGrowroom() string {
	return d.Growroom
}
