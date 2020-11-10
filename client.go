// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package polar

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kovacou/go-convert"
	"github.com/kovacou/go-env"
	"github.com/kovacou/go-types"
)

var (
	// cfgEnviron contains the loaded configuration from environment.
	cfgEnviron Config
)

// Client is the client interface of Polar service.
type Client interface {
	// AuthorizationURL returns the URL to get the authorization from the user.
	AuthorizationURL() string

	// AuthorizationAccessToken get an access token from the user code.
	AuthorizationAccessToken(string) (AccessToken, error)

	// SetBearer
	SetBearer(string)

	// SetUserID
	SetUserID(uint64)

	// RegisterUser register a new user to polar the application.
	RegisterUser() error

	// UnregisterUser remove an user from the polar application.
	UnregisterUser() error

	// User returns information about the current user.
	User() (u User, err error)
}

// RequestParams define the parameter to request the API.
type RequestParams struct {
	Queries    types.Map
	Values     types.Map
	WithBearer bool
}

// AccessToken is the response of the Authorization.
type AccessToken struct {
	Value     string `json:"access_token"`
	Type      string `json:"token_type"`
	ExpiresIn uint   `json:"expires_in"`
	XUserID   uint64 `json:"x_user_id"`
}

// init loads the global configuration.
func init() {
	env.Unmarshal(&cfgEnviron)
}

// NewEnv create a new Polar client from environment variables.
func NewEnv() Client {
	return New(cfgEnviron)
}

// New create a new Polar client from the given config.
func New(cfg Config) Client {
	return &polar{
		cfg: cfg,
		Client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// polar is the HTTP Client of the service.
type polar struct {
	*http.Client

	cfg       Config
	bearer    string
	userID    uint64
	LimitRate <-chan time.Time
}

// AuthorizationURL returns the URL to get the authorization from the user.
func (p *polar) AuthorizationURL() string {
	return fmt.Sprintf("https://flow.polar.com/oauth2/authorization?response_type=code&client_id=%s", p.cfg.ClientID)
}

// AuthorizationAccessToken get an access token from the user and authorize his profile to be fetched.
func (p *polar) AuthorizationAccessToken(code string) (result AccessToken, err error) {
	r, err := p.Request(http.MethodPost, "https://polarremote.com/v2/oauth2/token", RequestParams{
		Values: types.Map{
			"grant_type": "authorization_code",
			"code":       code,
		},
	})

	if err != nil {
		return
	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	body, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &result)
	return
}

// SetBearer set a new bearer token to the client.
func (p *polar) SetBearer(bearer string) {
	p.bearer = fmt.Sprintf("Bearer %s", bearer)
}

// SetUserID set an userID to the client.
func (p *polar) SetUserID(id uint64) {
	p.userID = id
}

// RegisterUser register a new user to polar the application.
func (p *polar) RegisterUser() (err error) {
	_, err = p.POST("/v3/users", RequestParams{
		Values: types.Map{
			"member-id": p.userID,
		},
	})

	return
}

// UnregisterUser remove an user from the polar application.
func (p *polar) UnregisterUser() (err error) {
	_, err = p.DELETE("/v3/users", RequestParams{
		Values: types.Map{
			"member-id": p.userID,
		},
	})

	return
}

// User is the representation of an Polar User.
type User struct {
	ID               uint64    `json:"polar-user-id"`
	RegistrationDate time.Time `json:"registration-date"`
	Firstname        string    `json:"first-name"`
	Lastname         string    `json:"last-name"`
	Birthday         string    `json:"birthday"`
	Gender           string    `json:"gender"`
}

// User returns information about the current user.
func (p *polar) User() (u User, err error) {
	r, err := p.GET("/v3/users", RequestParams{
		Values: types.Map{
			"member-id": p.userID,
		},
	})

	if r.StatusCode == http.StatusOK {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &u)
	}

	return
}

// SLeep is the representation of an Polar sleep.
type Sleep struct {
}

// Sleep returns information of sleep for the given date.
func (p *polar) Sleep(date time.Time) (s Sleep, err error) {

	return
}

// CreateTransaction creates a new transaction.
func (p *polar) CreateTransaction() {
}

// Commit a transaction.
func (p *polar) Commit() {
}

// Request
func (p *polar) Request(method, uri string, params RequestParams) (*http.Response, error) {
	var (
		values io.Reader
		token  string
	)

	if params.Queries == nil {
		params.Queries = types.Map{}
	}

	if params.Values == nil {
		params.Values = types.Map{}
	}

	if method == http.MethodPost {
		v := make(url.Values)
		for key, val := range params.Values {
			v.Set(key, convert.String(val))
		}

		values = strings.NewReader(v.Encode())
	}

	r, err := http.NewRequest(method, uri, values)
	if err != nil {
		return nil, err
	}

	// Manage authorization to use.
	if params.WithBearer {
		token = p.bearer
	} else {
		token = p.cfg.AuthorizationToken()
	}

	r.Header.Set("Authorization", token)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept", "application/json;charset=UTF-8")

	return p.Do(r)
}

// GET
func (p *polar) GET(endpoint string, params RequestParams) (*http.Response, error) {
	return p.Request(http.MethodGet, p.cfg.Host+endpoint, params)
}

// POST
func (p *polar) POST(endpoint string, params RequestParams) (*http.Response, error) {
	return p.Request(http.MethodPost, p.cfg.Host+endpoint, params)
}

// DELETE
func (p *polar) DELETE(endpoint string, params RequestParams) (*http.Response, error) {
	return p.Request(http.MethodDelete, p.cfg.Host+endpoint, params)
}
