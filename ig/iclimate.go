package ig

import (
	"fmt"
	"time"

	"github.com/autogrow/go-jelly/ig/datastructs"
)

const (
	// IClimate - type reported by a device if it is an IntelliClimate
	IClimate = "iclimate"
)

// IntelliClimate - Intelliclimate object
type IntelliClimate struct {
	*Device     `json:"device"`
	ValidConfig bool                         `json:"valid_config"`
	Config      *datastructs.ConfigIClimate  `json:"config"`
	Metrics     *datastructs.MetricsIClimate `json:"metrics"`
	ValidStatus bool                         `json:"valid_status"`
	Status      *datastructs.StatusIClimate  `json:"status"`
	History     *datastructs.ClimateHistory  `json:"history"`
}

// NewIntelliClimate - returns a new intelliclimate for the device passed in
func NewIntelliClimate(dev *Device) *IntelliClimate {
	return &IntelliClimate{
		dev,
		false,
		&datastructs.ConfigIClimate{},
		&datastructs.MetricsIClimate{},
		false,
		&datastructs.StatusIClimate{},
		&datastructs.ClimateHistory{},
	}
}

// SetTempTarget will set the temperature that the room should be kept to
func (ic *IntelliClimate) SetTempTarget(target float64) error {
	return fmt.Errorf("not implemented")
}

// SetCO2Target will set the CO2 levels in PPM that the room should be kept to
func (ic *IntelliClimate) SetCO2Target(target float64) error {
	return fmt.Errorf("not implemented")
}

// SetRHTarget will set the RH target that the room should be kept to
func (ic *IntelliClimate) SetRHTarget(target float64) error {
	return fmt.Errorf("not implemented")
}

// EnableCO2Dosing will enable the CO2 dosing
func (ic *IntelliClimate) EnableCO2Dosing() error {
	return fmt.Errorf("not implemented")
}

// DisableCO2Dosing will disable the CO2 dosing
func (ic *IntelliClimate) DisableCO2Dosing() error {
	return fmt.Errorf("not implemented")
}

// getEndpoint the device by quering the endpoint passed in
func (ic *IntelliClimate) getEndpoint(endpoint string) error {
	switch endpoint {
	case MetricsEP:
		url := igMetricsURI + ic.GetID()
		msi, err := ic.client.get(url)
		if err != nil {
			return err
		}
		ic.Readings = msi
		return updateStruct(msi, ic.Metrics)

	case ConfigEP:
		url := igConfigURI + ic.GetID()
		msi, err := ic.client.get(url)
		if err != nil {
			return err
		}
		return updateStruct(msi, ic.Config)

	case StateEP:
		url := igDevicesURI + ic.GetID()
		msi, err := ic.client.get(url)
		if err != nil {
			return err
		}
		return updateStruct(msi, ic.Status)

	default:
		return fmt.Errorf("Unknown field %s requested to be updated", endpoint)
	}
}

// GetAll will get the config, state and metrics from the API
func (ic *IntelliClimate) GetAll() error {
	if err := ic.GetMetrics(); err != nil {
		return err
	}

	if err := ic.GetConfig(); err != nil {
		return err
	}

	if err := ic.GetState(); err != nil {
		return err
	}

	return nil
}

// GetMetrics the device by quering the endpoint passed in
func (ic *IntelliClimate) GetMetrics() error {
	endpoint := igMetricsURI + ic.GetID()
	msi, err := ic.client.get(endpoint)
	if err != nil {
		return err
	}

	metrics, updated, err := validResponse(msi, ic.Type)

	if err != nil {
		return err
	}

	ic.LastUpdated = updated
	ic.Readings = metrics

	return updateStruct(metrics, ic.Metrics)
}

// GetConfig - this pulls both the config and state from the device endpoint
func (ic *IntelliClimate) GetConfig() error {
	endpoint := igConfigURI + ic.GetID()

	response, err := ic.client.get(endpoint)
	if err != nil {
		ic.ValidConfig = false
		return err
	}

	cfg, _, err := validResponse(response, ic.Type)

	err = updateStruct(cfg, ic.Config)
	if err != nil {
		ic.ValidConfig = false
		return err
	}
	ic.ValidConfig = true

	return nil
}

// GetState - this pulls both the state from the device endpoint
func (ic *IntelliClimate) GetState() error {
	endpoint := igDeviceStateURI + ic.GetID()

	msi, err := ic.client.get(endpoint)
	if err != nil {
		ic.ValidStatus = false
		return err
	}

	// Check that repsonce contains an iclimate readings field
	rawResponse, exist := msi[ic.Type]

	if !exist {
		return fmt.Errorf("Data doesn't contain any %s readings", ic.Type)
	}

	// Convert Raw Readings to a map
	response, valid := rawResponse.(map[string]interface{})
	if !valid {
		return fmt.Errorf("readings is not a map[string]interface{}")
	}

	err = updateStruct(response, ic.Status)
	if err != nil {
		ic.ValidStatus = false
		return err
	}
	ic.ValidStatus = true

	return nil
}

// GetConfigState - this pulls both the config and state from the device endpoint
func (ic *IntelliClimate) GetConfigState() error {
	cfgError := ic.GetConfig()
	if cfgError != nil {

	}

	stateError := ic.GetState()
	if stateError != nil {

	}

	return nil
}

// GetHistory the device by quering the history endpont for the time period specified
func (ic *IntelliClimate) GetHistory(to, from time.Time, points int) error {
	msi, err := getHistory(ic.client, ic.GetID(), to, from, points)
	if err != nil {
		return err
	}
	return updateStruct(msi, ic.History)
}

// UpdateState - Push the state to the IG
func (ic *IntelliClimate) UpdateState() error {
	if !ic.ValidConfig {
		return fmt.Errorf("Device doesn't have a valid config")
	}

	if !ic.ValidStatus {
		return fmt.Errorf("Device doesn't have a valid status")
	}

	msi := make(map[string]interface{})
	msi["device"] = ic.GetID()
	msi["state"] = ic.Status
	msi["config"] = ic.Config

	return ic.client.put(StateEP, msi)
}

// AverageClimateReadings - returns an average for the field specified from a list of IntelliDose
func AverageClimateReadings(climates []*IntelliClimate, field string) float64 {
	var sum float64
	var validDevices int

	for _, climate := range climates {
		if climate.IsValid() {
			sum += climate.Readings[field].(float64)
			validDevices++
		}
	}

	if validDevices != 0 {
		return sum / float64(validDevices)
	}

	return 0
}
