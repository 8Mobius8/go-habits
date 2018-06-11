package api

import (
	"log"
)

// Authenticate will return Habitica ID and APIToken with given username
// and password.
func (api *HabiticaAPI) Authenticate(user string, password string) UserToken {
	var creds userCredentials
	creds.Username = user
	creds.Password = password

	var resp UserToken
	err := api.Post("/user/auth/local/login", creds, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	if api.userAuth.ID != resp.ID {
		api.userAuth.ID = resp.ID
	}
	if api.userAuth.APIToken != resp.APIToken {
		api.userAuth.APIToken = resp.APIToken
	}

	return resp
}

type userCredentials struct {
	Username string
	Password string
}

// UserToken contains user ID and Token to make API calls.
type UserToken struct {
	ID       string
	APIToken string
}
