package commands

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	api "github.com/8Mobius8/go-habits/api"
)

// DeleteServer interface does everything a `TaskServer` does
// and Deletes tasks.
type DeleteServer interface {
	TasksServer
	DeleteTask(api.Task) error
}

// Remove will remove tasks by order as given in arguments from a `DeleteServer`
func Remove(in io.Reader, out io.Writer, args []string, server DeleteServer, force bool) error {
	pArg, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	t, err := GetTaskAtPosition(server, pArg-1)
	if err != nil {
		return err
	}

	deleteTask := force
	if !force {
		taskYorN := fmt.Sprint(formatTask(t), `[Y\n]?`)
		deleteTask = decideYesOrNo(in, out, taskYorN)
	}

	if deleteTask {
		err = server.DeleteTask(t)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, "Removed tasks:")
		fmt.Fprintln(out, fmt.Sprintf("%d%s %s ", t.Order, "X", t.Title))
	}

	return nil
}

func decideYesOrNo(in io.Reader, out io.Writer, question string) bool {
	answer := getAnswer(in, out, question)
	matched, _ := regexp.MatchString("Y", answer)
	return matched
}

func getAnswer(in io.Reader, out io.Writer, question string) (answer string) {
	fmt.Fprint(out, "Remove?\n", question)
	fmt.Fscan(in, &answer)
	return
}
