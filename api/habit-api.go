package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// HabiticaAPI Main client for interacting with Habitica API via HTTP
type HabiticaAPI struct {
	client  *http.Client
	hostURL string
}

// NewHabiticaAPI is a function for creating a new client api. Can pass in prexisting client
// for proxies or what not.
func NewHabiticaAPI(client *http.Client, hosturl string) *HabiticaAPI {
	var api HabiticaAPI

	if client == nil {
		api.client = &http.Client{}
	}

	if hosturl == "" {
		api.hostURL = `https://habitica.com/api`
	} else {
		api.hostURL = hosturl
	}

	return &api
}

// Get will return response from the passed in route of Habitica Api.
// It will also return errors in either HTTP Protocol or if status
// code is equal to or above 400.
func (api *HabiticaAPI) Get(route string) ([]byte, error) {
	res, protoerr := api.client.Get(api.hostURL + "/v3" + route)

	if protoerr != nil {
		return nil, protoerr
	}

	return parseHTTPBody(res), parseStatusErrors(res)
}

func parseHTTPBody(res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func parseStatusErrors(res *http.Response) error {
	if res.StatusCode >= 400 {
		return &GoHabitsError{http.StatusText(res.StatusCode), res.StatusCode}
	}

	return nil
}

// ParseResponse will unmarshal any byte[] array with the passed in `responsetype`
// parameter.
func (api *HabiticaAPI) ParseResponse(body []byte, responseType interface{}) {
	err := json.Unmarshal(body, &responseType)
	if err != nil {
		log.Fatal(err)
	}
}
