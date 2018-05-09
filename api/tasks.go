package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type TasksResponse struct {
	Success bool
	Data    []Task
}

type Task struct {
	Text string
}

func GetTasks(resp TasksResponse) []Task {

	return resp.Data
}

func ParseTaskResponse(res *http.Response) TasksResponse {

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	var tasks TasksResponse

	err = json.Unmarshal(body, &tasks)
	if err != nil {
		log.Fatal(err)
	}

	return tasks
}
