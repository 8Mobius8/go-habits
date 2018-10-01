package cmd

import (
	"fmt"

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
	Long:  `Authenicates with Habits server and saves api token in config file.`,
	Run:   Login,
}

// Login or `go-habits login` allows habiters to logon to a
// habitica server and save their api id and token to a config
// file. The file must be previously created.
func Login(cmd *cobra.Command, args []string) {
	client := habitsServer
	user := scanForUserCreds()
	creds := client.Authenticate(user.UserName, user.Password)
	setAuthConfig(creds)
	saveToConfigFile()
}

func scanForUserCreds() api.UserToken {
	var user api.UserToken
	fmt.Print("Username:")
	fmt.Scanln(&user.UserName)
	fmt.Print("Password:")
	fmt.Scanln(&user.Password)
	return user
}

func setAuthConfig(creds api.UserToken) {
	viper.Set("auth.local.id", creds.ID)
	viper.Set("auth.local.apitoken", creds.APIToken)
}

func saveToConfigFile() {
	fmt.Println(viper.ConfigFileUsed())
	if viper.ConfigFileUsed() == "" {
		fmt.Println("Didn't find config file. Create one at ~/.go-habits.yml to save api key for later use.")
	} else {
		fmt.Printf("Updating config at %s\n", viper.ConfigFileUsed())
		viper.WriteConfig()
	}
}
