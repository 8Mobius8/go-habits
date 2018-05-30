package api

// Status will return response from `/status` route of Habitica Api.
// It will also return errors in either HTTP Protocol or if status
// code is equal to or above 400.
func (api *HabiticaAPI) Status() (StatusResponse, error) {
	body, err := api.Get("/status")

	var status StatusResponse
	if err == nil {
		api.ParseResponse(body, &status)
	}

	return status, err
}

type StatusResponse struct {
	Success bool
	Data    struct {
		Status string
	}
}
