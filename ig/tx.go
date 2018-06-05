package ig

import "sync"

type transaction struct {
	lock    *sync.Mutex
	running bool
}

type configGetterSaver interface {
	GetConfigState() error
	SaveConfigState() error
}

func (tx *transaction) guard(cgs configGetterSaver, runner func()) error {
	if !tx.running {
		if err := cgs.GetConfigState(); err != nil {
			return err
		}
	}
	runner()
	if !tx.running {
		return cgs.SaveConfigState()
	}
	return nil
}

// Transaction allows multiple changes to be modified and pushed in one API request
// instead of doing multiple API requests as happens when those methods are called
// outside of a transation.
//
//     err := ic.Transaction(func() error) {
//       ic.SetCO2Target(2500)
//       ic.EnableCO2Dosing()
//       return nil
//     })
//
// When the methods are called inside the transaction, the changes will only by sent
// at the end of the callback and if the callback doesn't return an error.
//
// Using the transaction will also pull down the config and state immediately prior to
// and push it up immediately after, making the changes.  This small window helps to
// ensure that changes from other parties accessing the API are not overwritten.
func (ic *IntelliClimate) Transaction(runner func() error) error {
	ic.tx.lock.Lock()
	ic.tx.running = true
	defer ic.tx.lock.Unlock()
	defer func() { ic.tx.running = false }()

	if err := ic.GetConfigState(); err != nil {
		return err
	}

	if err := runner(); err != nil {
		return err
	}

	return ic.client.SaveDevice(ic)
}

// Transaction allows multiple changes to be modified and pushed in one API request
// instead of doing multiple API requests as happens when those methods are called
// outside of a transation.
//
//     err := id.Transaction(func() error) {
//       id.ForceIrrigation()
//       id.ForcePHDose()
//       id.ForceECDose()
//       return nil
//     })
//
// When the methods are called inside the transaction, the changes will only by sent
// at the end of the callback and if the callback doesn't return an error.
//
// Using the transaction will also pull down the config and state immediately prior to
// and push it up immediately after, making the changes.  This small window helps to
// ensure that changes from other parties accessing the API are not overwritten.
func (id *IntelliDose) Transaction(runner func() error) error {
	id.tx.lock.Lock()
	id.tx.running = true
	defer id.tx.lock.Unlock()
	defer func() { id.tx.running = false }()

	if err := id.GetConfigState(); err != nil {
		return err
	}

	if err := runner(); err != nil {
		return err
	}

	return id.client.SaveDevice(id)
}
