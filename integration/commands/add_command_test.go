package commands

import (
	"io/ioutil"
	"os"

	. "github.com/8Mobius8/go-habits/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits add command", func() {
	It("actually exists", func() {
		session := GoHabits("add", "-h")
		Eventually(session).Should(gexec.Exit(0))
	})

	Context("Given the user is already signed-in", func() {
		BeforeEach(func() {
			TouchConfigFile()
			expectSuccessfulLogin(UserName, Password)
		})
		AfterEach(func() {
			ResetUser()
			RemoveConfigFile()
		})

		Describe("Creating a new todo interactively", func() {
			It("will print newly created task as confirmation", func() {
				s := GoHabits("add", "clean my fishbowl")
				Eventually(s).Should(gbytes.Say("clean my fishbowl"))
				Eventually(s).Should(gexec.Exit(0))
			})

			It("will print newly created task and order as confirmation", func() {
				s := GoHabits("add", "clean my fishbowl")
				Eventually(s).Should(gexec.Exit(0))
				s = GoHabits("add", "make bed")
				Eventually(s).Should(gbytes.Say("1"))
				Eventually(s).Should(gbytes.Say("make bed"))
				Eventually(s).Should(gexec.Exit(0))
			})

			It("will print new tasks with tags in order as confirmation", func() {
				s := GoHabits("add", "clean my fishbowl #chores")
				Eventually(s).Should(gbytes.Say("1"))
				Eventually(s).Should(gbytes.Say("clean my fishbowl #chores"))
				Eventually(s).Should(gexec.Exit(0))
			})
		})

		Describe("Creating a new todo using file", func() {
			tasksFile := "tasks.txt"
			BeforeEach(func() {
				AddTaskToFile(tasksFile, "A new task")
			})
			AfterEach(func() {
				RemoveTaskFile(tasksFile)
			})

			It("will print new created task as confirmation", func() {
				s := GoHabits("add", "-f", "tasks.txt", "--log", "debug")
				Eventually(s).Should(gbytes.Say("A new task"))
				Eventually(s).Should(gexec.Exit(0))
			})
		})
	})
})

func AddTaskToFile(filePath, taskTitle string) {
	err := ioutil.WriteFile(filePath, []byte("[ ] "+taskTitle+"\n"), 0644)
	Expect(err).ShouldNot(HaveOccurred())
}

func RemoveTaskFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
	}
	Expect(err).ShouldNot(HaveOccurred())
}
