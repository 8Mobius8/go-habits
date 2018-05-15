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
	hostUrl string
}

func NewHabiticaApi(client *http.Client, hosturl string) *HabiticaAPI {
	var api HabiticaAPI

	if client == nil {
		api.client = &http.Client{}
	}

	if hosturl == "" {
		api.hostUrl = "https://habitica.com/api"
	} else {
		api.hostUrl = hosturl
	}

	return &api
}

func (api *HabiticaAPI) Status() (*http.Response, error) {
	return api.get("/status")
}

func (api *HabiticaAPI) Tasks() (*http.Response, error) {
	return api.get("/user")
}

func (api *HabiticaAPI) get(route string) (*http.Response, error) {
	res, err := api.client.Get(api.hostUrl + "/v3" + route)

	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, &ApiError{http.StatusText(res.StatusCode), res.StatusCode}
	}

	return res, nil
}

type ApiError struct {
	msg  string
	code int
}

func (err *ApiError) Error() string {
	return err.msg
}

func (api *HabiticaAPI) ParseResponse(res *http.Response, responseType interface{}) interface{} {

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &responseType)
	if err != nil {
		log.Fatal(err)
	}

	return responseType
}
