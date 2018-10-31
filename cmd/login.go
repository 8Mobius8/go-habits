package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	api "github.com/8Mobius8/go-habits/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenicates with Habits server and saves api token in config file.",
	Long: "Authenicates with Habits server and saves api token in config file.\n" +
		"You will need create an account on https://habitica.com",
	Run: func(cmd *cobra.Command, args []string) {
		Login(os.Stdin, cmd.OutOrStdout(), habitsServer, args)
	},
}

type authServer interface {
	Authenticate(string, string) api.UserToken
}

// Login or `go-habits login` allows habiters to logon to a
// habitica server and save their api id and token to a config
// file. The file must be previously created.
func Login(in io.Reader, out io.Writer, as authServer, args []string) {
	user := scanForUserCreds(in, out)
	creds := as.Authenticate(user.UserName, user.Password)
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
