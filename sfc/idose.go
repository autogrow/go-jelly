package sfc

import (
	"encoding/json"
	"time"
)

type empty struct{}

// NewIntelliDose returns a new IntelliDose with the given serial number
func NewIntelliDose(sn string) *IntelliDose {
	return &IntelliDose{iDoseShadow: &iDoseShadow{}, serial: sn}
}

// IntelliDose represents the IntelliDose single function controller
type IntelliDose struct {
	*iDoseShadow
	updateChan  chan empty
	LastUpdated time.Time
	serial      string
}

// Serial returns the serial number of this IntelliDose
func (id *IntelliDose) Serial() string {
	return id.serial
}

func (id *IntelliDose) WaitForUpdate() {
	last := id.LastUpdated
	for {
		time.Sleep(time.Second / 10)
		if last != id.LastUpdated {
			return
		}
	}
}

// Update the IntelliDose from the given JSON payload, WaitForUpdate will stop
// blocking after this
func (id *IntelliDose) Update(b []byte) error {
	err := json.Unmarshal(b, &id)
	if err != nil {
		return err
	}

	id.LastUpdated = time.Now()

	return nil
}

// Readings returns the current readings for the IntelliDose
func (id *IntelliDose) Readings() MetricsIDose {
	return id.State.Reported.Metrics
}

// Config returns the configuration of the IntelliDose
func (id *IntelliDose) Config() ConfigIDose {
	return id.State.Reported.Config
}

// Settings returns the settings of the IntelliDose
func (id *IntelliDose) Settings() SettingsIDose {
	return id.State.Reported.Settings
}

// IsDayTime returns true if the IntelliDose thinks that it is currently day time
func (id *IntelliDose) IsDayTime() bool {
	_now := time.Now()
	_ds, _ := time.Parse("", id.State.Reported.Config.Times.DayStart)
	_de, _ := time.Parse("", id.State.Reported.Config.Times.DayEnd)

	now := (_now.Hour() * 100) + _now.Minute()
	ds := (_ds.Hour() * 100) + _ds.Minute()
	de := (_de.Hour() * 100) + _de.Minute()

	if now > ds && now < de {
		return true
	}

	return false
}
