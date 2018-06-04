package ig

import (
	"fmt"
	"sync"
	"time"
)

const (
	grAirTemp        = "air_temp"
	grRH             = "rh"
	grVPD            = "vpd"
	grLight          = "light"
	grPowerFail      = "power_fail"
	grFailSafe       = "fail_safe_alarms"
	grDayNight       = "day_night"
	grCO2            = "co2"
	grLastUpdate     = "last_update"
	grEC             = "ec"
	grPH             = "ph"
	grTemp           = "nut_temp"
	grReadingUnknown = "value not known in climate"
)

// GrowroomClimate - climate data for a single room
type GrowroomClimate struct {
	AirTemp        float64 `json:"air_temp"`
	RH             float64 `json:"rh"`
	VPD            float64 `json:"vpd"`
	Light          float64 `json:"light"`
	PowerFail      bool    `json:"power_fail"`
	FailSafeAlarms bool    `json:"fail_safe_alarms"`
	DayNight       string  `json:"day_night"`
	CO2            float64 `json:"co2"`
	LastUpdate     float64 `json:"last_update"`
}

// GrowroomRootzone - rootzone data for a single room
type GrowroomRootzone struct {
	EC         float64 `json:"ec"`
	PH         float64 `json:"pH"`
	Temp       float64 `json:"nut_temp"`
	LastUpdate float64 `json:"last_update"`
}

// Growroom - object containing all relative information about a single growroom
type Growroom struct {
	*Devices          `json:"devices"`
	lock              *sync.RWMutex
	Name              string            `json:"name"`
	Climate           *GrowroomClimate  `json:"climate"`
	Rootzone          *GrowroomRootzone `json:"rootzone"`
	IntruderAlarm     float64           `json:"intruder_alarm"`
	OutsideTempSensor float64           `json:"outside_temp_sensor"`
}

// NewGrowroom - return a new growroom with the name specified
func NewGrowroom(name string) *Growroom {
	gr := &Growroom{
		NewDevices(),
		new(sync.RWMutex),
		name,
		&GrowroomClimate{},
		&GrowroomRootzone{},
		0,
		0,
	}
	return gr
}

// GetName - returns the name for the growroom
func (g *Growroom) GetName() string {
	g.lock.Lock()
	defer g.lock.Unlock()

	return g.Name
}

// AddDevice - add device to grow room
func (g *Growroom) AddDevice(dev *Device) {
	g.Add(dev)
}

// Update - updated the devices and readings inside the growroom
func (g *Growroom) Update() error {
	var errMsg string
	var aError bool

	err := g.UpdateClimateMetrics()
	if err != nil {
		aError = true
		errMsg = err.Error()
	}

	err = g.UpdateClimate()
	if err != nil {
		aError = true
		errMsg = err.Error()
	}

	err = g.UpdateDoserMetrics()
	if err != nil {
		if !aError {
			errMsg += err.Error()
		} else {
			errMsg = err.Error()
		}
		aError = true
	}

	err = g.UpdateRootzone()
	if err != nil {
		aError = true
		errMsg = err.Error()
	}

	if !aError {
		return nil
	}

	return fmt.Errorf("%s", errMsg)
}

// GetDevices - updated the device inside the growroom
func (g *Growroom) GetDevices() ([]string, []string) {
	var climates []string
	var dosers []string

	for _, d := range g.Climates() {
		climates = append(climates, d.GetID())
	}
	for _, d := range g.Dosers() {
		dosers = append(dosers, d.GetID())
	}
	return climates, dosers
}

// UpdateClimate - takes the readings dict from an intelliclimate and update them into the room climate
func (g *Growroom) UpdateClimate() error {
	// Get my Climates
	climates := g.Climates()
	switch len(climates) {
	case 0:
		return fmt.Errorf("There are no intelliclimates in the growroom")
	case 1:
		// Only one intelliclimate exists don't bother
		g.Climate.LastUpdate = climates[0].LastUpdated
		return updateStruct(climates[0].Readings, g.Climate)
	default:
		g.Climate.AirTemp = AverageClimateReadings(climates, "air_temp")
		g.Climate.RH = AverageClimateReadings(climates, "rh")
		g.Climate.VPD = AverageClimateReadings(climates, "vpd")
		g.Climate.Light = AverageClimateReadings(climates, "light")
		g.Climate.CO2 = AverageClimateReadings(climates, "co2")
		g.Climate.AirTemp = AverageClimateReadings(climates, "air_temp")
		g.Climate.LastUpdate = climates[0].LastUpdated
		g.Climate.FailSafeAlarms = climates[0].Readings["fail_safe_alarms"].(bool)
		g.Climate.PowerFail = climates[0].Readings["power_fail"].(bool)
		g.Climate.DayNight = climates[0].Readings["day_night"].(string)
		return nil
	}
}

// UpdateRootzone - takes the readings dict from an intellidose and update them into the room rootzone
func (g *Growroom) UpdateRootzone() error {
	// Get my Dosers
	dosers := g.Dosers()
	switch len(dosers) {
	case 0:
		return fmt.Errorf("There are no intellicdosers in the growroom")
	case 1:
		// Only one intelliclimate exists don't bother
		g.Rootzone.LastUpdate = dosers[0].LastUpdated
		return updateStruct(dosers[0].Readings, g.Rootzone)
	default:
		g.Rootzone.LastUpdate = dosers[0].LastUpdated
		g.Rootzone.EC = AverageDoseReadings(dosers, "ec")
		g.Rootzone.PH = AverageDoseReadings(dosers, "pH")
		g.Rootzone.Temp = AverageDoseReadings(dosers, "temp")
		return nil
	}
}

// GetReading - returns an avaliable flag and the reading as a string
func (g *Growroom) GetReading(reading string) (bool, string) {
	switch reading {
	case grAirTemp:
		return true, fmt.Sprintf("%.2f", g.Climate.AirTemp)
	case grRH:
		return true, fmt.Sprintf("%.2f", g.Climate.RH)
	case grCO2:
		return true, fmt.Sprintf("%.2f", g.Climate.CO2)
	case grLight:
		return true, fmt.Sprintf("%.2f", g.Climate.Light)
	case grPowerFail:
		return true, fmt.Sprintf("power fall: %t", g.Climate.PowerFail)
	case grFailSafe:
		return true, fmt.Sprintf("fail safe alarm: %t", g.Climate.FailSafeAlarms)
	case grDayNight:
		return true, fmt.Sprintf("day/night: %s", g.Climate.DayNight)
	case grEC:
		return true, fmt.Sprintf("%.2f", g.Rootzone.EC)
	case grPH:
		return true, fmt.Sprintf("%.2f", g.Rootzone.PH)
	case grTemp:
		return true, fmt.Sprintf("%.2f", g.Rootzone.Temp)
	default:
		return false, grReadingUnknown
	}
}

// GetRootzoneReading - returns an avaliable flag and the reading as a string
func (g *Growroom) GetRootzoneReading(reading string) (bool, string) {
	if (time.Now().Unix() - int64(g.Rootzone.LastUpdate)) > 60 {
		return false, ""
	}

	switch reading {
	case grEC:
		return true, fmt.Sprintf("%.2f", g.Rootzone.EC)
	case grPH:
		return true, fmt.Sprintf("%.2f", g.Rootzone.PH)
	case grTemp:
		return true, fmt.Sprintf("%.2f", g.Rootzone.Temp)
	default:
		return false, grReadingUnknown
	}
}

// GetClimateReading - returns an avaliable flag and the reading as a string
func (g *Growroom) GetClimateReading(reading string) (bool, string) {
	if (time.Now().Unix() - int64(g.Climate.LastUpdate)) > 60 {
		return false, ""
	}

	switch reading {
	case grAirTemp:
		return true, fmt.Sprintf("%.2f", g.Climate.AirTemp)
	case grRH:
		return true, fmt.Sprintf("%.2f", g.Climate.RH)
	case grCO2:
		return true, fmt.Sprintf("%.2f", g.Climate.CO2)
	case grLight:
		return true, fmt.Sprintf("%.2f", g.Climate.Light)
	case grPowerFail:
		return true, fmt.Sprintf("power fall: %t", g.Climate.PowerFail)
	case grFailSafe:
		return true, fmt.Sprintf("fail safe alarm: %t", g.Climate.FailSafeAlarms)
	case grDayNight:
		return true, fmt.Sprintf("day/night: %s", g.Climate.DayNight)
	default:
		return false, "value not known in climate"
	}
}

// AirTemp - returns the air temperature for the climate as a string
func (g *Growroom) AirTemp() (bool, string) {
	return g.GetClimateReading(grAirTemp)
}

// RH - returns the relative humidity the climate as a string
func (g *Growroom) RH() (bool, string) {
	return g.GetClimateReading(grRH)
}

// Light - returns the Light level for the climate as a string
func (g *Growroom) Light() (bool, string) {
	return g.GetClimateReading(grLight)
}

// CO2 - returns the CO2 for the climate as a string
func (g *Growroom) CO2() (bool, string) {
	return g.GetClimateReading(grCO2)
}

// EC - returns the ec for the rootzone as a string
func (g *Growroom) EC() (bool, string) {
	return g.GetRootzoneReading(grEC)
}

// PH - returns the ph for the rootzone as a string
func (g *Growroom) PH() (bool, string) {
	return g.GetRootzoneReading(grPH)
}

// WaterTemp - returns the water temp for the rootzone as a string
func (g *Growroom) WaterTemp() (bool, string) {
	return g.GetRootzoneReading(grTemp)
}

// GetClimateHistory - returns the history for the growroom
func (g *Growroom) GetClimateHistory(from, to time.Time, points int) error {
	if len(g.Climates()) > 0 {
		return g.Climates()[0].GetHistory(from, to, points)
	}
	return fmt.Errorf("Growroom has no Intelliclimates")
}

// GetDoserHistory - returns the history for the growroom
func (g *Growroom) GetDoserHistory(from, to time.Time, points int) error {
	if len(g.Dosers()) > 0 {
		return g.Dosers()[0].GetHistory(from, to, points)
	}
	return fmt.Errorf("Growroom has no Intellidosers")
}
