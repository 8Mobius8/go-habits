package commands_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	api "github.com/8Mobius8/go-habits/api"
	cmd "github.com/8Mobius8/go-habits/cmd/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Add command", func() {
	Describe("FilterTask", func() {
		Context("given an empty id with some tasks", func() {
			It("should return empty array", func() {
				id := ""
				tasks := generateTasks(5)
				filtered := cmd.FilterTask(id, tasks)
				Expect(filtered).Should(BeEmpty())
			})
		})
		Context("given an id that matches one task", func() {
			It("should return an array with that task", func() {
				id := randomID()
				tasks := generateTasks(3)

				tasks[2].ID = id

				filtered := cmd.FilterTask(id, tasks)
				Expect(filtered).Should(HaveLen(1))
				Expect(filtered[0].ID).Should(Equal(id))
			})
		})
	})

	Describe("ParseTask", func() {
		DescribeTable("should return a task with it's title set",
			func(args []string, expectedTitle string) {
				task := cmd.ParseTask(args)
				Expect(task).To(Equal(api.Task{Title: expectedTitle, Type: api.TodoType.String()}))
			},
			Entry("when given a single word", []string{"eat"}, "eat"),
			Entry("when given a multiple words", []string{"eat", "treats"}, "eat treats"),
			Entry("when given a multiple words as single argument", []string{"eat treats"}, "eat treats"),
		)

		DescribeTable("should return a task with it's title and tags set",
			func(args []string, expectedTitle string, expectedTags []string) {
				task := cmd.ParseTask(args)
				Expect(task).To(Equal(api.Task{Title: expectedTitle, Tags: expectedTags, Type: api.TodoType.String()}))
			},
			Entry("when given a single word and single tag", []string{"eat", "#chore"}, "eat", []string{"chore"}),
			Entry("when given a single word and single tag in reverse", []string{"#chore", "eat"}, "eat", []string{"chore"}),
			Entry("when given a single word and multiple tags", []string{"eat", "#chore", "#delight"}, "eat", []string{"chore", "delight"}),
		)
	})

	Describe("ParseTasksFromFile contents of file and parse tasks", func() {
		var taskfileName string
		var fileContents string
		JustBeforeEach(func() {
			taskfileName = randomTaskFileName()
			err := ioutil.WriteFile(taskfileName, []byte(fileContents), 0644)
			Expect(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			err := os.Remove(taskfileName)
			if err != nil {
				if os.IsNotExist(err) {
					return
				}
			}
			Expect(err).ShouldNot(HaveOccurred())
			fileContents = ""
		})
		Context("when a single task is inside the file", func() {
			BeforeEach(func() {
				fileContents = "[ ] soak dreams in coke-a-cola"
			})
			It("no error and returns correct task name", func() {
				t, err := cmd.ParseTasksFromFile(taskfileName)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(t)).To(Equal(1))
				Expect(t[0].Title).To(Equal("soak dreams in coke-a-cola"))
			})
		})
		Context("when a random task is inside the file", func() {
			var expectedName string
			BeforeEach(func() {
				expectedName = randomTaskName()
				fileContents = "[ ] " + expectedName
			})
			It("no error and returns correct task name", func() {
				t, err := cmd.ParseTasksFromFile(taskfileName)
				Expect(len(t)).To(Equal(1))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(t[0].Title).To(Equal(expectedName))
			})
		})
		Context("when 5 random task is inside the file", func() {
			var expectedNames = []string{}
			BeforeEach(func() {
				for i := 0; i < 5; i++ {
					taskName := randomTaskName()
					expectedNames = append(expectedNames, taskName)
					fileContents += "[ ] " + taskName + "\n"
				}
			})
			It("no error and returns correct task names", func() {
				t, err := cmd.ParseTasksFromFile(taskfileName)
				Expect(len(t)).To(Equal(len(expectedNames)))
				Expect(err).ShouldNot(HaveOccurred())
				for i := 0; i < 5; i++ {
					Expect(t[i].Title).To(Equal(expectedNames[i]))
				}
			})
		})
		Context("when a bad line is inside the file", func() {
			BeforeEach(func() {
				fileContents = "bad line" + randomTaskName()
			})
			It("Will return with an empty list", func() {
				t, err := cmd.ParseTasksFromFile(taskfileName)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(t)).To(Equal(0))
			})
		})
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
				err := cmd.Add(out, server, args, "")
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

func randomTaskFileName() string {
	return "taskfile" + randomString(5) + ".txt"
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
