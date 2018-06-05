package ig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

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

func (c *Client) getRefreshTime() time.Duration {
	rTime := c.auth.ExpiresIn
	if rTime > 60 {
		rTime -= 60
	}
	return time.Duration(rTime)
}
