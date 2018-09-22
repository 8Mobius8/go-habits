package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits list", func() {

	It("actually exists", func() {
		session := GoHabits("list")
		Eventually(session).Should(gexec.Exit(0))
	})

	Context("User is already signed in", func() {
		var configPath string
		BeforeEach(func() {
			configPath = TouchConfigFile()
			expectSuccessfulLogin(userName, password)
		})
		AfterEach(func() {
			RemoveConfigFile(configPath)
		})
		It("lists a task", func() {
			addTask("Do the dishes")

			session := GoHabits("list")
			Eventually(session).Should(gbytes.Say("Do the dishes"))
			Eventually(session).Should(gexec.Exit(0))
		})
	})
})

type Task struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type TaskResponse struct {
	Id   string
	Text string
}

func addTask(taskTitle string) {
	todoTask := Task{taskTitle, "todo"}
	res := TaskResponse{}
	err := apiClient.Post("/tasks/user", todoTask, res)
	Expect(err).ShouldNot(HaveOccurred())
}

func expectSuccessfulLogin(user, password string) {
	s, in := GoHabitsWithStdin("login")
	EventuallyLogin(s, in, user, password)
	Eventually(s).Should(gexec.Exit(0))
}
