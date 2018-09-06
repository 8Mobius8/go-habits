package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	HabitApi "github.com/8Mobius8/go-habits/api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(isupCmd)
}

var isupCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if Habitica api is reachable.",
	Run: func(cmd *cobra.Command, args []string) {
		api := HabitApi.NewHabiticaAPI(nil, viper.GetString("server"))

		res, err := api.Status()
		if err != nil {
			fmt.Println(err)
			os.Exit(5)
		}

		fmt.Println(StatusMessage(res))
	},
}

// StatusMessage returns text based on Status message
func StatusMessage(resp HabitApi.Status) string {
	if resp.Status != "up" {
		return ":( Habitica is unreachable."
	}

	return "Habitica is reachable, GO catch all those pets!"
}
