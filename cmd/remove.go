package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	api "github.com/8Mobius8/go-habits/api"
	"github.com/spf13/cobra"
)

var (
	ForceRemove bool
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove tasks from your list. Does NOT complete them.",
	Long: `Remove tasks from your list. Does NOT complete them.
You will not recieve and awards for removing tasks.`,
	Aliases: []string{"remove todo", "rm t", "rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return Remove(os.Stdin, cmd.OutOrStdout(), args, habitsServer)
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&ForceRemove, "force", "f", false, "Remove task without confirmation")
	rootCmd.AddCommand(removeCmd)
}

// DeleteServer interface does everything a `TaskServer` does
// and Deletes tasks.
type DeleteServer interface {
	TasksServer
	DeleteTask(api.Task) error
}

// Remove will remove tasks by order as given in arguments from a `DeleteServer`
func Remove(in io.Reader, ino io.Writer, args []string, server DeleteServer) error {
	pArg, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	t, err := GetTaskAtPosition(server, pArg-1)
	if err != nil {
		return err
	}

	if !ForceRemove {
		fmt.Fprintln(ino, "Remove?")
		fmt.Fprintf(ino, "%s [Y\\n]?", formatTask(t))
		answer := "n"
		fmt.Fscanln(in, &answer)
	}

	err = server.DeleteTask(t)
	if err != nil {
		return err
	}

	fmt.Fprintln(ino, "Removed tasks:")
	fmt.Fprintln(ino, fmt.Sprintf("%d%s %s ", t.Order, "X", t.Title))
	return nil
}
