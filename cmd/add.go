package cmd

import (
	"fmt"
	"os"
	"strings"

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

	t := ParseTask(args)

	t, err := client.AddTask(t)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tasks := client.GetTasks(api.TodoType)
	printTasks(filterTask(t.ID, tasks))
}

func filterTask(id string, tasks []api.Task) []api.Task {
	var filtered []api.Task
	for _, task := range tasks {
		if task.ID == id {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

func parseTaskTitle(args []string) string {
	title := args[0]
	for _, arg := range args[1:] {
		title += " " + arg
	}
	return title
}

// ParseTask parses an api.Task from []string
func ParseTask(args []string) api.Task {
	var titleArgs, tagsArgs []string

	for _, arg := range args {
		if arg[0] == '#' {
			tagsArgs = append(tagsArgs, arg)
		} else {
			titleArgs = append(titleArgs, arg)
		}
	}

	t := api.NewTask(parseTaskTitle(titleArgs), api.TodoType)
	t.Tags = parseTags(tagsArgs)
	return t
}

func parseTags(args []string) []string {
	for i, arg := range args {
		args[i] = strings.Split(arg, "#")[1]
	}
	return args
}
