// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

		var user HabitApi.UserToken
		fmt.Print("Username:")
		fmt.Scanln(&user.UserName)
		fmt.Print("Password:")
		fmt.Scanln(&user.Password)
		creds := api.Authenticate(user.UserName, user.Password)

		viper.Set("apiauth.id", creds.ID)
		viper.Set("apiauth.apitoken", creds.APIToken)
		fmt.Println(viper.ConfigFileUsed())
		if viper.ConfigFileUsed() == "" {
			fmt.Println("Didn't find config file. Create one at ~/.go-habits.yaml to save api key for later use")
		} else {
			fmt.Printf("Updating config at %s\n", viper.ConfigFileUsed())
			viper.WriteConfig()
		}
	},
}
