// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package polar

import (
	"bytes"
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
	AuthorizationURL(string) string

	// AuthorizationAccessToken get an access token from the user code.
	AuthorizationAccessToken(string) (AccessToken, error)

	// Exercises returns the exercises of the user.
	Exercises() ([]Exercise, error)

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

	// LastRecharges returns information about the last 28 night recharges.
	LastRecharges() ([]Recharge, error)

	// LastSleeps returns information about the last 28 days of sleep.
	LastSleeps() ([]Sleep, error)
}

// RequestParams define the parameter to request the API.
type RequestParams struct {
	Queries            types.Map
	Values             types.Map
	WithBearer         bool
	WithFormURLEncoded bool
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

// close is used as defer to automatically close the body and prevent memory leak.
func closeHTTPResponse(r *http.Response) {
	if r != nil && r.Body != nil {
		r.Body.Close()
	}
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
func (p *polar) AuthorizationURL(state string) string {
	return fmt.Sprintf("https://flow.polar.com/oauth2/authorization?response_type=code&client_id=%s&state=%s", p.cfg.ClientID, state)
}

// AuthorizationAccessToken get an access token from the user and authorize his profile to be fetched.
func (p *polar) AuthorizationAccessToken(code string) (result AccessToken, err error) {
	r, err := p.Request(http.MethodPost, "https://polarremote.com/v2/oauth2/token", RequestParams{
		WithFormURLEncoded: true,
		Values: types.Map{
			"grant_type": "authorization_code",
			"code":       code,
		},
	})

	defer closeHTTPResponse(r)
	if err != nil {
		return
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

// CreateTransaction creates a new transaction.
func (p *polar) CreateTransaction() {
}

// Commit a transaction.
func (p *polar) Commit() {
}

// Request build a new request from the input and return the response.
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

	// To indicate Values must be encoded as FormURLEncoded, please pass WithFormURLEncoded with true.
	// By default, JSON will be used.
	if method == http.MethodPost {
		if params.WithFormURLEncoded {
			val := url.Values{}
			for k, v := range params.Values {
				val.Set(k, convert.String(v))
			}
			values = strings.NewReader(val.Encode())
		} else {
			b, _ := json.Marshal(params.Values)
			values = bytes.NewReader(b)
		}
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

	// Managing the content type to use : some endpoint need JSON and some need form encoded.
	// To indicate Values must be encoded as FormURLEncoded, please pass WithFormURLEncoded with true.
	contentType := "application/x-www-form-urlencoded"
	if method == http.MethodPost && !params.WithFormURLEncoded {
		contentType = "application/json"
	}

	r.Header.Set("Authorization", token)
	r.Header.Set("Content-Type", contentType)
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
