package polar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Jeffail/gabs"
)

// Recharge is the representation of a Polar "Nightly Recharge".
type Recharge struct {
	Date                    string  `json:"date"`
	HeartRateAvg            uint64  `json:"heart_rate_avg"`
	BeatToBeatAvg           uint64  `json:"beat_to_beat_avg"`
	HeartRateVariabilityAvg uint64  `json:"heart_rate_variability_avg"`
	BreathingRateAvg        float64 `json:"breathing_rate_avg"`
	ANSCharge               float64 `json:"ans_charge"`
}

// Recharge returns information about the recharge of the given date.
func (p *polar) Recharge(date time.Time) (out Recharge, err error) {
	r, err := p.GET(fmt.Sprintf("/v3/users/nightly-recharge/%s", date), RequestParams{
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

// LastRecharges returns information about the last 28 night recharges.
func (p *polar) LastRecharges() (out []Recharge, err error) {
	r, err := p.GET("/v3/users/nightly-recharge", RequestParams{
		WithBearer: true,
	})

	defer closeHTTPResponse(r)
	if err != nil {
		return
	}

	if r.StatusCode == http.StatusOK {
		j, _ := gabs.ParseJSONBuffer(r.Body)
		err = json.Unmarshal(j.Path("recharges").Bytes(), &out)
	}
	return
}
