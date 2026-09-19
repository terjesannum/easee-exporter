package easee

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

const (
	apiUrl = "https://api.easee.cloud"
	// The observations API that replaces /api/chargers/{id}/state is served
	// from api.easee.com.
	stateApiUrl = "https://api.easee.com"
)

type Client struct {
	HttpClient *http.Client
	username   string
	password   string
	token      *oauth2.Token
	mu         sync.Mutex
}

func NewClient(ctx context.Context, username, password string) *Client {
	c := &Client{
		username: username,
		password: password,
	}
	c.HttpClient = oauth2.NewClient(ctx, c)
	return c
}

func (c *Client) Token() (*oauth2.Token, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error
	if c.token == nil {
		data := struct {
			Username string `json:"userName"`
			Password string `json:"password"`
		}{
			Username: c.username,
			Password: c.password,
		}
		err = c.getToken("/api/accounts/login", data)
	} else if time.Now().After(c.token.Expiry) {
		log.Printf("Token expired: %v\n", c.token.Expiry)
		data := struct {
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
		}{
			AccessToken:  c.token.AccessToken,
			RefreshToken: c.token.RefreshToken,
		}
		err = c.getToken("/api/accounts/refresh_token", data)
	}
	return c.token, err
}

func (c *Client) getToken(endpoint string, data interface{}) error {
	log.Printf("Getting token from %s\n", endpoint)
	var auth struct {
		AccessToken  string `json:"accessToken"`
		ExpiresIn    int32  `json:"expiresIn"`
		TokenType    string `json:"tokenType"`
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.doRequest(fmt.Sprintf("%s%s", apiUrl, endpoint), data, &auth); err != nil {
		return err
	}
	if auth.ExpiresIn == 0 {
		return errors.New("Getting token failed: no expiry in response")
	}
	c.token = &oauth2.Token{
		AccessToken:  auth.AccessToken,
		TokenType:    auth.TokenType,
		RefreshToken: auth.RefreshToken,
		Expiry:       time.Now().Add(time.Second * time.Duration(auth.ExpiresIn)),
	}
	return nil
}

func (c *Client) doRequest(url string, body, result interface{}) error {
	var res *http.Response
	var err error
	if body == nil {
		res, err = c.HttpClient.Get(url)
	} else {
		var reqBody []byte
		if reqBody, err = json.Marshal(body); err != nil {
			return err
		}
		// Deliberately not c.HttpClient: this is the token request itself, and
		// the oauth2 client would call back into Token to authenticate it.
		res, err = http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// Without this, an error response decodes into result and surfaces as a
	// confusing json error instead of the status the API actually returned.
	if res.StatusCode < 200 || res.StatusCode > 299 {
		msg, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return fmt.Errorf("%s returned %s: %s", url, res.Status, strings.TrimSpace(string(msg)))
	}
	dec := json.NewDecoder(res.Body)
	return dec.Decode(result)
}

func (c *Client) Chargers() ([]Charger, error) {
	var chargers []Charger
	err := c.doRequest(fmt.Sprintf("%s/api/chargers", apiUrl), nil, &chargers)
	return chargers, err
}

// ChargerState reads a charger's state from the observations API. It replaces
// /api/chargers/{id}/state, which was removed on 2026-09-01. All observations
// are requested in one call to stay well inside the endpoint's rate limit of
// 100 requests per 5 minutes.
func (c *Client) ChargerState(charger string) (ChargerState, error) {
	var res observationsResponse
	if err := c.doRequest(chargerStateUrl(charger), nil, &res); err != nil {
		return ChargerState{}, err
	}
	return newChargerState(&res), nil
}

func chargerStateUrl(charger string) string {
	return fmt.Sprintf("%s/state/%s/observations?ids=%s", stateApiUrl, charger, chargerStateObservationIds)
}
