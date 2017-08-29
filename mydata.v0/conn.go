package mydata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Dial creates a new connection to the API that can be used to get data points
func Dial(deviceID string) *Conn {
	url := fmt.Sprintf("https://api.autogrow.com/multi/v0/%s", deviceID)
	return &Conn{&http.Client{}, deviceID, url}
}

// Conn represents a connection to the API
type Conn struct {
	client *http.Client
	uuid   string
	URL    string
}

// DeviceID returns the device ID that the API will retrieve data points for
func (c *Conn) DeviceID() string {
	return c.uuid
}

func (c *Conn) getPoints(uri string) (Points, error) {
	req, err := http.NewRequest("GET", uri, strings.NewReader(""))
	if err != nil {
		return Points{}, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return Points{}, err
	}

	if res.StatusCode != 200 {
		return Points{}, fmt.Errorf("unexpected status %d", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Points{}, err
	}

	parsedErr, err := parseError(b)
	if err != nil {
		return Points{}, err
	}

	if parsedErr != nil {
		return Points{}, parsedErr
	}

	points := Points{}
	err = json.Unmarshal(b, &points)

	return points, err
}

func (c *Conn) getRecord(uri string, object interface{}) error {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	fmt.Printf("req %+v\n", req)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	// if res.StatusCode != 200 {
	// return fmt.Errorf("unexpected status %d", res.StatusCode)
	// }

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", res)
	fmt.Println(string(b))
	parsedErr, err := parseError(b)
	if err != nil {
		return err
	}

	if parsedErr != nil {
		return parsedErr
	}

	return json.Unmarshal(b, object)
}

// Labels returns all of the section labels known by the API for the given device ID
func (c *Conn) Labels() (Labels, error) {
	uri := fmt.Sprintf("%s/labels", c.URL)
	res, err := http.Get(uri)
	if err != nil {
		return Labels{}, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Labels{}, err
	}

	parsedErr, err := parseError(b)
	if err != nil {
		return Labels{}, err
	}

	if parsedErr != nil {
		return Labels{}, parsedErr
	}

	labels := Labels{}
	err = json.Unmarshal(b, &labels)

	return labels, err
}

// Site starts a site query for site level metrics such as weather data
func (c *Conn) Site() *Query {
	return &Query{env: "site", name: "", days: 45, conn: c}
}

// Compartment starts a compartment query
func (c *Conn) Compartment(name string) *Query {
	return &Query{env: "comp", name: name, days: 45, conn: c}
}

// GrowRoom starts a grow room query
func (c *Conn) GrowRoom(name string) *Query {
	return &Query{env: "room", name: name, days: 45, conn: c}
}

// Irrigator starts an irrigation system query
func (c *Conn) Irrigator(name string) *Query {
	return &Query{env: "irrigator", name: name, days: 45, conn: c}
}

// Monitor starts an environment monitor query
func (c *Conn) Monitor(name string) *Query {
	return &Query{env: "monitor", name: name, days: 45, conn: c}
}

// Field starts an open field query
func (c *Conn) Field(name string) *Query {
	return &Query{env: "field", name: name, days: 45, conn: c}
}
