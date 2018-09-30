package cmd

import (
	api "github.com/8Mobius8/go-habits/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a todo to Habitica",
	Aliases: []string{"a", "a t"},
	Run: func(cmd *cobra.Command, args []string) {
		client := api.NewHabiticaAPI(nil, viper.GetString("server"))
		client.UpdateUserAuth(getAuthConfig())

		t := api.Todo{}
		t.Title = parseTodoTitle(args)
		t = client.AddTodo(t)

		todos := client.GetTodos()
		PrintTodo(filterTodo(t.ID, todos))
	},
}

func filterTodo(id string, todos []api.Todo) api.Todo {
	for _, todo := range todos {
		if todo.ID == id {
			return todo
		}
	}
	return api.Todo{}
}

func parseTodoTitle(args []string) string {
	title := ""
	for _, arg := range args {
		title += arg
	}
	return title
}
