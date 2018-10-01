package cmd

import (
	"fmt"

	api "github.com/8Mobius8/go-habits/api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List todos",
	Aliases: []string{"l", "l t"},
	Run:     List,
}

// List or `go-habits list` command allows habiters to see
// their list of todos currently needing to be completed.
func List(cmd *cobra.Command, args []string) {
	apiClient := habitsServer
	todos := apiClient.GetTodos()
	printTodos(todos)
}

func printTodos(todos []api.Todo) {
	for _, todo := range todos {
		fmt.Println(formatTodo(todo))
	}
}

func formatTodo(t api.Todo) string {
	s := fmt.Sprintf("%d[ ] %s ", t.Order, t.Title)
	for _, tag := range t.Tags {
		s += fmt.Sprintf(" #%s", tag)
	}
	return s
}
