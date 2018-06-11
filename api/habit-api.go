package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// HabiticaAPI Main client for interacting with Habitica API via HTTP
type HabiticaAPI struct {
	client   *http.Client
	hostURL  string
	userAuth UserToken
}

// NewHabiticaAPI is a function for creating a new client api. Can pass in prexisting client
// for proxies or what not.
func NewHabiticaAPI(client *http.Client, hosturl string) *HabiticaAPI {
	var api HabiticaAPI

	if client == nil {
		api.client = &http.Client{}
	}

	api.hostURL = hosturl
	if hosturl == "" {
		api.hostURL = `https://habitica.com/api`
	}

	return &api
}

// Do is a wrapper function around the api's http.client.Do but Marshals any json struct
// given to it. Also, it will parse http status errors over 400 and return an error.
func (api *HabiticaAPI) Do(req *http.Request, responseType interface{}) error {
	api.addAuthHeaders(req)

	body, err := api.doRequest(req)
	if err != nil {
		return err
	}

	return parseResponse(body, responseType)
}

func (api *HabiticaAPI) addAuthHeaders(req *http.Request) {
	if api.userAuth.APIToken != "" {
		req.Header.Add("x-api-key", api.userAuth.APIToken)
	}
	if api.userAuth.ID != "" {
		req.Header.Add("x-api-user", api.userAuth.ID)
	}
}

func (api *HabiticaAPI) doRequest(req *http.Request) ([]byte, error) {
	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	return readBody(res), parseStatusErrors(res)
}

func readBody(res *http.Response) []byte {
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

func parseResponse(body []byte, object interface{}) error {
	var hres habiticaResponse
	err := json.Unmarshal(body, &hres)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hres.Data, &object)
	if err != nil {
		return err
	}

	return nil
}

type habiticaResponse struct {
	Data   json.RawMessage
	Error  string
	Errors []struct {
		Message string
		Path    string
	}
}

// Get will return response from the passed in route of Habitica Api.
// It will also return errors in either HTTP Protocol or if status
// code is equal to or above 400.
func (api *HabiticaAPI) Get(route string, responseType interface{}) error {
	req, err := http.NewRequest("GET", api.hostURL+"/v3"+route, nil)
	if err != nil {
		return err
	}

	err = api.Do(req, responseType)
	return err
}

// Post will take in route, request data as a struct, and response as a struct and output errors for
// marshalling either.
func (api *HabiticaAPI) Post(url string, requestObject interface{}, responseObject interface{}) error {
	data, merr := json.Marshal(requestObject)
	if merr != nil {
		return merr
	}

	req, rerr := http.NewRequest("POST", api.hostURL+"/v3"+url, bytes.NewBuffer(data))
	if rerr != nil {
		return rerr
	}

	err := api.Do(req, &responseObject)
	return err
}
