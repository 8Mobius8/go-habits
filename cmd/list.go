package cmd

import (
	"fmt"
	"io"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		return List(cmd.OutOrStdout(), habitsServer, args)
	},
}

// ListServer interface GetTasks will return []Task
type ListServer interface {
	GetTasks(api.TaskType) []api.Task
}

// List or `go-habits list` command allows habiters to see
// their list of todos currently needing to be completed.
func List(out io.Writer, server ListServer, args []string) error {
	tasks := server.GetTasks(api.TodoType)
	printTasks(out, tasks)
	return nil
}

func printTasks(out io.Writer, tasks []api.Task) {
	for _, task := range tasks {
		fmt.Fprintln(out, formatTask(task))
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
