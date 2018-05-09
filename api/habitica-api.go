package api

import (
	"log"
	"net/http"
)

func GetHabiticaAPIStatus() *http.Response {
	res, err := http.Get("https://habitica.com/api/v3/status")
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func GetHabiticaAPITasks() *http.Response {
	res, err := http.Get("https://habitica.com/api/v3/tasks/user")
	if err != nil {
		log.Fatal(err)
	}

	return res
}
