package api

import (
	"log"
)

// Authenticate will return Habitica ID and APIToken with given username
// and password.
func (api *HabiticaAPI) Authenticate(user string, password string) UserToken {
	creds := struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}{
		user,
		password,
	}

	var resp UserToken
	err := api.Post("/user/auth/local/login", creds, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	api.UpdateUserAuth(resp)

	return resp
}

func (api *HabiticaAPI) UpdateUserAuth(creds UserToken) {
	if api.userAuth.ID != creds.ID {
		api.userAuth.ID = creds.ID
	}
	if api.userAuth.APIToken != creds.APIToken {
		api.userAuth.APIToken = creds.APIToken
	}
}

// UserToken contains user ID and Token to make API calls.
type UserToken struct {
	ID       string `json:"id"`
	APIToken string `json:"apitoken"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
