package cmd

import (
	"fmt"

	HabitApi "github.com/8Mobius8/go-habits/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenicates with Habits server and saves api token in config file.",
	Long:  `Authenicates with Habits server and saves api token in config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		api := HabitApi.NewHabiticaAPI(nil, viper.GetString("server"))
		user := scanForUserCreds()
		creds := api.Authenticate(user.UserName, user.Password)
		setAuthConfig(creds)
		saveToConfigFile()
	},
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
		fmt.Println("Didn't find config file. Create one at ~/.go-habits.yaml to save api key for later use.")
	} else {
		fmt.Printf("Updating config at %s\n", viper.ConfigFileUsed())
		viper.WriteConfig()
	}
}
