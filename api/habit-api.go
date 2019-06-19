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
	logger   log.Logger
}

// NewHabiticaAPI is a function for creating a new client api. Can pass in prexisting client
// for proxies or what not.
func NewHabiticaAPI(client *http.Client, hosturl string, logger log.Logger) *HabiticaAPI {
	var api HabiticaAPI

	api.logger = logger
	if api.logger == nil {
		api.logger = log.NewNullLogger("go-habits")
		api.logger.Debugln("No logger configured in, using /dev/null")
	}

	api.client = client
	if client == nil {
		api.client = &http.Client{}
		api.logger.Debugln("No client configured in, using default client")
	}

	api.hostURL = hosturl
	if hosturl == "" {
		api.hostURL = `https://habitica.com/api`
	}
	api.logger.Debugln("Habitica url configured to: ", api.hostURL)

	return &api
}

// GetHostURL returns client's hostURL configured.
func (api *HabiticaAPI) GetHostURL() string {
	return api.hostURL
}

// Do is a wrapper function around the api's http.client.Do but Marshals any json struct
// given to it. Also, it will parse http status errors over 400 and return an error.
func (api *HabiticaAPI) Do(req *http.Request, responseType interface{}) error {
	api.addAuthHeaders(req)
	req.Header.Add("content-type", "application/json")

	res, err := api.client.Do(req)
	if err != nil {
		api.logger.Errorln("Request failed on do ", err)
		statusCode := 1
		if res != nil {
			statusCode = res.StatusCode
		}
		return NewGoHabitsError(err.Error(), statusCode, "")
	}

	body := api.readBody(res)
	api.logger.Debugln(res.Status, req.Method, req.URL, "Body: \n", string(body))

	return api.parseResponse(body, res, responseType)
}

func (api *HabiticaAPI) readBody(res *http.Response) []byte {
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		api.logger.Errorln(err)
	}
	return body
}

func (api *HabiticaAPI) parseResponse(body []byte, res *http.Response, object interface{}) error {
	var hres habiticaResponse
	var err error

	if len(body) > 0 {
		err = json.Unmarshal(body, &hres)
		if err != nil {
			api.logger.Debugln("Unmarshal habitica response failed", string(body))
			api.logger.Errorln(err)
			return NewGoHabitsError("Unmarshal habitica response failed", res.StatusCode, res.Request.URL.EscapedPath())
		}
	}

	if hres.Error != "" || res.StatusCode >= 400 {
		err := parseHabitsServerError(hres, res)
		api.logger.Errorln("Error status from Habitica API", err)
		return err
	}

	if len(body) > 0 {
		api.logger.Debugln("hres.Data", string(hres.Data))
		err = json.Unmarshal(hres.Data, &object)
		if err != nil {
			api.logger.Debugln("Unmarshal data failed", string(body))
			api.logger.Errorln(err)
			return NewGoHabitsError("Unmarshal data failed", 1, res.Request.URL.EscapedPath())
		}
		api.logger.Debugln("Object", object)
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
	errMessage := hres.Message
	for _, errorMessage := range hres.Errors {
		errMessage += "\n" + errorMessage.Message
	}
	return NewGoHabitsError(errMessage, res.StatusCode, res.Request.URL.EscapedPath())
}

func (api *HabiticaAPI) createRequest(method, route string, requestObject interface{}) (*http.Request, error) {
	switch method {
	case http.MethodGet, http.MethodDelete:
		return http.NewRequest(method, api.hostURL+"/v3"+route, nil)
	case http.MethodPost, http.MethodPut:
		data, err := json.Marshal(requestObject)
		if err != nil {
			api.logger.Errorln(err)
			return nil, err
		}
		return http.NewRequest(method, api.hostURL+"/v3"+route, bytes.NewBuffer(data))
	}
	return nil, NewGoHabitsError("Method '"+method+"' not supported", 1, route)
}

func (api *HabiticaAPI) runRequest(method, route string, requestObject interface{}, responseType interface{}) error {
	req, err := api.createRequest(method, route, requestObject)
	if err != nil {
		api.logger.Errorln(err)
		return err
	}

	return api.Do(req, responseType)
}

// Get will return response from the passed in route of Habitica Api.
// It will also return errors in either HTTP Protocol or if status
// code is equal to or above 400.
func (api *HabiticaAPI) Get(route string, responseType interface{}) error {
	return api.runRequest(http.MethodGet, route, nil, responseType)
}

// Post will take in url, request data as a struct, and response as pointer to struct and output
// errors for marshalling either.
func (api *HabiticaAPI) Post(url string, requestObject interface{}, responseObject interface{}) error {
	return api.runRequest(http.MethodPost, url, requestObject, responseObject)
}

// Put will take in url, request data as a struct and repsonse as a struct and output errors for
// marshaling either
func (api *HabiticaAPI) Put(url string, requestObject interface{}, responseObject interface{}) error {
	return api.runRequest(http.MethodPut, url, requestObject, responseObject)
}

// Delete will take an url, and response as a struct and output errors for marshalling either.
func (api *HabiticaAPI) Delete(route string) error {
	return api.runRequest(http.MethodDelete, route, nil, nil)
}
