package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	api "github.com/8Mobius8/go-habits/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var defaultGoHabitsConfigPath string

// init will gather system info to run go-habits cmd on.
func init() {
	userHomePath, _ := homedir.Dir()
	defaultGoHabitsConfigPath = userHomePath + "/.go-habits.yml"
}

// Authenticator small interface for Authenticating users by password
// and username to return a Token
type Authenticator interface {
	Authenticate(string, string) (api.UserToken, error)
}

// Login or `go-habits login` allows habiters to logon to a
// habitica server and save their api id and token to a config
// file. The file must be previously created.
func Login(in io.Reader, out io.Writer, as Authenticator, args []string) {
	user := scanForUserCreds(in, out)
	creds, _ := as.Authenticate(user.UserName, user.Password)
	fmt.Fprintln(out, "Login Successful <3")
	setAuthConfig(creds)
	saveToConfigFile(out)
}

func scanForUserCreds(in io.Reader, out io.Writer) api.UserToken {
	var user api.UserToken
	fmt.Fprint(out, "Username:")
	fmt.Fscanln(in, &user.UserName)
	fmt.Fprint(out, "Password:")
	fmt.Fscanln(in, &user.Password)
	return user
}

func setAuthConfig(creds api.UserToken) {
	viper.Set("auth.local.id", creds.ID)
	viper.Set("auth.local.apitoken", creds.APIToken)
}

func saveToConfigFile(out io.Writer) {
	if viper.ConfigFileUsed() == "" {
		viper.SetConfigFile(defaultGoHabitsConfigPath)
	}

	_, err := os.Stat(viper.ConfigFileUsed())
	viper.WriteConfig()
	if !os.IsNotExist(err) {
		fmt.Fprintf(out, "Overridden config at %s\n", viper.ConfigFileUsed())
	}
	if os.IsNotExist(err) {
		fmt.Fprintln(out, "Didn't find config file.")
		fmt.Fprintf(out, "Created a new config file at %s\n", viper.ConfigFileUsed())
	}
}

func touchConfigFile(configPath string) error {
	return ioutil.WriteFile(configPath, []byte(`test: hai`), 0644)
}
