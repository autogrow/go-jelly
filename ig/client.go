package ig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	igBaseURL = "https://api.autogrow.com/v1"

	igTokenURI       = igBaseURL + "/auth/token"
	igRefreshURI     = igBaseURL + "/auth/token/refresh"
	igDeviceURI      = "/intelligrow/devices?username="
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
	auth               authResponse
	CheckedDevices     float64
	growrooms          map[string]*Growroom
	devices            *Devices
	tokenRefresherQuit chan bool
	url                url.URL
}

// NewClient creates a new client with the given username and password.  It will
// return an error if the authentication fails
func NewClient(user, pass string) (*Client, error) {
	c := &Client{
		Client:    &http.Client{Timeout: time.Second * 30},
		lock:      new(sync.RWMutex),
		username:  user,
		password:  pass,
		growrooms: make(map[string]*Growroom),
	}

	c.url.Scheme = "https"
	c.url.Host = "api.autogrow.com"

	// Initialize the devices object in the structure, this is blank object
	c.devices = NewDevices()

	err := c.authenticate()
	if err != nil {
		return c, err
	}

	go c.autoAuthExtender()

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
	data, err := json.Marshal(map[string]string{"username": c.username, "password": c.password})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", igTokenURI, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Unable to get tokens %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Unable to get tokens %s", err)
	}

	return c.processAuthResponse(res)
}

// Token returns the current token for the client
func (c *Client) getToken() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.auth.Token
}

type authResponse struct {
	Token        string  `json:"api_access_token"`
	ExpiresIn    float64 `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
}

func (auth authResponse) reauthPayload(user string) io.Reader {
	data, _ := json.Marshal(map[string]string{"username": user, "refresh_token": auth.RefreshToken})
	return bytes.NewBuffer(data)
}

func (auth authResponse) validate() error {
	if auth.Token == "" {
		return fmt.Errorf("Response from server doesn't contain an api token")
	}

	if auth.ExpiresIn == 0 {
		return fmt.Errorf("Response from server doesn't contain a refresh time")
	}

	if auth.RefreshToken == "" {
		return fmt.Errorf("Response from server doesn't contain a refresh token")
	}

	return nil
}

func (c *Client) processAuthResponse(res *http.Response) error {
	auth := authResponse{}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("couldn't read response body: %s", err)
	}

	err = json.Unmarshal(data, &auth)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal response body: %s", err)
	}

	if err := auth.validate(); err != nil {
		return err
	}

	c.auth = auth
	return nil
}

func (c *Client) extendAuth() error {
	req, err := http.NewRequest("POST", igRefreshURI, c.auth.reauthPayload(c.username))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		return err
	}

	err = c.processAuthResponse(res)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) autoAuthExtender() {
	c.tokenRefresherQuit = make(chan bool)
	defer func() { c.tokenRefresherQuit = nil }()
	timer := time.NewTimer(c.getRefreshTime())

	for {
		select {
		case <-c.tokenRefresherQuit:
			return

		case <-timer.C:
			if err := c.extendAuth(); err != nil {
				// TODO: handle err
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
	rTime := c.auth.ExpiresIn
	if rTime > 60 {
		rTime -= 60
	}
	return time.Duration(rTime)
}

// GetDevices is deprecated in favour of RefreshDevices
func (c *Client) GetDevices() error {
	return c.RefreshDevices()
}

// SaveDevice will save the config and state of the given device
func (c *Client) SaveDevice(i Intelli) error {
	payload, err := i.StatePayload()
	if err != nil {
		return err
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	res, err := c.doRequest("PUT", c.buildURL("/intelligrow/devices"), bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to save state/config: %s", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected http status: %d", res.StatusCode)
	}

	return nil
}

func (c *Client) doRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.getToken())

	return c.Do(req)
}

func (c *Client) buildURL(path string, queries ...string) string {
	u := c.url
	u.Path = filepath.Join("v1", path)
	u.RawQuery = strings.Join(queries, "&")
	return u.String()
}

// RefreshDevices will get the latest data from the API and update all known structs
func (c *Client) RefreshDevices() error {
	res, err := c.doRequest("GET", c.buildURL("/intelligrow/devices", "username="+c.username), nil)
	if err != nil {
		return fmt.Errorf("failed to refresh devices; %s", err)
	}

	msi := make(map[string]interface{})
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("couldn't read response body: %s", err)
	}

	if err := json.Unmarshal(data, &msi); err != nil {
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

	data, err = json.Marshal(devices)
	if err != nil {
		return fmt.Errorf("Error marshalling devices interface: %s", err)
	}

	var igDevices []*Device
	if err := json.Unmarshal(data, &igDevices); err != nil {
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
