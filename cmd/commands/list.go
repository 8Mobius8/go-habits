package commands

import (
	"fmt"

	api "github.com/8Mobius8/go-habits/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// List or `go-habits list` command allows habiters to see
// their list of todos currently needing to be completed.
func List(cmd *cobra.Command, args []string) {
	apiClient := api.NewHabiticaAPI(nil, viper.GetString("server"))
	apiClient.UpdateUserAuth(getAuthConfig())
	todos := apiClient.GetTodos()
	printTodos(todos)
}

func getAuthConfig() api.UserToken {
	id := viper.GetString("auth.local.id")
	token := viper.GetString("auth.local.apitoken")
	creds := api.UserToken{}
	creds.ID = id
	creds.APIToken = token
	return creds
}

func printTodos(todos []api.Todo) {
	for _, todo := range todos {
		printTodo(todo)
	}
}

func printTodo(todo api.Todo) {
	fmt.Printf("%d[ ] %s", todo.Order, todo.Title)
	for _, tag := range todo.Tags {
		fmt.Printf(" #%s", tag)
	}
	fmt.Println()
}
