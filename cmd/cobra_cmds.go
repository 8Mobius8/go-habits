package cmd

import (
	"fmt"
	"os"

	"github.com/8Mobius8/go-habits/cmd/commands"
	"github.com/spf13/cobra"
)

func init() {
	addCmd.Flags().StringVarP(&file, "file", "f", "", "File name to parse for tasks")
	removeCmd.Flags().BoolVarP(&forceRemove, "force", "f", false, "Remove task without confirmation")
	rootCmd.AddCommand(
		addCmd,
		completeCmd,
		removeCmd,
		listCmd,
		loginCmd,
		statusCmd,
		versionCmd,
	)
}

var file string
var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a todo to Habitica",
	Aliases: []string{"a", "a t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.Add(cmd.OutOrStdout(), habitsServer, args, file)
	},
}

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark tasks as complete using this command.",
	Long:  `Mark tasks as complete using this command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.Complete(cmd.OutOrStdout(), habitsServer, args)
	},
}

var forceRemove bool
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove tasks from your list. Does NOT complete them.",
	Long: `Remove tasks from your list. Does NOT complete them.
You will not recieve and awards for removing tasks.`,
	Aliases: []string{"remove todo", "rm t", "rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if forceRemove {
			return commands.RemoveForced(os.Stdin, cmd.OutOrStdout(), args, habitsServer)
		}
		return commands.Remove(os.Stdin, cmd.OutOrStdout(), args, habitsServer)
	},
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List todos",
	Aliases: []string{"l", "l t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.List(cmd.OutOrStdout(), habitsServer, args)
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenicates with Habits server and saves api token in config file.",
	Long: "Authenicates with Habits server and saves api token in config file.\n" +
		"You will need create an account on https://habitica.com",
	Run: func(cmd *cobra.Command, args []string) {
		commands.Login(os.Stdin, cmd.OutOrStdout(), habitsServer, args)
	},
}

var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Check if Habitica api is reachable.",
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return commands.Status(cmd.OutOrStdout(), habitsServer)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version of go-habits.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), rootCmd.Use+" version "+rootCmd.Version)
	},
}
