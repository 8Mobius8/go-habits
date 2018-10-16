package cmd

import (
	"io"

	"github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/spf13/cobra"
)

var _ = Describe("Complete cmd", func() {
	Describe("GetTaskAtPosition", func() {
		It("will return an error when no tasks exist", func() {
			ts := &MockTaskServer{}
			_, err := GetTaskAtPosition(ts, 1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no tasks"))
		})
		It("will return an error when position doesn't exist", func() {
			ts := &MockTaskServer{}
			ts.GetTasksOut = []api.Task{
				{Title: "Task 1"},
			}
			_, err := GetTaskAtPosition(ts, 1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("bad index"))
		})
		It("will return an error when position is negitive", func() {
			ts := &MockTaskServer{}
			ts.GetTasksOut = []api.Task{
				{Title: "Task 1"},
			}
			_, err := GetTaskAtPosition(ts, -1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("bad index"))
		})
		It("will return the task at the right position.", func() {
			ts := &MockTaskServer{}
			ts.GetTasksOut = []api.Task{
				{Title: "Task 1", Order: 1},
				{Title: "Task 2", Order: 2},
				{Title: "Task 3", Order: 3},
			}
			t, err := GetTaskAtPosition(ts, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(t.Title).Should(BeEquivalentTo("Task 2"))
			Expect(t.Order).Should(BeEquivalentTo(2))
		})
	})
})

type MockTaskServer struct {
	GetTasksOut []api.Task
	AddTaskOut  api.Task
	ErrorOut    error
}

func (ts *MockTaskServer) GetTasks(tt api.TaskType) []api.Task {
	return ts.GetTasksOut
}
func (ts *MockTaskServer) AddTask(t api.Task) (api.Task, error) {
	return ts.AddTaskOut, ts.ErrorOut
}
func (ts *MockTaskServer) ScoreTaskUp(t api.Task) error {
	return ts.ErrorOut
}

var _ = Describe("Complete integration", func() {
	Describe("Complete command", func() {
		var buf *gbytes.Buffer
		var mockCmd *cobra.Command
		var out io.Writer
		BeforeEach(func() {
			buf = gbytes.NewBuffer()
			mockCmd = &cobra.Command{}
			mockCmd.SetOutput(buf)
			out = mockCmd.OutOrStdout()
		})
		It("will return error when given no tasks", func() {
			server.SetAllowUnhandledRequests(true)

			err := completeCmd.RunE(mockCmd, []string{"1"})
			Expect(err).To(HaveOccurred())
			ghe := err.(*api.GoHabitsError)
			Expect(ghe.StatusCode).Should(Equal(1))
			Eventually(out).Should(gbytes.Say("You have no tasks"))
		})
	})
})
