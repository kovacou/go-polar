// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package polar

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/kovacou/go-convert"
	"github.com/kovacou/go-env"
	"github.com/kovacou/go-types"
)

// Client is the client interface of Polar service.
type Client interface {
	// AuthorizationURL returns the URL to get the authorization from the user.
	AuthorizationURL() string
}

// RequestParams define the parameter to request the API.
type RequestParams struct {
	Queries types.Map
	Values  types.Map
}

// NewEnv create a new Polar client from environment variables.
func NewEnv() Client {
	cfg := Config{}
	env.Unmarshal(&cfg)

	return New(cfg)
}

// New create a new Polar client from the given config.
func New(cfg Config) Client {
	cfg.Bearer = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.ClientID, cfg.ClientSecret)))

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
	LimitRate <-chan time.Time
}

func (p *polar) AuthorizationURL() string {
	return fmt.Sprintf("https://flow.polar.com/oauth2/authorization?response_type=code&client_id=%s", p.cfg.ClientID)
}

func (p *polar) AuthenticationHandler(r http.Request) {
	code := r.URL.Query().Get("code")
}

func (p *polar) Request(method, endpoint string, params RequestParams) (*http.Response, error) {
	var values url.Values

	if params.Queries == nil {
		params.Queries = types.Map{}
	}

	if params.Values == nil {
		params.Values = types.Map{}
	}

	if method == http.MethodPost {
		for key, val := range params.Values {
			values.Set(key, convert.String(val))
		}
	}

	r, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Authorisation", fmt.Sprintf("Basic %s", p.cfg.Bearer))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept", "application/json;charset=UTF-8")

	return p.Do(r)
}

func (p *polar) GET(endpoint string, params RequestParams) (*http.Response, error) {
	return p.Request(http.MethodGet, endpoint, params)
}

func (p *polar) POST(endpoint string, params RequestParams) (*http.Response, error) {
	return p.Request(http.MethodPost, endpoint, params)
}
