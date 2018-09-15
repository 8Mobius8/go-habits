package cmd

import (
	"fmt"

	api "github.com/8Mobius8/go-habits/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List",
	Aliases: []string{"l", "l t"},
	Run: func(cmd *cobra.Command, args []string) {
		apiClient := api.NewHabiticaAPI(nil, viper.GetString("server"))
		todos := apiClient.GetTodos()
		PrintTodos(todos)
	},
}

// PrintTodos ...
func PrintTodos(todos []api.Todo) {
	for _, todo := range todos {
		fmt.Printf("%d %s\n", todo.Order, todo.Title)
	}
}
