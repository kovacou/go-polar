// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package polar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/kovacou/go-types"
)

// Sleep is the representation of an Polar sleep.
type Sleep struct {
	Continuity   float64 `json:"continuity"`
	Score        uint64  `json:"sleep_score"`
	Light        uint64  `json:"light_sleep"`
	Deep         uint64  `json:"deep_sleep"`
	Rem          uint64  `json:"rem_sleep"`
	Interruption uint64  `json:"total_interruption_duration"`
	DeviceID     string  `json:"device_id"`

	// Dates
	Start types.DateTime `json:"sleep_start_time"`
	End   types.DateTime `json:"sleep_end_time"`
}

// Sleep returns information of sleep for the given date.
func (p *polar) Sleep(date string) (out Sleep, err error) {
	r, err := p.GET(fmt.Sprintf("/v3/users/sleep/%s", date), RequestParams{
		WithBearer: true,
	})

	defer closeHTTPResponse(r)
	if err != nil {
		return
	}

	if r.StatusCode == http.StatusOK {
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &out)
	}
	return
}

// LastSleeps returns information about the last 28 days of sleep.
func (p *polar) LastSleeps() (out []Sleep, err error) {
	r, err := p.GET("/v3/users/sleep", RequestParams{
		WithBearer: true,
	})

	defer closeHTTPResponse(r)
	if err != nil {
		return
	}

	if r.StatusCode == http.StatusOK {
		j, _ := gabs.ParseJSONBuffer(r.Body)
		err = json.Unmarshal(j.Path("nights").Bytes(), &out)
	}
	return
}
