package commands

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	api "github.com/8Mobius8/go-habits/api"
)

// TasksServer interface is able to Get, Add, and Score Up tasks
type TasksServer interface {
	GetTasks(api.TaskType) []api.Task
	AddTask(api.Task) (api.Task, error)
	ScoreTaskUp(api.Task) (api.ScoreUpDelta, error)
}

// Complete or `go-habits complete` allows habiters to complete
// todos on their existing list. `complete` task in a interger
// number representing the position of the todo they would like
// to complete.
func Complete(out io.Writer, ts TasksServer, args []string) error {
	parseArg, err := strconv.Atoi(args[0])
	t, err := GetTaskAtPosition(ts, parseArg-1)
	if err != nil {
		if strings.Contains(err.Error(), "no tasks") {
			fmt.Fprintln(out, "You have no tasks.")
			fmt.Fprintln(out, "Create tasks before trying to complete them.")
		} else if strings.Contains(err.Error(), "bad index") {
			fmt.Fprintf(out, "There is no task at %d", parseArg)
		}
		return api.NewGoHabitsError(err.Error(), 1, "")
	}

	delta, err := ts.ScoreTaskUp(t)
	if err != nil {
		return err
	}
	t.Completed = true
	fmt.Fprintln(out, formatTask(t))
	fmt.Fprintf(out, "MP: %f GP: %f XP: %f\n", delta.Mp, delta.Gp, delta.Exp)
	return nil
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
