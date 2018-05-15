package api

type StatusResponse struct {
	Success bool
	Data    struct {
		Status string
	}
}

func HabiticaStatusMessage(resp StatusResponse) string {
	if resp.Data.Status != "up" {
		return ":( Habitica is unreachable."
	}

	return "Habitica is reachable, GO catch all those pets!"
}
