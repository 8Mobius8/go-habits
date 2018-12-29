package cmd

import (
	api "github.com/8Mobius8/go-habits/api"
	log "github.com/amoghe/distillog"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var defaultGoHabitsConfigPath string
var habitsServer *api.HabiticaAPI
var habitsServerURL string
var logLevel string

// init will gather system info to run go-habits cmd on.
func init() {
	userHomePath, _ := homedir.Dir()
	defaultGoHabitsConfigPath = userHomePath + "/.go-habits.yml"
}

// setupAPIClient creates a new API client for running go-habits cmd.
// Uses logLevel variable to determine how to configure API client
// to log. Additionaly will fetch users config and set it in the API client
func setupAPIClient() {
	clientLogger := log.NewNullLogger("go-habits")
	switch logLevel {
	case "ERROR":
		clientLogger = log.NewStderrLogger("go-habits")
	}
	habitsServer = api.NewHabiticaAPI(nil, viper.GetString("server"), clientLogger)
	habitsServer.UpdateUserAuth(getAuthConfig())
}

// getAuthConfig will use viper for configuration of users' id and token
// to habitica.
func getAuthConfig() api.UserToken {
	id := viper.GetString("auth.local.id")
	token := viper.GetString("auth.local.apitoken")
	creds := api.UserToken{}
	creds.ID = id
	creds.APIToken = token
	return creds
}

// handleRootError will check if it's ghe error and return the exit code
// accordingly, otherwise as long as error is not nil will return 1.
func handleRootError(err error) int {
	exitCode := 1
	ghe, ok := err.(*api.GoHabitsError)
	if ok {
		exitCode = ghe.StatusCode
	}
	return exitCode
}
