package cmd

import (
	"fmt"
	"os"

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
	Run:     Status,
}

// Status or `go-habits status` allows habiters to check if the
// habitica server api is up and running.
func Status(cmd *cobra.Command, args []string) {
	fmt.Println("Using " + viper.GetString("server") + " as api server")
	client := habitsServer
	res, err := client.GetServerStatus()
	if err != nil {
		fmt.Println(err)
		fmt.Println(StatusMessage(res))
		os.Exit(5)
	}

	fmt.Println(StatusMessage(res))
}

// StatusMessage returns text based on Status message
func StatusMessage(resp api.Status) string {
	if resp.Status != "up" {
		return ":( Habitica is unreachable."
	}

	return "Habitica is reachable, GO catch all those pets!"
}
