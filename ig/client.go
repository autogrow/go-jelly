package ig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	igBaseURL = "https://api.autogrow.com/v1"

	igTokenURI       = igBaseURL + "/auth/token"
	igRefreshURI     = igBaseURL + "/auth/token/refresh"
	igDeviceURI      = igBaseURL + "/intelligrow/devices?username="
	igDevicesURI     = igBaseURL + "/intelligrow/devices"
	igDeviceStateURI = igBaseURL + "/intelligrow/devices/state?device="
	igMetricsURI     = igBaseURL + "/intelligrow/devices/metrics?device="
	igConfigURI      = igBaseURL + "/intelligrow/devices/config?device="
	igHistoryURI     = igBaseURL + "/intelligrow/devices/history?device="
)

// Client - object that can be used to communicate directly with intelligrow
type Client struct {
	*http.Client
	lock               *sync.RWMutex
	username           string
	password           string
	token              string
	refreshToken       string
	refreshTime        float64
	CheckedDevices     float64
	growrooms          map[string]*Growroom
	devices            *Devices
	tokenRefresherQuit chan bool
}

// NewClient - Return a new IGClient object as well as any errors detected
// The function takes three inputs
// First - Intelligrow username as a string
// Second - Intelligrow Password as a string
// The Function will return a new client object with the appropriate tokens as well as any error encountered.
func NewClient(user, pass string) (*Client, error) {
	c := &Client{
		Client:    &http.Client{Timeout: time.Second * 30},
		lock:      new(sync.RWMutex),
		username:  user,
		password:  pass,
		growrooms: make(map[string]*Growroom),
	}

	// Initialize the devices object in the structure, this is blank object
	c.devices = NewDevices()

	err := c.authenticate()
	if err != nil {
		return c, err
	}

	go c.refresher()

	return c, nil
}

// Close the client (read: stop trying to refresh the auth token every hour)
func (c *Client) Close() error {
	if c.tokenRefresherQuit != nil {
		close(c.tokenRefresherQuit)
	}
	return nil
}

// GetToken - updates the token using an existing client
func (c *Client) authenticate() error {
	values := map[string]string{"username": c.username, "password": c.password}

	jsonValue, _ := json.Marshal(values)

	body := bytes.NewBuffer(jsonValue)

	req, err := http.NewRequest("POST", igTokenURI, body)

	if err != nil {
		// handle err
		return fmt.Errorf("Unable to get tokens %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		// handle err
		return fmt.Errorf("Unable to get tokens %s", err)
	}

	defer resp.Body.Close()

	return c.processBody(resp)
}

// Token returns the current token for the client
func (c *Client) getToken() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.token
}

func (c *Client) processBody(resp *http.Response) error {
	msi := make(map[string]interface{})

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("couldn't read response body: %s", err)
	}

	err = json.Unmarshal(data, &msi)

	if err != nil {
		return fmt.Errorf("couldn't unmarshal response body: %s", err)
	}

	_, exists := msi["api_access_token"]
	if !exists {
		return fmt.Errorf("Response from server doesn't contain an api token")
	}

	token, isString := msi["api_access_token"].(string)
	if !isString {
		return fmt.Errorf("Recieved Token is no a string")
	}

	_, exists = msi["expires_in"]
	if !exists {
		return fmt.Errorf("Response from server doesn't contain a refresh time")
	}

	expiresIn, isFloat := msi["expires_in"].(float64)
	if !isFloat {
		return fmt.Errorf("Refresh time couldn't be converted to a int")
	}

	_, exists = msi["refresh_token"]
	if !exists {
		return fmt.Errorf("Response from server doesn't contain a refresh token")
	}
	refresh, isString := msi["refresh_token"].(string)
	if !isString {
		return fmt.Errorf("Refresh token is not a string")
	}

	c.token = token
	c.refreshTime = expiresIn
	c.refreshToken = refresh

	return nil
}

func (c *Client) refresher() {
	c.tokenRefresherQuit = make(chan bool)
	defer func() { c.tokenRefresherQuit = nil }()
	timer := time.NewTimer(c.getRefreshTime())

	for {
		select {
		case <-c.tokenRefresherQuit:
			return
		case <-timer.C:

			body := strings.NewReader(`{"username": c.username, "refresh_token": c.refreshToken}`)
			req, err := http.NewRequest("POST", igRefreshURI, body)
			if err != nil {
				// handle err

			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := c.Do(req)
			if err != nil {
				// handle err
			}

			err = c.processBody(resp)

			resp.Body.Close()

			if err != nil {
				// handle err
			}
			timer.Reset(c.getRefreshTime())
		}
	}
}

// AutoUpdater - automatically updates the device connected to the client
func (c *Client) AutoUpdater(pollInterval int, quit chan bool, updateInterval chan int) {
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)

	for {
		select {
		case <-ticker.C:
			c.UpdateAllGrowrooms()
		case stop, ok := <-quit:
			if !ok || stop {
				return
			}
		case newInt, ok := <-updateInterval:
			if !ok {
				return
			}
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(newInt) * time.Second)
		}
	}
}

func (c *Client) getRefreshTime() time.Duration {
	rTime := c.refreshTime
	if rTime > 60 {
		rTime -= 60
	}
	return time.Duration(rTime)
}

// GetDevices is deprecated in favour of RefreshDevices
func (c *Client) GetDevices() error {
	return c.RefreshDevices()
}

// RefreshDevices will get the latest data from the API and update all known structs
func (c *Client) RefreshDevices() error {

	url := igDeviceURI + c.username

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return fmt.Errorf("Unable to get devices %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.token)

	resp, err := c.Do(req)

	if err != nil {
		return fmt.Errorf("Get request return an error %s", err)
	}

	defer resp.Body.Close()

	msi := make(map[string]interface{})

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("couldn't read response body: %s", err)
	}

	err = json.Unmarshal(data, &msi)

	if err != nil {
		return fmt.Errorf("couldn't unmarshal response body: %s", err)
	}

	checkedDevices, exists := msi["checked_devices"]

	if !exists {
		return fmt.Errorf("no checked devices in response")
	}

	c.CheckedDevices = checkedDevices.(float64)

	devices, exists := msi["devices"]

	if !exists {
		return fmt.Errorf("no devices in response")
	}

	var igDevices []*Device

	jBytes, err := json.Marshal(devices)
	if err != nil {
		return fmt.Errorf("Error marshalling devices interface: %s", err)
	}

	err = json.Unmarshal(jBytes, &igDevices)
	if err != nil {
		return fmt.Errorf("Error unmarshalling devices: %s", err)
	}

	c.lock.Lock()
	for _, d := range igDevices {
		d.AttachClient(c)
		c.devices.Add(d)
		c.addDeviceToGrowroom(d)
	}
	c.lock.Unlock()
	return err
}

func (c *Client) addDeviceToGrowroom(dev *Device) {
	grName := dev.GetGrowroom()
	gr, exists := c.growrooms[grName]
	if !exists {
		gr = NewGrowroom(grName)
		c.growrooms[grName] = gr
	}
	gr.AddDevice(dev)
}

// IntelliClimate will return the IntelliClimate with the given name or serial
func (c *Client) IntelliClimate(nameOrID string) (*IntelliClimate, error) {
	if c.devices.IsEmpty() {
		c.RefreshDevices()
	}

	ic, err := c.devices.GetClimateByID(nameOrID)
	if err != nil {
		ic, err = c.devices.GetClimateByName(nameOrID)
		if err != nil {
			return nil, fmt.Errorf("no IntelliClimate found with name or serial of %s", nameOrID)
		}
	}

	return ic, nil
}

// IntelliDose will return the IntelliDose with the given name  or serial
func (c *Client) IntelliDose(nameOrID string) (*IntelliDose, error) {
	if c.devices.IsEmpty() {
		c.RefreshDevices()
	}

	id, err := c.devices.GetDoserByID(nameOrID)
	if err != nil {
		id, err = c.devices.GetDoserByName(nameOrID)
		if err != nil {
			return nil, fmt.Errorf("no IntelliDose found with name or serial of %s", nameOrID)
		}
	}

	return id, nil
}

// IntelliDoses will return all known IntelliDoses
func (c *Client) IntelliDoses() ([]*IntelliDose, error) {
	if c.devices.IsEmpty() {
		c.RefreshDevices()
	}
	return c.devices.Dosers(), nil
}

// IntelliClimates will return all known IntelliClimates
func (c *Client) IntelliClimates() ([]*IntelliClimate, error) {
	if c.devices.IsEmpty() {
		c.RefreshDevices()
	}
	return c.devices.Climates(), nil
}

// Devices returns the devices found in IntelliGrow for the user
func (c *Client) Devices() []*Device {
	if c.devices.IsEmpty() {
		c.RefreshDevices()
	}

	devs := []*Device{}
	for _, ic := range c.devices.IntelliClimates {
		devs = append(devs, ic.Device)
	}

	for _, id := range c.devices.IntelliDoses {
		devs = append(devs, id.Device)
	}

	return devs
}

// Growrooms returns a slice of all the known growrooms
func (c *Client) Growrooms() []*Growroom {
	grs := []*Growroom{}
	for _, gr := range grs {
		grs = append(grs, gr)
	}
	return grs
}

// ListDevicesBySerial will return the serial numbers of all known devices
func (c *Client) ListDevicesBySerial() []string {
	if len(c.growrooms) == 0 {
		c.RefreshDevices()
	}

	serials := []string{}
	for _, ic := range c.devices.IntelliClimates {
		serials = append(serials, ic.ID)
	}

	for _, id := range c.devices.IntelliDoses {
		serials = append(serials, id.ID)
	}

	return serials
}

// ListGrowrooms returns a list of the names known growrooms
func (c *Client) ListGrowrooms() []string {
	if len(c.growrooms) == 0 {
		c.RefreshDevices()
	}

	grs := []string{}
	for n := range c.growrooms {
		grs = append(grs, n)
	}

	return grs
}

// Growroom returns the growroom with the given name, and a false if it
// wasn't found
func (c *Client) Growroom(name string) (*Growroom, bool) {
	if len(c.growrooms) == 0 {
		c.RefreshDevices()
	}

	growroom, exists := c.growrooms[name]
	return growroom, exists
}

// GetGrowroom is deprecated in favour of Growroom
func (c *Client) GetGrowroom(name string) (*Growroom, bool) {
	return c.Growroom(name)
}

// UpdateAllGrowrooms - Updated all growrooms
func (c *Client) UpdateAllGrowrooms() {
	for _, gr := range c.growrooms {
		gr.Update()
	}
}

// UpdateGrowroom - returns the growroom specified and an error
func (c *Client) UpdateGrowroom(gr string) error {
	growroom, exists := c.growrooms[gr]

	if !exists {
		return fmt.Errorf("No growroom called %s found", gr)
	}

	return growroom.Update()
}

// GetGrowroomReading - returns the reading for the growroom specified as a string, also return an error
func (c *Client) GetGrowroomReading(gr, reading string) (string, error) {
	growroom, exists := c.growrooms[gr]

	if !exists {
		return "no growroom found", fmt.Errorf("No growroom called %s found", gr)
	}

	avaliable, r := growroom.GetReading(reading)

	if !avaliable {
		return "no reading found", fmt.Errorf("No reading called %s found in growroom %s", r, gr)
	}

	return r, nil
}
