package api

// Status will return response from `/status` route of Habitica Api.
// It will also return errors in either HTTP Protocol or if status
// code is equal to or above 400.
func (api *HabiticaAPI) Status() (Status, error) {
	body, err := api.Get("/status")

	var status statusResponse
	if err == nil {
		api.ParseResponse(body, &status)
	}

	return status.Data.Status, err
}

// Status is string that is usually 'up' when Habitica API is fully available
type Status string

type statusResponse struct {
	Success bool
	Data    struct {
		Status Status
	}
}
