package ig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ValidResponse - checks the map[string]interface{} contains information for the given device, also needs to contain a Last Updated time
func validResponse(msi map[string]interface{}, devType string) (map[string]interface{}, float64, error) {
	// Check that repsonce contains an iclimate readings field
	rawResponse, exist := msi[devType]

	if !exist {
		return nil, 0, fmt.Errorf("Data doesn't contain any %s readings", devType)
	}

	// Convert Raw Readings to a map
	response, valid := rawResponse.(map[string]interface{})
	if !valid {
		return nil, 0, fmt.Errorf("readings is not a map[string]interface{}")
	}

	// Check that repsonse contains a last updated field
	rawLastUpdated, exist := msi["last_updated"]

	if !exist {
		return nil, 0, fmt.Errorf("Data doesn't contain a last updated time")
	}
	last := rawLastUpdated.(float64) / 1000

	return response, last, nil
}

// UpdateStruct - converts a IC reading maps to a struct as well as lastupdated and
func updateStruct(src map[string]interface{}, target interface{}) error {
	jBytes, err := json.Marshal(src)
	if err != nil {
		return fmt.Errorf("Error marshalling src: %s", err)
	}

	err = json.Unmarshal(jBytes, target)
	if err != nil {
		return fmt.Errorf("Error unmarshalling src: %s", err)
	}
	return nil
}

// Get - Returns a map[string]interface{} and error for a specified endpoint for a device
func (c *Client) get(endpoint string) (map[string]interface{}, error) {
	// Make Request
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		// handle err
		return nil, fmt.Errorf("Unable to get build request for endpoint %s: %s", endpoint, err)
	}

	// Set Request Headers to include client token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.getToken())

	// Do the request
	resp, err := c.Do(req)

	if err != nil {
		// handle err
		return nil, fmt.Errorf("Get request return an error for endpoint %s: %s", endpoint, err)
	}

	// Process response, this should contain a last_updated and iclimate fields
	defer resp.Body.Close()

	msi := make(map[string]interface{})

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %s", err)
	}

	err = json.Unmarshal(data, &msi)

	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response body: %s", err)
	}

	return msi, nil
}

// GetHistory - Returns a map[string]interface{} and error for a specified endpoint for a device
func getHistory(c *Client, device string, from, to time.Time, points int) (map[string]interface{}, error) {
	// Build URL
	startTStamp := from.Unix() * 1000
	endTStamp := to.Unix() * 1000

	url := fmt.Sprintf("%s%s&points=%d&to_date=%d&from_date=%d", igHistoryURI, device, points, startTStamp, endTStamp)

	// Make Request
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		// handle err
		return nil, fmt.Errorf("Unable to get build request for dev %s: %s", device, err)
	}

	// Set Request Headers to include client token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.getToken())

	// Do the request
	resp, err := c.Do(req)

	if err != nil {
		// handle err
		return nil, fmt.Errorf("Get request return an error for device %s: %s", device, err)
	}

	// Process response, this should contain a last_updated and iclimate fields
	defer resp.Body.Close()

	msi := make(map[string]interface{})

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("couldn't read response body: %s", err)
	}

	err = json.Unmarshal(data, &msi)

	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response body: %s", err)
	}

	// Check that repsonce contains an iclimate readings field
	responseDevice, exist := msi["device"].(string)

	if !exist {
		return nil, fmt.Errorf("Response doesn't reference any device")
	}

	if responseDevice != device {
		return nil, fmt.Errorf("history is not for me - this shouldn't be possible")
	}

	// Convert Raw Readings to a map
	response, valid := msi["history"].(map[string]interface{})
	if !valid {
		return nil, fmt.Errorf("history doesn't exists or is not a map[string]interface{}")
	}

	return response, nil
}

// Put - takes the state passed to it and pushes it to Intelligrow
func (c *Client) put(endpoint string, payload map[string]interface{}) error {
	jsonValue, _ := json.Marshal(payload)

	// fmt.Println(string(jsonValue))

	body := bytes.NewBuffer(jsonValue)

	req, err := http.NewRequest("PUT", endpoint, body)

	if err != nil {
		// handle err
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.getToken())

	resp, err := c.Do(req)

	if err != nil {
		// handle err
		return fmt.Errorf("Put request return an error for endpoint %s: %s", endpoint, err)
	}

	// fmt.Println("HTTP Response Status: " + strconv.Itoa(resp.StatusCode))
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		return fmt.Errorf("Bad HTTP response %d from Put request to endpoint %s", resp.StatusCode, endpoint)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("couldn't read response body: %s", err)
	}

	respMap := make(map[string]interface{})

	err = json.Unmarshal(data, &respMap)

	if err != nil {
		return fmt.Errorf("couldn't unmarshal response body: %s", err)
	}

	return nil
}
