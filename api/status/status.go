package status

type StatusResponse struct {
	Success bool
	Data    struct {
		Status string
	}
}
