package cmd

import (
	"fmt"
	"strconv"
	"strings"

	api "github.com/8Mobius8/go-habits/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completeCmd)
}

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark tasks as complete using this command.",
	Long:  `Mark tasks as complete using this command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := Complete(habitsServer, args)
		cmd.Println(out)
		return err
	},
}

// TasksServer interface is able to Get, Add, and Score Up tasks
type TasksServer interface {
	GetTasks(api.TaskType) []api.Task
	AddTask(api.Task) (api.Task, error)
	ScoreTaskUp(api.Task) error
}

// Complete or `go-habits complete` allows habiters to complete
// todos on their existing list. `complete` task in a interger
// number representing the position of the todo they would like
// to complete.
func Complete(ts TasksServer, args []string) (string, error) {
	out := ""
	parseArg, err := strconv.Atoi(args[0])
	t, err := GetTaskAtPosition(ts, parseArg-1)
	if err != nil {
		if strings.Contains(err.Error(), "no tasks") {
			out += "You have no tasks.\n"
			out += "Create tasks before trying to complete them."
		} else if strings.Contains(err.Error(), "bad index") {
			out += fmt.Sprintf("There is no task at %d", parseArg)
		}
		return out, api.NewGoHabitsError(err.Error(), 1, "")
	}

	err = ts.ScoreTaskUp(t)
	if err != nil {
		err = api.NewGoHabitsError(err.Error(), 1, "")
		return err.Error(), err
	}
	t.Completed = true
	return formatTask(t), nil
}

// GetTaskAtPosition returns todo from task server at a given position or
// index
func GetTaskAtPosition(ts TasksServer, p int) (api.Task, error) {
	tasks := ts.GetTasks(api.TodoType)
	if len(tasks) == 0 {
		return api.Task{}, fmt.Errorf("no tasks")
	}
	if p < 0 || p >= len(tasks) {
		return api.Task{}, fmt.Errorf("%d is a bad index", p+1)
	}
	return tasks[p], nil
}
