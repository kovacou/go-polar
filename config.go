// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package polar

// Config is the configuration for the client to request Polar API.
type Config struct {
	Host         string `json:"host"`
	LimitRate    uint8  `json:"limit_rate"`
	Timeout      uint16 `json:"timeout"`
	Bearer       string `json:"bearer"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
