package cmd

import (
	"io"
	"strings"

	api "github.com/8Mobius8/go-habits/api"
	log "github.com/amoghe/distillog"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a todo to Habitica",
	Aliases: []string{"a", "a t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return Add(cmd.OutOrStdout(), habitsServer, args)
	},
}

type AddTaskServer interface {
	AddTask(api.Task) (api.Task, error)
	GetTasks(api.TaskType) []api.Task
}

// Add or `go-habits add` command allows habiters to add
// new tasks to their list of todos. When running this command
// the format expected is like follows:
// go-habits a {{ TaskTitle }}
func Add(out io.Writer, server AddTaskServer, args []string) error {
	t := ParseTask(args)

	t, err := server.AddTask(t)
	if err != nil {
		log.Errorln(err)
		return err
	}

	tasks := server.GetTasks(api.TodoType)
	printTasks(out, filterTask(t.ID, tasks))
	return nil
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

func parseTaskTitle(args []string) string {
	title := args[0]
	for _, arg := range args[1:] {
		title += " " + arg
	}
	return title
}

func parseTags(args []string) []string {
	for i, arg := range args {
		args[i] = strings.Split(arg, "#")[1]
	}
	return args
}
