package ig

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	// IDose - type reported by a device if it is an IntelliDose
	IDose = "idose"
	// NutrientDosingFunction - string used to identify the status field for nutrient dosing
	NutrientDosingFunction = "Nutrient Dosing"
	// PHDosingFunction - string used to identify the status field for ph dosing
	PHDosingFunction = "ph"
	// IrrigationFunction - string used to identify the status field for irrigation
	IrrigationFunction = "irrigation"
	// StationFunction - string used to identify the status field for irrigation station 1
	StationFunction = "Irrigation Station "
)

// IntelliDose - IntelliDose object
type IntelliDose struct {
	*Device     `json:"device"`
	ValidConfig bool          `json:"valid_config"`
	Config      *ConfigIDose  `json:"config"`
	Metrics     *MetricsIDose `json:"metrics"`
	ValidStatus bool          `json:"valid_status"`
	Status      *StatusIDose  `json:"status"`
	History     *DoserHistory `json:"history"`
}

// NewIntelliDose - returns a new intellidose for the device passed in
func NewIntelliDose(dev *Device) *IntelliDose {
	return &IntelliDose{
		dev,
		false,
		&ConfigIDose{},
		&MetricsIDose{},
		false,
		&StatusIDose{},
		&DoserHistory{},
	}
}

// MarshalJSON will marshal the IntelliDose into JSON
func (id *IntelliDose) MarshalJSON() ([]byte, error) {
	return json.Marshal(id)
}

// SetPHTarget will set the target pH the system should dose to
func (id *IntelliDose) SetPHTarget(target float64) error {
	return fmt.Errorf("not implemented")
}

// SetNutrientTarget will set the target EC the system should dose to
func (id *IntelliDose) SetNutrientTarget(target float64) error {
	return fmt.Errorf("not implemented")
}

// DisableNutrientDosing will disable the nutrient dosing
func (id *IntelliDose) DisableNutrientDosing(target float64) error {
	return fmt.Errorf("not implemented")
}

// EnableNutrientDosing will enable the nutrient dosing
func (id *IntelliDose) EnableNutrientDosing(target float64) error {
	return fmt.Errorf("not implemented")
}

// DisablePHDosing will disable the pH dosing
func (id *IntelliDose) DisablePHDosing(target float64) error {
	return fmt.Errorf("not implemented")
}

// EnablePHDosing will enable the pH dosing
func (id *IntelliDose) EnablePHDosing(target float64) error {
	return fmt.Errorf("not implemented")
}

// getEndpoint the device by quering the endpoint passed in
func (id *IntelliDose) getEndpoint(endpoint string) error {
	switch endpoint {
	case MetricsEP:
		url := igMetricsURI + id.GetID()
		msi, err := id.client.get(url)
		if err != nil {
			return err
		}
		id.Readings = msi
		return updateStruct(msi, id.Metrics)

	case ConfigEP:
		url := igConfigURI + id.GetID()
		msi, err := id.client.get(url)
		if err != nil {
			return err
		}
		return updateStruct(msi, id.Config)

	case StateEP:
		url := igDevicesURI + id.GetID()
		msi, err := id.client.get(url)
		if err != nil {
			return err
		}
		return updateStruct(msi, id.Status)

	default:
		return fmt.Errorf("Unknown field %s requested to be updated", endpoint)
	}
}

// GetMetrics the device by quering the endpoint passed in
func (id *IntelliDose) GetMetrics() error {
	endpoint := igMetricsURI + id.GetID()
	msi, err := id.client.get(endpoint)
	if err != nil {
		return err
	}

	metrics, updated, err := validResponse(msi, id.Type)

	if err != nil {
		return err
	}

	id.LastUpdated = updated
	metrics["ec"] = metrics["ec"].(float64) / 100.0
	id.Readings = metrics

	return updateStruct(metrics, id.Metrics)
}

// GetConfig - this pulls both the config and state from the device endpoint
func (id *IntelliDose) GetConfig() error {
	endpoint := igConfigURI + id.GetID()

	response, err := id.client.get(endpoint)
	if err != nil {
		id.ValidConfig = false
		return err
	}

	cfg, _, err := validResponse(response, id.Type)

	err = updateStruct(cfg, id.Config)
	if err != nil {
		id.ValidConfig = false
		return err
	}
	id.ValidConfig = true

	return nil
}

// GetState - this pulls both the state from the device endpoint
func (id *IntelliDose) GetState() error {
	endpoint := igDeviceStateURI + id.GetID()
	msi, err := id.client.get(endpoint)
	if err != nil {
		id.ValidStatus = false
		return err
	}

	// Check that repsonce contains an iclimate readings field
	rawResponse, exist := msi[id.Type]

	if !exist {
		return fmt.Errorf("Data doesn't contain any %s readings", id.Type)
	}

	// Convert Raw Readings to a map
	response, valid := rawResponse.(map[string]interface{})
	if !valid {
		return fmt.Errorf("readings is not a map[string]interface{}")
	}

	err = updateStruct(response, id.Status)
	if err != nil {
		id.ValidStatus = false
		return err
	}
	id.ValidStatus = true

	return nil
}

// GetConfigState - this pulls both the config and state from the device endpoint
func (id *IntelliDose) GetConfigState() error {
	err := id.GetConfig()
	if err != nil {
		return err
	}

	err = id.GetState()
	if err != nil {
		return err
	}

	return nil
}

// GetHistory the device by quering the history endpont for the time period specified
func (id *IntelliDose) getHistory(to, from time.Time, points int) error {
	msi, err := getHistory(id.client, id.GetID(), to, from, points)
	if err != nil {
		return err
	}
	return updateStruct(msi, id.History)
}

// UpdateState - Push the state to the IG
func (id *IntelliDose) UpdateState() error {
	if !id.ValidConfig {
		return fmt.Errorf("Device doesn't have a valid config")
	}

	if !id.ValidStatus {
		return fmt.Errorf("Device doesn't have a valid status")
	}

	msi := make(map[string]interface{})
	msi["device"] = id.GetID()
	msi["state"] = id.Status
	msi["config"] = id.Config

	return id.client.put(igDevicesURI, msi)
}

// ForceNutrient - sets the appropriate field in the Status to force an nutrient Dose
func (id *IntelliDose) ForceNutrient() bool {
	err := id.GetConfigState()
	if err != nil {
		return false
	}

	for num, status := range id.Status.Status {
		if status.Function == NutrientDosingFunction {
			id.Status.Status[num].ForceOn = true
			return true
		}
	}
	return false
}

// ForcePH - sets the appropriate field in the Status to force an pH Dose
func (id *IntelliDose) ForcePH() bool {
	err := id.GetConfigState()
	if err != nil {
		return false
	}

	for num, status := range id.Status.Status {
		if status.Function == PHDosingFunction {
			id.Status.Status[num].ForceOn = true
			return true
		}
	}
	return false
}

// ForceIrrigation - sets the appropriate field in the Status to force an irrigation
func (id *IntelliDose) ForceIrrigation() bool {
	err := id.GetConfigState()

	if err != nil {
		return false
	}

	for num, status := range id.Status.Status {
		if status.Function == IrrigationFunction {
			id.Status.Status[num].ForceOn = true
			return true
		}
	}
	return false
}

// ForceStation - sets the appropriate field in the Status to force an irrigation on the station specified (1-4)
func (id *IntelliDose) ForceStation(stn string) bool {
	funcName := StationFunction + stn
	for num, status := range id.Status.Status {
		if status.Function == funcName {
			id.Status.Status[num].ForceOn = true
			return true
		}
	}
	return false
}

// AverageDoseReadings - returns an average for the field specified from a list of IntelliDose
func AverageDoseReadings(dosers []*IntelliDose, field string) float64 {
	var sum float64
	var validDevices int

	for _, doser := range dosers {
		if doser.IsValid() {
			sum += doser.Readings[field].(float64)
			validDevices++
		}
	}

	if validDevices != 0 {
		return sum / float64(validDevices)
	}

	return 0
}

// DoserHistory - consists of a slice of history points
type DoserHistory struct {
	Points []*DoserHistoryPoint `json:"points"`
}

// DoseMetricsHistory - Metrics the history point contains
type DoseMetricsHistory struct {
	EC   float64 `json:"ec"`
	PH   float64 `json:"pH"`
	Temp float64 `json:"nut_temp"`
}

// DoserHistoryPoint - defines a single history point reported for a IntelliDose
type DoserHistoryPoint struct {
	Timestamp float64            `json:"timestamp"`
	Status    Status             `json:"status"`
	Metrics   DoseMetricsHistory `json:"metrics"`
}

// IDose
type iDoseShadow struct {
	State StateIDose `json:"state"`
}

// StateIDose represents the State data structure from an IntelliDose packet
type StateIDose struct {
	Reported ReportedIDose `json:"reported"`
}

// ReportedIDose represents the Reported data structure from an IntelliDose packet
type ReportedIDose struct {
	Config    ConfigIDose  `json:"config"`
	Metrics   MetricsIDose `json:"metrics"`
	Status    StatusIDose  `json:"status"`
	Source    string       `json:"source"`
	Device    string       `json:"device"`
	Timestamp int64        `json:"timestamp"`
	Connected bool         `json:"connected"`
}

// ConfigIDose represents the Config data structure from an IntelliDose packet
type ConfigIDose struct {
	Units      UnitsIDose      `json:"units"`
	Times      TimesIDose      `json:"times"`
	Functions  FunctionsIDose  `json:"functions"`
	Advanced   AdvancedIDose   `json:"advanced"`
	General    GeneralIDose    `json:"general"`
	Scheduling SchedulingIDose `json:"scheduling"`
	Reminders  RemindersIDose  `json:"reminder"`
}

// MetricsIDose represents the Metrics data structure from an IntelliDose packet
type MetricsIDose struct {
	Ec      float64 `json:"ec"`
	NutTemp float64 `json:"nut_temp"`
	PH      float64 `json:"pH"`
}

// StatusIDose represents the Status data structure from an IntelliDose packet
type StatusIDose struct {
	General   GeneralStatusIDose  `json:"general"`
	Nutrient  NutrientIDose       `json:"nutrient"`
	SetPoints SetPointsIDose      `json:"set_points"`
	Status    []StatusStatusIDose `json:"status"`
}

// GeneralStatusIDose represents the GeneralStatus data structure from an IntelliDose packet
type GeneralStatusIDose struct {
	DoseInterval        byte                    `json:"dose_interval"`
	NutrientDoseTime    byte                    `json:"nutrient_dose_time"`
	WaterOnTime         byte                    `json:"water_on_time"`
	IrrigationInterval1 IrrigationIntervalIDose `json:"irrigation_interval_1"`
	IrrigationInterval2 IrrigationIntervalIDose `json:"irrigation_interval_2"`
	IrrigationInterval3 IrrigationIntervalIDose `json:"irrigation_interval_3"`
	IrrigationInterval4 IrrigationIntervalIDose `json:"irrigation_interval_4"`
	IrrigationDuration1 int                     `json:"irrigation_duration_1"`
	IrrigationDuration2 int                     `json:"irrigation_duration_2"`
	IrrigationDuration3 int                     `json:"irrigation_duration_3"`
	IrrigationDuration4 int                     `json:"irrigation_duration_4"`
	MaxNutrientDoseTime byte                    `json:"max_nutrient_dose_time"`
	MaxPhDoseTime       byte                    `json:"max_ph_dose_time"`
	Mix1                byte                    `json:"mix_1"`
	Mix2                byte                    `json:"mix_2"`
	Mix3                byte                    `json:"mix_3"`
	Mix4                byte                    `json:"mix_4"`
	Mix5                byte                    `json:"mix_5"`
	Mix6                byte                    `json:"mix_6"`
	Mix7                byte                    `json:"mix_7"`
	Mix8                byte                    `json:"mix_8"`
	PhDoseTime          byte                    `json:"ph_dose_time"`
}

// IrrigationIntervalIDose represents the IrrigationInterval data structure from an IntelliDose packet
type IrrigationIntervalIDose struct {
	Day   int `json:"day"`
	Night int `json:"night"`
	Every int `json:"every"`
}

// NutrientIDose represents the Nutrient data structure from an IntelliDose packet
type NutrientIDose struct {
	Detent  byte         `json:"detent"`
	Ec      EcIDose      `json:"ec"`
	NutTemp NutTempIDose `json:"nut_temp"`
	Ph      PhIDose      `json:"ph"`
}

// EcIDose represents the Ec data structure from an IntelliDose packet
type EcIDose struct {
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

// NutTempIDose represents the NutTemp data structure from an IntelliDose packet
type NutTempIDose struct {
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

// PhIDose represents the Ph data structure from an IntelliDose packet
type PhIDose struct {
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

// SetPointsIDose represents the SetPoints data structure from an IntelliDose packet
type SetPointsIDose struct {
	Nutrient      float64 `json:"nutrient"`
	NutrientNight float64 `json:"nutrient_night"`
	PhDosing      string  `json:"ph_dosing"`
	Ph            float64 `json:"ph"`
}

// StatusStatusIDose represents the StatusStatus data structure from an IntelliDose packet
type StatusStatusIDose struct {
	Active   bool   `json:"active"`
	Enabled  bool   `json:"enabled"`
	ForceOn  bool   `json:"force_on"`
	Function string `json:"function"`
}

// UnitsIDose represents the Units data structure from an IntelliDose packet
type UnitsIDose struct {
	DateFormat              string `json:"date_format"`
	Temperature             string `json:"temperature"`
	Ec                      string `json:"ec"`
	TdsConversationStandart int    `json:"tds_conversation_standart"`
}

// TimesIDose represents the Times data structure from an IntelliDose packet
type TimesIDose struct {
	DayStart int `json:"day_start"`
	DayEnd   int `json:"day_end"`
}

// FunctionsIDose represents the Functions data structure from an IntelliDose packet
type FunctionsIDose struct {
	NutrientsParts     byte   `json:"nutrients_parts"`
	PhDosing           string `json:"ph_dosing"`
	IrrigationMode     string `json:"irrigation_mode"`
	IrrigationStations byte   `json:"irrigation_stations"`
	SeparatePumpOutput bool   `json:"separate_pump_output"`
	UseWater           bool   `json:"use_water"`
	ExternalAlarm      bool   `json:"external_alarm"`
	DayNightEc         bool   `json:"day_night_ec"`
	IrrigationStation1 string `json:"irrigation_station_1"`
	IrrigationStation2 string `json:"irrigation_station_2"`
	IrrigationStation3 string `json:"irrigation_station_3"`
	IrrigationStation4 string `json:"irrigation_station_4"`
	Scheduling         bool   `json:"scheduling"`
	MuteBuzzer         bool   `json:"mute_buzzer"`
}

// AdvancedIDose represents the Advanced data structure from an IntelliDose packet
type AdvancedIDose struct {
	ProportinalDosing bool   `json:"proportinal_dosing"`
	SequentialDosing  bool   `json:"sequential_dosing"`
	DisableEc         bool   `json:"disable_ec"`
	DisablePh         bool   `json:"disable_ph"`
	MntnReminderFreq  string `json:"mntn_reminder_freq"`
}

// GeneralIDose represents the General data structure from an IntelliDose packet
type GeneralIDose struct {
	DeviceName string  `json:"device_name"`
	Firmware   float64 `json:"firmware"`
	Growroom   string  `json:"growroom"`
}

// SchedulingIDose represents the scheduling data structure from an IntelliDose packet
type SchedulingIDose struct {
	LastUpdated float64 `json:"last_updated"`
	Mode        string  `json:"mode"`
}

// ReminderIDose represents the a single reminder data structure from an IntelliDose packet
type ReminderIDose struct {
	CleanECProbe float64 `json:"clean_ec_probe"`
	CleanpHProbe float64 `json:"clean_ph_electrode"`
	CheckECProbe float64 `json:"check_ec_probe"`
	CalibratePH  float64 `json:"calibrate_ph"`
	CleanFilters float64 `json:"clean_filters"`
}

// RemindersIDose represents the reminders data structure from an IntelliDose packet
type RemindersIDose struct {
	Frequency    string        `json:"frequency"`
	StartDate    float64       `json:"start_date"`
	ReminderList ReminderIDose `json:"reminder_list"`
}
