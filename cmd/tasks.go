// Copyright © 2018 Mark Odell <m4rk.odell@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	HabitApi "github.com/8Mobius8/go-habits/api"
	"github.com/spf13/cobra"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List your undone tasks",
	Long:  `List your undone tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		api := HabitApi.NewHabiticaApi(nil, "https://habitica.com")

		res, err := api.Tasks()
		if err != nil {
			fmt.Println(err)
		}

		var taskResponse HabitApi.TasksResponse
		api.ParseResponse(res, &taskResponse)

		fmt.Println(HabitApi.GetTasks(taskResponse))
	},
}

func init() {
	rootCmd.AddCommand(tasksCmd)

}
