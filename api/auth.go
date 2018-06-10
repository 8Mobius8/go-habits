package api

import "log"

// Authenticate will return Habitica ID and APIToken with given username
// and password.
func (api *HabiticaAPI) Authenticate(user string, password string) UserToken {
	var resp UserTokenResponse
	var creds userCredentials
	creds.Username = user
	creds.Password = password

	err := api.Post("/user/auth/local/login", creds, &resp)
	if err != nil {
		log.Fatalln(err)
	}

	return resp.Data
}

type userCredentials struct {
	Username string
	Password string
}

type UserTokenResponse struct {
	Data UserToken
}

// UserToken contains user ID and Token to make API calls.
type UserToken struct {
	ID       string
	APIToken string
}
