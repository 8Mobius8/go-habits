package cmd_test

import (
	"github.com/8Mobius8/go-habits/api"
	"github.com/8Mobius8/go-habits/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type MockCompleteTaskServer struct {
	GetTasksFunc    func(api.TaskType) []api.Task
	AddTaskFunc     func(api.Task) (api.Task, error)
	ScoreTaskUpFunc func(api.Task) error
}

func (ts MockCompleteTaskServer) GetTasks(tt api.TaskType) []api.Task {
	return ts.GetTasksFunc(tt)
}
func (ts MockCompleteTaskServer) AddTask(t api.Task) (api.Task, error) {
	return ts.AddTaskFunc(t)
}
func (ts MockCompleteTaskServer) ScoreTaskUp(t api.Task) error {
	return ts.ScoreTaskUpFunc(t)
}

var _ = Describe("Complete integration", func() {
	Describe("Complete command", func() {
		var out *gbytes.Buffer
		BeforeEach(func() {
			out = gbytes.NewBuffer()
		})
		It("will return error when given no tasks", func() {
			mockServer := MockCompleteTaskServer{
				GetTasksFunc: func(tt api.TaskType) []api.Task {
					return []api.Task{}
				},
			}

			err := cmd.Complete(out, mockServer, []string{"1"})
			Expect(err).To(HaveOccurred())
			ghe := err.(*api.GoHabitsError)
			Expect(ghe.StatusCode).Should(Equal(1))
			Eventually(out).Should(gbytes.Say("You have no tasks"))
		})
	})

	Describe("GetTaskAtPosition", func() {
		It("will return an error when no tasks exist", func() {
			taskServer := &MockCompleteTaskServer{
				GetTasksFunc: func(api.TaskType) []api.Task {
					return []api.Task{}
				},
			}
			_, err := cmd.GetTaskAtPosition(taskServer, 1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no tasks"))
		})
		It("will return an error when position doesn't exist", func() {
			task := api.Task{Title: "Task 1"}
			ts := &MockCompleteTaskServer{
				GetTasksFunc: func(api.TaskType) []api.Task {
					return []api.Task{task}
				},
			}

			_, err := cmd.GetTaskAtPosition(ts, 1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("bad index"))
		})
		It("will return an error when position is negitive", func() {
			task := api.Task{Title: "Task 1"}
			ts := &MockCompleteTaskServer{
				GetTasksFunc: func(api.TaskType) []api.Task {
					return []api.Task{task}
				},
			}
			_, err := cmd.GetTaskAtPosition(ts, -1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("bad index"))
		})
		It("will return the task at the right position.", func() {
			tasks := []api.Task{
				{Title: "Task 1", Order: 1},
				{Title: "Task 2", Order: 2},
				{Title: "Task 3", Order: 3},
			}
			ts := &MockCompleteTaskServer{
				GetTasksFunc: func(api.TaskType) []api.Task {
					return tasks
				},
			}
			t, err := cmd.GetTaskAtPosition(ts, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(t.Title).Should(BeEquivalentTo("Task 2"))
			Expect(t.Order).Should(BeEquivalentTo(2))
		})
	})
})
