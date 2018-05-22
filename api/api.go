package api

import (
	"github.com/8Mobius8/go-habits/api/status"
)

// Status will return response from `/status` route of Habitica Api.
// It will also return errors in either HTTP Protocol or if status
// code is equal to or above 400.
func (api *HabiticaAPI) Status() (status.StatusResponse, error) {
	body, err := api.Get("/status")

	var status status.StatusResponse
	if err == nil {
		api.ParseResponse(body, &status)
	}

	return status, err
}

// Tasks will return response from `/task` route of Habitica Api.
// It will also return errors in either HTTP Protocol or if status
// code is equal to or above 400.
func (api *HabiticaAPI) Tasks() ([]byte, error) {
	return api.Get("/user")
}
