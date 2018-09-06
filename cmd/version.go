package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version of go-habits.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(rootCmd.Use + " version " + rootCmd.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
