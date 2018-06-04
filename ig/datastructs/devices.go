package datastructs

// DeviceStatus - a generic layout of how functions are reported from the history endpoint
type DeviceStatus struct {
	Active    bool   `json:"active"`
	Enabled   bool   `json:"enabled"`
	ForceOn   bool   `json:"force_on"`
	Function  string `json:"function"`
	Installed bool   `json:"installed"`
}

// Status - the status from the histpry endpoint is a collect of DeviceStatus'
type Status struct {
	Status []DeviceStatus `json:"status"`
}
