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

		BeforeEach(func() {
			TouchConfigFile()
			expectSuccessfulLogin(userName, password)
		})
		AfterEach(func() {
			RemoveConfigFile()
		})
		It("lists a task", func() {
			addTask("Do the dishes")

			session := GoHabits("list")
			Eventually(session).Should(gbytes.Say("Do the dishes"))
			Eventually(session).Should(gexec.Exit(0))
		})
		It("lists a task with it's tag", func() {
			task := addTask("Clean the bed")
			tag := addTag("chores")
			addTagToTask(task.Id, tag.Id)

			session := GoHabits("list")
			Eventually(session).Should(gbytes.Say("Clean the bed"))
			Eventually(session).Should(gbytes.Say("#chore"))
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

func addTask(taskTitle string) TaskResponse {
	todoTask := Task{taskTitle, "todo"}
	var res TaskResponse
	err := apiClient.Post("/tasks/user", todoTask, &res)
	Expect(err).ShouldNot(HaveOccurred())
	return res
}

type Tag struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func addTag(tagTitle string) Tag {
	tagInput := struct {
		Name string `json:"name"`
	}{Name: tagTitle}
	var tagResponse Tag
	err := apiClient.Post("/tags", tagInput, &tagResponse)
	Expect(err).ShouldNot(HaveOccurred())
	return tagResponse
}

func addTagToTask(taskId, tagId string) {
	var e struct{}
	err := apiClient.Post("/tasks/"+taskId+"/tags/"+tagId, e, e)
	Expect(err).ShouldNot(HaveOccurred())
}

func expectSuccessfulLogin(user, password string) {
	s, in := GoHabitsWithStdin("login")
	EventuallyLogin(s, in, user, password)
	Eventually(s).Should(gexec.Exit(0))
}
