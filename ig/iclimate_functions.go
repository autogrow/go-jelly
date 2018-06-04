package ig

import "fmt"

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
