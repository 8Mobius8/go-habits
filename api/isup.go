package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func ParseResponse(res *http.Response) StatusResponse {

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var isUp StatusResponse

	err = json.Unmarshal(body, &isUp)
	if err != nil {
		log.Fatal(err)
	}

	return isUp
}

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
