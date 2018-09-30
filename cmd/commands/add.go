package commands

import (
	api "github.com/8Mobius8/go-habits/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Add or `go-habits add` command allows habiters to add
// new tasks to their list of todos. When running this command
// the format expected is like follows:
// go-habits a {{ TaskTitle }}
func Add(cmd *cobra.Command, args []string) {
	client := api.NewHabiticaAPI(nil, viper.GetString("server"))
	client.UpdateUserAuth(getAuthConfig())

	t := api.Todo{}
	t.Title = parseTodoTitle(args)
	t = client.AddTodo(t)

	todos := client.GetTodos()
	printTodos(filterTodo(t.ID, todos))
}

func filterTodo(id string, todos []api.Todo) []api.Todo {
	var filtered []api.Todo
	for _, todo := range todos {
		if todo.ID == id {
			filtered = append(filtered, todo)
		}
	}
	return filtered
}

func parseTodoTitle(args []string) string {
	title := args[0]
	for _, arg := range args[1:] {
		title += " " + arg
	}
	return title
}
