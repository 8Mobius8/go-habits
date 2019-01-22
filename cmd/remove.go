package cmd

import (
	"fmt"
	"io"
	"strconv"

	api "github.com/8Mobius8/go-habits/api"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove tasks from your list. Does NOT complete them.",
	Long: `Remove tasks from your list. Does NOT complete them.
You will not recieve and awards for removing tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return Remove(cmd.OutOrStdout(), args, habitsServer)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

type DeleteServer interface {
	TasksServer
	DeleteTask(api.Task) error
}

func Remove(ino io.Writer, args []string, server DeleteServer) error {
	pArg, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	t, err := GetTaskAtPosition(server, pArg-1)
	if err != nil {
		return err
	}
	err = server.DeleteTask(t)
	if err != nil {
		return err
	}

	fmt.Fprintln(ino, "Removed tasks:")
	fmt.Fprintln(ino, fmt.Sprintf("%d%s %s ", t.Order, "X", t.Title))
	return nil
}
