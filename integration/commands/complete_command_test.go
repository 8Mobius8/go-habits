package commands

import (
	"github.com/8Mobius8/go-habits/api"
	. "github.com/8Mobius8/go-habits/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits complete command", func() {
	BeforeEach(func() {
		TouchConfigFile()
		expectSuccessfulLogin(UserName, Password)
	})
	AfterEach(func() {
		ResetUser(ApiID, ApiToken)
		RemoveConfigFile()
	})
	It("exits safely when showing usage", func() {
		session := GoHabits("complete", "--help")
		Eventually(session).Should(gbytes.Say("Usage:"))
		Eventually(session).Should(gbytes.Say(`go-habits complete \[flags\]`))
		Eventually(session).Should(gexec.Exit(0))
	})
	Describe("when completing a task by number", func() {
		Context("given a task has been already created", func() {
			var task api.Task
			BeforeEach(func() {
				task = api.NewTask("A Task to complete", api.TodoType)
				t, err := ApiClient.AddTask(task)
				task = t
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("will print task completed", func() {
				s := GoHabits("complete", "1")
				Eventually(s).Should(gbytes.Say("[X]"))
				Eventually(s).Should(gbytes.Say(task.Title))
				Eventually(s).Should(gexec.Exit(0))
			})
			It("will mark the task completed on the server", func() {
				s := GoHabits("complete", "1")
				Eventually(s).Should(gexec.Exit(0))

				t := api.Task{}
				err := ApiClient.Get("/tasks/"+task.ID, &t)
				Expect(err).ToNot(HaveOccurred())

				Expect(t.Completed).Should(Equal(true))
			})
		})
		Context("given no task has been created", func() {
			It("will print no tasks have been created.", func() {
				s := GoHabits("complete", "1")
				Eventually(s).Should(gbytes.Say("You have no tasks."))
				Eventually(s).Should(gbytes.Say("Create tasks before trying to complete them."))
				Eventually(s).Should(gexec.Exit(1))
			})
		})
	})
})
