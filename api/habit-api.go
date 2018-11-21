package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/amoghe/distillog"
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
	req.Header.Add("content-type", "application/json")

	res, err := api.client.Do(req)
	if err != nil {
		return NewGoHabitsError(err.Error(), 1, "")
	}

	body := readBody(res)
	return parseResponse(body, res, responseType)
}

func readBody(res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Errorln(err)
	}
	return body
}

func parseResponse(body []byte, res *http.Response, object interface{}) error {
	var hres habiticaResponse
	var err error

	if len(body) > 0 {
		err = json.Unmarshal(body, &hres)
		if err != nil {
			return err
		}
	}

	if hres.Error != "" || res.StatusCode >= 400 {
		return parseHabitsServerError(hres, res)
	}

	if len(body) > 0 {
		err = json.Unmarshal(hres.Data, &object)
		if err != nil {
			return err
		}
	}
	return nil
}

type habiticaResponse struct {
	Data    json.RawMessage
	Error   string
	Message string
	Errors  []struct {
		Message string
	}
}

func parseHabitsServerError(hres habiticaResponse, res *http.Response) error {
	errMessage := hres.Error + "\n" + hres.Message
	for _, errorMessage := range hres.Errors {
		errMessage += "\n" + errorMessage.Message
	}
	return NewGoHabitsError(errMessage, res.StatusCode, res.Request.URL.EscapedPath())
}

// Get will return response from the passed in route of Habitica Api.
// It will also return errors in either HTTP Protocol or if status
// code is equal to or above 400.
func (api *HabiticaAPI) Get(route string, responseType interface{}) error {
	req, err := http.NewRequest("GET", api.hostURL+"/v3"+route, nil)
	if err != nil {
		return err
	}

	return api.Do(req, responseType)
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
	err := api.Do(req, responseObject)
	return err
}
