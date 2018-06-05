package ig

import (
	"fmt"
	"sync"
	"time"

	"github.com/autogrow/go-jelly/ig/datastructs"
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
	ValidConfig bool                      `json:"valid_config"`
	Config      *datastructs.ConfigIDose  `json:"config"`
	Metrics     *datastructs.MetricsIDose `json:"metrics"`
	ValidStatus bool                      `json:"valid_status"`
	Status      *datastructs.StatusIDose  `json:"status"`
	History     *datastructs.DoserHistory `json:"history"`
	tx          *transaction
}

// NewIntelliDose - returns a new intellidose for the device passed in
func NewIntelliDose(dev *Device) *IntelliDose {
	return &IntelliDose{
		dev,
		false,
		&datastructs.ConfigIDose{},
		&datastructs.MetricsIDose{},
		false,
		&datastructs.StatusIDose{},
		&datastructs.DoserHistory{},
		&transaction{new(sync.Mutex), false},
	}
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

// GetAll will get the config, state and metrics from the API
func (id *IntelliDose) GetAll() error {
	if err := id.GetMetrics(); err != nil {
		return err
	}

	if err := id.GetConfig(); err != nil {
		return err
	}

	if err := id.GetState(); err != nil {
		return err
	}

	return nil
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
func (id *IntelliDose) GetHistory(to, from time.Time, points int) error {
	msi, err := getHistory(id.client, id.GetID(), to, from, points)
	if err != nil {
		return err
	}
	return updateStruct(msi, id.History)
}

// SaveConfigState will save the config and state
func (id *IntelliDose) SaveConfigState() error {
	return id.client.SaveDevice(id)
}

// StatePayload builds and returns the state payload for updating a devices state or config
func (id *IntelliDose) StatePayload() (interface{}, error) {
	if !id.ValidConfig {
		return nil, fmt.Errorf("invalid config")
	}

	if !id.ValidStatus {
		return nil, fmt.Errorf("invalid status")
	}

	msi := make(map[string]interface{})
	msi["device"] = id.GetID()
	msi["state"] = id.Status
	msi["config"] = id.Config
	return msi, nil
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
