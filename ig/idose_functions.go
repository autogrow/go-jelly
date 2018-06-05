package ig

import "fmt"

// ForceNutrientDose will force a nutrient dose on the controller
func (id *IntelliDose) ForceNutrientDose() error {
	return id.tx.guard(id, func() {
		for num, status := range id.Status.Status {
			if status.Function == NutrientDosingFunction {
				id.Status.Status[num].ForceOn = true
			}
		}
	})
}

// ForcePHDose will force a pH dose on the controller
func (id *IntelliDose) ForcePHDose() error {
	return id.tx.guard(id, func() {
		for num, status := range id.Status.Status {
			if status.Function == PHDosingFunction {
				id.Status.Status[num].ForceOn = true
			}
		}
	})
}

// ForceIrrigation will force an irrigation on the controller
func (id *IntelliDose) ForceIrrigation() error {
	return id.tx.guard(id, func() {
		for num, status := range id.Status.Status {
			if status.Function == IrrigationFunction {
				id.Status.Status[num].ForceOn = true
			}
		}
	})
}

// ForceStation will force an irrigation on the station specified (1-4)
func (id *IntelliDose) ForceStation(stn string) error {
	return id.tx.guard(id, func() {
		funcName := StationFunction + stn
		for num, status := range id.Status.Status {
			if status.Function == funcName {
				id.Status.Status[num].ForceOn = true
			}
		}
	})
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
