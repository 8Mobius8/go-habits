package cmd

import (
	api "github.com/8Mobius8/go-habits/api"
	log "github.com/8Mobius8/go-habits/log"
	"github.com/8Mobius8/go-habits/log/levels"
	dlog "github.com/amoghe/distillog"
	"github.com/spf13/viper"
)

var cfgFile string

var habitsServer *api.HabiticaAPI
var habitsServerURL string

// logLevel is set from flag on root command in root.go
var logLevel string

// setupAPIClient creates a new API client for running go-habits cmd.
// Uses logLevel variable to determine how to configure API client
// to log. Additionaly will fetch users config and set it in the API client
func setupAPIClient() {
	clientLogger := dlog.NewNullLogger("go-habits")
	if logLevel != "" {
		switch logLevel {
		case "ERROR", "error", "1":
			clientLogger = log.NewLevelLogger(levels.Error, dlog.NewStdoutLogger("go-habits"))
		case "WARN", "warn", "2":
			clientLogger = log.NewLevelLogger(levels.Warning, dlog.NewStdoutLogger("go-habits"))
		case "INFO", "info", "3":
			clientLogger = log.NewLevelLogger(levels.Info, dlog.NewStdoutLogger("go-habits"))
		case "DEBUG", "debug", "4":
			clientLogger = log.NewLevelLogger(levels.Debug, dlog.NewStdoutLogger("go-habits"))
		}
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
