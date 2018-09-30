package commands

import (
	"fmt"

	HabitApi "github.com/8Mobius8/go-habits/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Login or `go-habits login` allows habiters to logon to a
// habitica server and save their api id and token to a config
// file. The file must be previously created.
func Login(cmd *cobra.Command, args []string) {
	api := HabitApi.NewHabiticaAPI(nil, viper.GetString("server"))
	user := scanForUserCreds()
	creds := api.Authenticate(user.UserName, user.Password)
	setAuthConfig(creds)
	saveToConfigFile()
}

func scanForUserCreds() HabitApi.UserToken {
	var user HabitApi.UserToken
	fmt.Print("Username:")
	fmt.Scanln(&user.UserName)
	fmt.Print("Password:")
	fmt.Scanln(&user.Password)
	return user
}

func setAuthConfig(creds HabitApi.UserToken) {
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
