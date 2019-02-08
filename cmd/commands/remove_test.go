package commands_test

import (
	"github.com/8Mobius8/go-habits/api"
	cmd "github.com/8Mobius8/go-habits/cmd/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Remove cmd", func() {
	Context("when api returns successful task return and deletion", func() {
		var out, in *gbytes.Buffer
		BeforeEach(func() {
			out = gbytes.NewBuffer()
			in = gbytes.NewBuffer()
		})
		AfterEach(func() {
			out.Close()
			in.Close()
		})
		Context("given removing with force", func() {
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

				err := cmd.Remove(in, out, []string{"1"}, server, true)
				Expect(err).ToNot(HaveOccurred())
				Eventually(out).Should(gbytes.Say("Removed tasks"))
				Eventually(out).Should(gbytes.Say(aTask.Title))
			})
		})
		Context("Given removing with/out force", func() {
			FIt("Asks user if they want to remove list of tasks and removes them given yes", func() {
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

				err := cmd.Remove(in, out, []string{"1"}, server, false)
				Expect(err).ToNot(HaveOccurred())
				Eventually(out).Should(gbytes.Say("Remove?"))
				Eventually(out).Should(gbytes.Say(aTask.Title))
				Eventually(out).Should(gbytes.Say("[Y\\n]?"))

				in.Write([]byte("Y\n"))

				Eventually(out).Should(gbytes.Say("Removed tasks"))
				Eventually(out).Should(gbytes.Say(aTask.Title))
			})
			It("Asks user if they want to remove list of tasks and does not removes them given no", func() {
				aTask := api.Task{
					Title: "A Chore of mine",
					Order: 1,
				}
				server := MockDeleteServer{
					DeleteTaskFunc: func(t api.Task) error {
						Fail("DeleteTask on api should not be called")
						return nil
					},
					GetTasksFunc: func(api.TaskType) []api.Task {
						return []api.Task{aTask}
					},
				}

				err := cmd.Remove(in, out, []string{"1"}, server, false)
				Expect(err).ToNot(HaveOccurred())
				Eventually(out).Should(gbytes.Say("Remove?"))
				Eventually(out).Should(gbytes.Say(aTask.Title))
				Eventually(out).Should(gbytes.Say("[Y\\n]"))

				in.Write([]byte("N\n"))
			})
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
