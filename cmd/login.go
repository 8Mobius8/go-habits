package cmd

import (
	"fmt"
	"io/ioutil"

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
	Run: Login,
}

// Login or `go-habits login` allows habiters to logon to a
// habitica server and save their api id and token to a config
// file. The file must be previously created.
func Login(cmd *cobra.Command, args []string) {
	client := habitsServer
	user := scanForUserCreds()
	creds := client.Authenticate(user.UserName, user.Password)
	cmd.Println("Login Successful <3")
	setAuthConfig(creds)
	err := saveToConfigFile()
	if err != nil {
		cmd.Println(err)
	}
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

func saveToConfigFile() error {
	if viper.ConfigFileUsed() == "" {
		fmt.Println("Didn't find config file. Creating a new config file at " + defaultGoHabitsConfigPath)
		touchConfigFile(defaultGoHabitsConfigPath)
		viper.SetConfigFile(defaultGoHabitsConfigPath)
	}
	fmt.Printf("Updating config at %s\n", viper.ConfigFileUsed())
	return viper.WriteConfig()
}

func touchConfigFile(configPath string) error {
	return ioutil.WriteFile(configPath, []byte(`test: hai`), 0644)
}
