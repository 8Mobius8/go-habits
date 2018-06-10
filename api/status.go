package api

// Status will return response from `/status` route of Habitica Api.
// It will also return errors in either HTTP Protocol or if status
// code is equal to or above 400.
func (api *HabiticaAPI) Status() (Status, error) {
	var res Status
	err := api.Get("/status", &res)

	return res, err
}

// Status is a string that is usualy 'up' when Habitica API is full available
type Status struct {
	Status string
}
