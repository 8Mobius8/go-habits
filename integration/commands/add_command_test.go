package commands

import (
	"fmt"
	"io/ioutil"
	"math/rand"
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
			var tasksFileName string
			BeforeEach(func() {
				f, err := ioutil.TempFile("./", "tasks")
				f.Close()
				Expect(err).ToNot(HaveOccurred())
				tasksFileName = f.Name()
			})
			AfterEach(func() {
				RemoveFile(tasksFileName)
			})

			It("will accept --file or -f to specify file", func() {
				s := GoHabits("add", "-f", tasksFileName)
				Eventually(s).Should(gexec.Exit(0))

				s = GoHabits("add", "--file", tasksFileName)
				Eventually(s).Should(gexec.Exit(0))
			})

			It("will print new created task as confirmation", func() {
				AppendToFile(tasksFileName, TaskLine("A new task"))
				AppendToFile(tasksFileName, TaskLine("A new task #dup"))

				s := GoHabits("add", "-f", tasksFileName)
				Eventually(s).Should(gbytes.Say("A new task"))
				Eventually(s).Should(gbytes.Say("A new task #dup"))
				Eventually(s).Should(gexec.Exit(0))
			})

			It("will ignore lines that do not have prefixed task list", func() {
				AppendToFile(tasksFileName, TaskLine("A new task"))
				for i := 1; i < rand.Intn(5); i++ {
					AppendToFile(tasksFileName, "a line that isn't a task\n")
				}

				s := GoHabits("add", "-f", tasksFileName)
				Eventually(s).Should(gbytes.Say("A new task"))
				Eventually(s).Should(gexec.Exit(0))
			})

			It("will print 'missing' lines that have an empty task title", func() {
				AppendToFile(tasksFileName, TaskLine(""))

				s := GoHabits("add", "-f", tasksFileName)
				Eventually(s).Should(gbytes.Say("missing text"))
				Eventually(s).Should(gexec.Exit(0))
			})
		})
	})
})

func AppendToFile(filepath, s string) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	Expect(err).ToNot(HaveOccurred())

	defer f.Close()

	_, err = f.WriteString(s)
	Expect(err).ToNot(HaveOccurred())
}

func TaskLine(title string) string {
	return fmt.Sprintf("[ ] %s\n", title)
}

func RemoveFile(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
	}
	Expect(err).ShouldNot(HaveOccurred())
}
