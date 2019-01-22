package cmd_test

import (
	"github.com/8Mobius8/go-habits/api"
	"github.com/8Mobius8/go-habits/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Remove cmd", func() {
	Context("when api returns successful task return and deletion", func() {
		var out *gbytes.Buffer
		BeforeEach(func() {
			out = gbytes.NewBuffer()
		})
		AfterEach(func() {
			out.Close()
		})
		It("prints 'Removed Tasks' and lists them", func() {
			aTask := api.Task{
				Title: "A Chore of mine",
				Order: 1,
			}
			server := MockDeleteServer{
				DeleteTaskFunc: func(t api.Task) error {
					return nil
				},
				GetTasksFunc: func(api.TaskType) []api.Task {
					return []api.Task{aTask}
				},
			}

			err := cmd.Remove(out, []string{"1"}, server)
			Expect(err).ToNot(HaveOccurred())
			Eventually(out).Should(gbytes.Say("Removed tasks:"))
			Eventually(out).Should(gbytes.Say(aTask.Title))
		})
	})
})

type MockDeleteServer struct {
	DeleteTaskFunc  func(api.Task) error
	GetTasksFunc    func(api.TaskType) []api.Task
	AddTaskFunc     func(api.Task) (api.Task, error)
	ScoreTaskUpFunc func(api.Task) error
}

func (s MockDeleteServer) DeleteTask(t api.Task) error {
	return s.DeleteTaskFunc(t)
}
func (s MockDeleteServer) GetTasks(t api.TaskType) []api.Task {
	return s.GetTasksFunc(t)
}
func (s MockDeleteServer) AddTask(t api.Task) (api.Task, error) {
	return s.AddTaskFunc(t)
}
func (s MockDeleteServer) ScoreTaskUp(t api.Task) error {
	return s.ScoreTaskUpFunc(t)
}
