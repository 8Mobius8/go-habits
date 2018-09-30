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
	Short:   "List todos",
	Aliases: []string{"l", "l t"},
	Run: func(cmd *cobra.Command, args []string) {
		apiClient := api.NewHabiticaAPI(nil, viper.GetString("server"))
		apiClient.UpdateUserAuth(getAuthConfig())
		todos := apiClient.GetTodos()
		PrintTodos(todos)
	},
}

func getAuthConfig() api.UserToken {
	id := viper.GetString("auth.local.id")
	token := viper.GetString("auth.local.apitoken")
	creds := api.UserToken{}
	creds.ID = id
	creds.APIToken = token
	return creds
}

func PrintTodos(todos []api.Todo) {
	for _, todo := range todos {
		PrintTodo(todo)
	}
}

func PrintTodo(todo api.Todo) {
	fmt.Printf("%d[ ] %s", todo.Order, todo.Title)
	for _, tag := range todo.Tags {
		fmt.Printf(" #%s", tag)
	}
	fmt.Println()
}
