package commands

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/8Mobius8/go-habits/api"
	. "github.com/8Mobius8/go-habits/integration"
)

var _ = Describe("go-habits remove command", func() {
	It("actually exists", func() {
		session := GoHabits("remove", "-h")
		Eventually(session).Should(gexec.Exit(0))
	})

	It("exits safely when showing usage", func() {
		session := GoHabits("remove", "--help")
		Eventually(session).Should(gbytes.Say("Usage:"))
		Eventually(session).Should(gbytes.Say(`go-habits remove \[flags\]`))
		Eventually(session).Should(gexec.Exit(0))
	})

	Context("when tasks has been created", func() {
		var task api.Task
		BeforeEach(func() {
			TouchConfigFile()
			expectSuccessfulLogin(UserName, Password)

			task = api.NewTask("A Task to complete", api.TodoType)
			t, err := APIClient.AddTask(task)
			task = t
			Expect(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			ResetUser()
			RemoveConfigFile()
		})
		Describe("and remove command is given first task by number", func() {
			It("should exit succesfully, and print removed tasks", func() {
				s := GoHabits("remove", "-f", "1")
				Eventually(s).Should(gbytes.Say("Removed tasks"))
				Eventually(s).Should(gbytes.Say("X"))
				Eventually(s).Should(gbytes.Say(task.Title))
				Eventually(s).Should(gexec.Exit(0))
			})
			It("should ask for confirmation and print removed tasks", func() {
				s, in := GoHabitsWithStdin("remove", "1")

				Eventually(s).Should(gbytes.Say("Remove?"))
				Eventually(s).Should(gbytes.Say("1"))
				Eventually(s).Should(gbytes.Say(task.Title))
				Eventually(s).Should(gbytes.Say("[Y\\n]?"))
				in.Write([]byte("Y\n"))

				Eventually(s).Should(gbytes.Say("Removed tasks"))
				Eventually(s).Should(gbytes.Say("X"))
				Eventually(s).Should(gbytes.Say(task.Title))
				Eventually(s).Should(gexec.Exit(0))
			})
		})
	})
})
