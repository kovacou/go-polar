// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package polar

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Exercise is the representation of an Polar exercise.
type Exercise struct {
	ID                     string  `json:"id"`
	Sport                  string  `json:"sport"`
	Device                 string  `json:"device"`
	HasRoute               bool    `json:"has_route"`
	Calories               float64 `json:"calories"`
	Distance               float64 `json:"distance"`
	TrainingLoad           float64 `json:"training_load"`
	FatPercentage          uint8   `json:"fat_percentage"`
	ProteinPercentage      uint8   `json:"protein_percentage"`
	CarbohydratePercentage uint8   `json:"carbohydrate_percentage"`
	Duration               string  `json:"duration"`
	StartTime              string  `json:"start_time"`
	HeartRate              struct {
		Average uint16 `json:"average"`
		Maximum uint16 `json:"maximum"`
	} `json:"heart_rate"`
}

// Exercises returns exercises of the user.
func (p *polar) Exercises() (out []Exercise, err error) {
	r, err := p.GET("/v3/exercises", RequestParams{
		WithBearer: true,
	})

	defer closeHTTPResponse(r)
	if err != nil {
		return
	}

	if r.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &out)
	}
	return
}
