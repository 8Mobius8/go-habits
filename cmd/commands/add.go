package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	api "github.com/8Mobius8/go-habits/api"
)

// TaskLine is the Perl regular expression `(?m)^\\[ \\] (.*)$`
// `(?m)` sets flags so that `^` `$` match line beginning endings.
var TaskLine = regexp.MustCompile("(?m)^\\[ \\] (.*)$")

// AddTaskServer interface handles adding and getting tasks
// from a habits server
type AddTaskServer interface {
	AddTask(api.Task) (api.Task, error)
	GetTasks(api.TaskType) []api.Task
}

// Add or `go-habits add` command allows habiters to add
// new tasks to their list of todos. When running this command
// the format expected is like follows:
// go-habits a {{ TaskTitle }}
func Add(out io.Writer, server AddTaskServer, args []string, filePath string) error {
	var tasks []api.Task
	var err error

	if filePath != "" {
		tasks, err = ParseTasksFromFile(filePath)
	} else {
		tasks = []api.Task{ParseTask(args)}
	}
	if err != nil {
		return err
	}

	ids := []string{}
	for _, t := range tasks {
		tt, err := server.AddTask(t)
		if err != nil {
			fmt.Fprintln(out, err)
			continue
		}
		ids = append(ids, tt.ID)
	}

	tasks = server.GetTasks(api.TodoType)
	for _, id := range ids {
		printTasks(out, FilterTask(id, tasks))
	}
	return nil
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

// FilterTask will filter a array of tasks by id.
func FilterTask(id string, tasks []api.Task) []api.Task {
	var filtered []api.Task
	for _, task := range tasks {
		if task.ID == id {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// ParseTasksFromFile will return an array of tasks that matches
// `TaskLine` in the file given, line by line.
func ParseTasksFromFile(filePath string) ([]api.Task, error) {
	tasks := []api.Task{}

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []api.Task{}, err
	}

	// Using FindAll to get just the lines that have
	// tasks in them.
	lines := TaskLine.FindAll(dat, -1)
	for _, line := range lines {
		t := api.NewTask(strings.TrimLeft(string(line), "[ ]"), api.TodoType)
		tasks = append(tasks, t)
	}

	return tasks, nil
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
