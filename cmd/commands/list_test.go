package commands_test

import (
	api "github.com/8Mobius8/go-habits/api"
	cmd "github.com/8Mobius8/go-habits/cmd/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type MockListServer struct {
	GetTasksFunc func(api.TaskType) []api.Task
}

func (s MockListServer) GetTasks(tt api.TaskType) []api.Task {
	return s.GetTasksFunc(tt)
}

var _ = Describe("List command", func() {
	Describe("List", func() {
		var out *gbytes.Buffer
		BeforeEach(func() {
			out = gbytes.NewBuffer()
		})
		Context("given a task from the server", func() {
			It("should print formated task", func() {
				mockServer := MockListServer{
					func(tt api.TaskType) []api.Task {
						task := api.Task{
							Title: "Some Task",
							Order: 1,
						}
						return []api.Task{task}
					},
				}

				cmd.List(out, mockServer, []string{})
				Eventually(out).Should(gbytes.Say("1"))
				Eventually(out).Should(gbytes.Say("[ ]"))
				Eventually(out).Should(gbytes.Say("Some Task"))
			})
		})
		Context("given several tasks", func() {
			It("should print several tasks", func() {
				mockServer := MockListServer{
					func(tt api.TaskType) []api.Task {
						task := api.Task{
							Title: "Some Task",
							Order: 1,
						}
						task2 := task
						task2.Order = 2
						task3 := task2
						task3.Order = 3
						return []api.Task{task, task2, task3}
					},
				}

				cmd.List(out, mockServer, []string{})
				Eventually(out).Should(gbytes.Say("1"))
				Eventually(out).Should(gbytes.Say("[ ]"))
				Eventually(out).Should(gbytes.Say("Some Task"))
				Eventually(out).Should(gbytes.Say("2"))
				Eventually(out).Should(gbytes.Say("[ ]"))
				Eventually(out).Should(gbytes.Say("Some Task"))
				Eventually(out).Should(gbytes.Say("3"))
				Eventually(out).Should(gbytes.Say("[ ]"))
				Eventually(out).Should(gbytes.Say("Some Task"))
			})
		})
		Context("given a task with a tag from the server", func() {
			It("should print formated task", func() {
				mockServer := MockListServer{
					func(tt api.TaskType) []api.Task {
						task := api.Task{
							Title: "Some Task",
							Order: 1,
							Tags:  []string{"tag"},
						}
						return []api.Task{task}
					},
				}

				cmd.List(out, mockServer, []string{})
				Eventually(out).Should(gbytes.Say("1"))
				Eventually(out).Should(gbytes.Say("[ ]"))
				Eventually(out).Should(gbytes.Say("Some Task"))
				Eventually(out).Should(gbytes.Say("#tag"))
			})
		})
		Context("given several tasks with a tag from the server", func() {
			It("should print the formated tasks", func() {
				mockServer := MockListServer{
					func(tt api.TaskType) []api.Task {
						task := api.Task{
							Title: "Some Task",
							Order: 1,
							Tags:  []string{"tag"},
						}
						task2 := task
						task2.Order = 2
						task3 := task2
						task3.Order = 3
						return []api.Task{task, task2, task3}
					},
				}

				cmd.List(out, mockServer, []string{})
				Eventually(out).Should(gbytes.Say("1"))
				Eventually(out).Should(gbytes.Say("[ ]"))
				Eventually(out).Should(gbytes.Say("Some Task"))
				Eventually(out).Should(gbytes.Say("#tag"))
				Eventually(out).Should(gbytes.Say("2"))
				Eventually(out).Should(gbytes.Say("[ ]"))
				Eventually(out).Should(gbytes.Say("Some Task"))
				Eventually(out).Should(gbytes.Say("#tag"))
				Eventually(out).Should(gbytes.Say("3"))
				Eventually(out).Should(gbytes.Say("[ ]"))
				Eventually(out).Should(gbytes.Say("Some Task"))
				Eventually(out).Should(gbytes.Say("#tag"))
			})
		})
	})
})
