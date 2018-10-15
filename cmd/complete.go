package cmd

import (
	"fmt"
	"os"
	"strconv"

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
	Run:   Complete,
}

// Complete or `go-habits complete` allows habiters to complete
// todos on their existing list. `complete` task in a interger
// number representing the position of the todo they would like
// to complete.
func Complete(cmd *cobra.Command, args []string) {
	client := habitsServer
	tasks := client.GetTasks(api.TodoType)
	if len(tasks) == 0 {
		fmt.Println("You have no tasks.")
		fmt.Println("Create tasks before trying to complete them.")
		os.Exit(1)
		return
	}

	parseArg, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	order := parseArg - 1
	if order < 0 || order >= len(tasks) {
		fmt.Printf("There is no task at %d\n", order+1)
		os.Exit(1)
		return
	}

	err = client.ScoreTaskUp(tasks[order])
	if err != nil {
		fmt.Println(err)
	}
	printTasks(tasks[order : order+1])
}
