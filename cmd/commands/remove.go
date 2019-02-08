package commands

import (
	"fmt"
	"io"
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
func Remove(in io.Reader, ino io.Writer, args []string, server DeleteServer, force bool) error {
	pArg, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	t, err := GetTaskAtPosition(server, pArg-1)
	if err != nil {
		return err
	}

	answer := "n"
	if !force {
		fmt.Fprintln(ino, "Remove?")
		fmt.Fprintf(ino, "%s [Y\\n]?", formatTask(t))
		fmt.Fscanln(in, &answer)
	}

	if answer == "Y" {
		err = server.DeleteTask(t)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(ino, "Removed tasks:")
	fmt.Fprintln(ino, fmt.Sprintf("%d%s %s ", t.Order, "X", t.Title))
	return nil
}
