package easee

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

const (
	apiUrl = "https://api.easee.cloud"
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
	c.doRequest(endpoint, data, &auth)
	if auth.ExpiresIn == 0 {
		return errors.New("Getting token failed")
	}
	c.token = &oauth2.Token{
		AccessToken:  auth.AccessToken,
		TokenType:    auth.TokenType,
		RefreshToken: auth.RefreshToken,
		Expiry:       time.Now().Add(time.Second * time.Duration(auth.ExpiresIn)),
	}
	return nil
}

func (c *Client) doRequest(endpoint string, body, result interface{}) error {
	var res *http.Response
	var err error
	if body == nil {
		res, err = c.HttpClient.Get(fmt.Sprintf("%s%s", apiUrl, endpoint))
	} else {
		reqBody, err := json.Marshal(body)
		if err == nil {
			res, err = http.Post(fmt.Sprintf("%s%s", apiUrl, endpoint), "application/json", bytes.NewBuffer(reqBody))
		}
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	return dec.Decode(&result)
}

func (c *Client) Chargers() ([]Charger, error) {
	var chargers []Charger
	err := c.doRequest("/api/chargers", nil, &chargers)
	return chargers, err
}

func (c *Client) ChargerState(charger string) (ChargerState, error) {
	var state ChargerState
	err := c.doRequest(fmt.Sprintf("/api/chargers/%s/state", charger), nil, &state)
	return state, err
}
