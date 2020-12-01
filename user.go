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

	"github.com/kovacou/go-types"
)

// User is the representation of an Polar User.
type User struct {
	ID               uint64 `json:"polar-user-id"`
	Firstname        string `json:"first-name"`
	Lastname         string `json:"last-name"`
	Birthday         string `json:"birthdate"`
	Gender           string `json:"gender"`
	RegistrationDate string `json:"registration-date"`
}

// User returns information about the current user.
func (p *polar) User() (u User, err error) {
	r, err := p.GET(fmt.Sprintf("/v3/users/%d", p.userID), RequestParams{
		WithBearer: true,
	})

	defer closeHTTPResponse(r)
	if err != nil {
		return
	}

	if r.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(body, &u)
	}

	return
}

// RegisterUser register a new user to polar the application.
func (p *polar) RegisterUser() (err error) {
	r, err := p.POST("/v3/users", RequestParams{
		WithBearer: true,
		Values: types.Map{
			"member-id": p.userID,
		},
	})

	closeHTTPResponse(r)
	return
}

// UnregisterUser remove an user from the polar application.
func (p *polar) UnregisterUser() (err error) {
	r, err := p.DELETE(fmt.Sprintf("/v3/users/%d", p.userID), RequestParams{
		WithBearer: true,
	})

	closeHTTPResponse(r)
	return
}
