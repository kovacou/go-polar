// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package polar

import (
	"encoding/base64"
	"fmt"

	"github.com/kovacou/go-env"
)

var (
	// cfgEnviron contains the loaded configuration from environment.
	cfgEnviron Config
)

// init loads the global configuration.
func init() {
	_ = env.Unmarshal(&cfgEnviron)
}

// Config is the configuration for the client to request Polar API.
type Config struct {
	Host         string `json:"host" env:"POLAR_HOST"`
	LimitRate    uint8  `json:"limit_rate" env:"POLAR_LIMITRATE"`
	Timeout      uint16 `json:"timeout" env:"POLAR_TIMETOUT"`
	ClientID     string `json:"client_id" env:"POLAR_ID"`
	ClientSecret string `json:"client_secret" env:"POLAR_SECRET"`

	// authorizationToken is the cache of the basic authentication.
	authorizationToken string
}

// AuthorizationToken return the token based on configuration.
func (c *Config) AuthorizationToken() string {
	if c.authorizationToken == "" {
		c.authorizationToken = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.ClientID, c.ClientSecret)))
	}

	return c.authorizationToken
}
