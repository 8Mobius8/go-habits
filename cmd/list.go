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
	tasks := apiClient.GetTasks(api.TodoType)
	printTasks(tasks)
}

func printTasks(tasks []api.Task) {
	for _, task := range tasks {
		fmt.Println(formatTask(task))
	}
}

func formatTask(t api.Task) string {
	completedString := "[ ]"
	if t.Completed {
		completedString = "[X]"
	}
	s := fmt.Sprintf("%d%s %s ", t.Order, completedString, t.Title)
	for _, tag := range t.Tags {
		s += fmt.Sprintf(" #%s", tag)
	}
	return s
}
