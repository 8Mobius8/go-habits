package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	api "github.com/8Mobius8/go-habits/api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Check if Habitica api is reachable.",
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := Status()
		cmd.Println(out)
		return err
	},
}

// Status or `go-habits status` allows habiters to check if the
// habitica server api is up and running.
func Status() (string, error) {
	output := fmt.Sprintln("Using " + viper.GetString("server") + " as api server")
	res, err := habitsServer.GetServerStatus()
	if err != nil {
		output += fmt.Sprintln(err)
		ghe, ok := err.(*api.GoHabitsError)
		if ok {
			ghe.StatusCode = 5
			err = ghe
		}
	}

	output += fmt.Sprintln(StatusMessage(res))
	return output, err
}

// StatusMessage returns text based on Status message
func StatusMessage(resp api.Status) string {
	if resp.Status != "up" {
		return ":( Habitica is unreachable."
	}

	return "Habitica is reachable, GO catch all those pets!"
}
