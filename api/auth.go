package api

import (
	"log"
	"net/http"
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

// UpdateUserAuth takes a UserToken object and updates apiclient's user's id and token.
func (api *HabiticaAPI) UpdateUserAuth(creds UserToken) {
	if api.userAuth.ID != creds.ID {
		api.userAuth.ID = creds.ID
	}
	if api.userAuth.APIToken != creds.APIToken {
		api.userAuth.APIToken = creds.APIToken
	}
}

func (api *HabiticaAPI) addAuthHeaders(req *http.Request) {
	if api.userAuth.APIToken != "" {
		req.Header.Add("x-api-key", api.userAuth.APIToken)
	}
	if api.userAuth.ID != "" {
		req.Header.Add("x-api-user", api.userAuth.ID)
	}
}

// UserToken contains user ID and Token to make API calls.
type UserToken struct {
	ID       string `json:"id"`
	APIToken string `json:"apitoken"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
