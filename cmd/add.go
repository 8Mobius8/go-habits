package cmd

import (
	api "github.com/8Mobius8/go-habits/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a todo to Habitica",
	Aliases: []string{"a", "a t"},
	Run:     Add,
}

// Add or `go-habits add` command allows habiters to add
// new tasks to their list of todos. When running this command
// the format expected is like follows:
// go-habits a {{ TaskTitle }}
func Add(cmd *cobra.Command, args []string) {
	client := habitsServer

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
