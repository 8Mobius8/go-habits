package cmd

import (
	"math/rand"
	"strings"

	api "github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Add command", func() {
	Describe("filterTask", func() {
		Context("given an empty id with some tasks", func() {
			It("should return empty array", func() {
				id := ""
				tasks := generateTasks(5)
				filtered := filterTask(id, tasks)
				Expect(filtered).Should(BeEmpty())
			})
		})
		Context("given an id that matches one task", func() {
			It("should return an array with that task", func() {
				id := randomID()
				tasks := generateTasks(3)

				tasks[2].ID = id

				filtered := filterTask(id, tasks)
				Expect(filtered).Should(HaveLen(1))
				Expect(filtered[0].ID).Should(Equal(id))
			})
		})
	})

	Describe("parseTaskTitle", func() {
		Context("given a single word as arguments", func() {
			It("should return the title as the word", func() {
				args := []string{"eat"}
				title := parseTaskTitle(args)

				Expect(title).To(Equal("eat"))
			})
		})

		Context("given a multiple words as arguments", func() {
			It("should return the title as the words separated by spaces", func() {
				args := []string{"eat", "breakfast"}
				title := parseTaskTitle(args)

				Expect(title).To(MatchRegexp(strings.Join(args, " ")))
			})
		})

		Context("given words and tags as arguments", func() {
			It("should return with title and tags set", func() {
				args := []string{"eat", "breakfast", "#health"}
				task := ParseTask(args)

				Expect(task.Title).To(MatchRegexp(strings.Join(args[0:2], " ")))
				Expect(task.Tags).To(ContainElement("health"))
			})
		})
	})

	Describe("ParseTask", func() {
		DescribeTable("should return a task with it's title set",
			func(args []string, expectedTitle string) {
				task := ParseTask(args)
				Expect(task).To(Equal(api.Task{Title: expectedTitle, Type: api.TodoType.String()}))
			},
			Entry("when given a single word", []string{"eat"}, "eat"),
			Entry("when given a multiple words", []string{"eat", "treats"}, "eat treats"),
			Entry("when given a multiple words as single argument", []string{"eat treats"}, "eat treats"),
		)

		DescribeTable("should return a task with it's title and tags set",
			func(args []string, expectedTitle string, expectedTags []string) {
				task := ParseTask(args)
				Expect(task).To(Equal(api.Task{Title: expectedTitle, Tags: expectedTags, Type: api.TodoType.String()}))
			},
			Entry("when given a single word and single tag", []string{"eat", "#chore"}, "eat", []string{"chore"}),
			Entry("when given a single word and single tag in reverse", []string{"#chore", "eat"}, "eat", []string{"chore"}),
			Entry("when given a single word and multiple tags", []string{"eat", "#chore", "#delight"}, "eat", []string{"chore", "delight"}),
		)
	})

	Describe("Add", func() {
		var out *gbytes.Buffer
		BeforeEach(func() {
			out = gbytes.NewBuffer()
		})
		AfterEach(func() {
			out.Close()
		})
		Context("given the task has regular character phrase", func() {
			It("will have no error return and prints tasks", func() {
				aTask := api.Task{Title: "A chore that must be completed", Order: 1}
				args := strings.Split(aTask.Title, " ")
				server := MockAddTaskServer{
					AddTaskFunc: func(t api.Task) (api.Task, error) {
						return aTask, nil
					},
					GetTasksFunc: func(t api.TaskType) []api.Task {
						return []api.Task{aTask}
					},
				}
				err := Add(out, server, args)
				Expect(err).ToNot(HaveOccurred())
				Eventually(out).Should(gbytes.Say("1"))
				Eventually(out).Should(gbytes.Say("[ ]"))
				Eventually(out).Should(gbytes.Say(aTask.Title))
			})
		})
	})
})

func generateTasks(num int) []api.Task {
	var tasks []api.Task
	for i := 0; i < num; i++ {
		t := api.Task{}
		t.ID = randomID()
		t.Title = randomTaskName()
		tasks = append(tasks, t)
	}
	return tasks
}

func randomID() string {
	id := randomString(8)
	id += "-"
	id += randomString(4)
	id += "-"
	id += randomString(4)
	id += "-"
	id += randomString(8)
	return id
}

const pool = "0987654321abcdefghijklmnopqrstuvwxyz"

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}
	return string(bytes)
}

func randomTaskName() string {
	return "task-" + randomString(10)
}

type MockAddTaskServer struct {
	AddTaskFunc  func(api.Task) (api.Task, error)
	GetTasksFunc func(api.TaskType) []api.Task
}

func (m MockAddTaskServer) AddTask(t api.Task) (api.Task, error) {
	return m.AddTaskFunc(t)
}
func (m MockAddTaskServer) GetTasks(t api.TaskType) []api.Task {
	return m.GetTasksFunc(t)
}
