package cmd

import (
	"math/rand"
	"strings"

	api "github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
				task := parseTask(args)

				Expect(task.Title).To(MatchRegexp(strings.Join(args[0:2], " ")))
				Expect(task.Tags).To(ContainElement("health"))
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
