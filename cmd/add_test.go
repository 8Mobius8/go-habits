package cmd

import (
	"math/rand"
	"strings"

	api "github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Add command", func() {
	Describe("filterTodo", func() {
		Context("given an empty id with some todos", func() {
			It("should return empty array", func() {
				id := ""
				todos := generateTodos(5)
				filtered := filterTodo(id, todos)
				Expect(filtered).Should(BeEmpty())
			})
		})
		Context("given an id that matches one todo", func() {
			It("should return an array with that todo", func() {
				id := randomId()
				todos := generateTodos(3)

				todos[2].ID = id

				filtered := filterTodo(id, todos)
				Expect(filtered).Should(HaveLen(1))
				Expect(filtered[0].ID).Should(Equal(id))
			})
		})
	})

	Describe("parseTodoTitle", func() {
		Context("given a single word as arguments", func() {
			It("should return the title as the word", func() {
				args := []string{"eat"}
				title := parseTodoTitle(args)

				Expect(title).To(Equal("eat"))
			})
		})

		Context("given a multiple words as arguments", func() {
			It("should return the title as the words separated by spaces", func() {
				args := []string{"eat", "breakfast"}
				title := parseTodoTitle(args)

				Expect(title).To(MatchRegexp(strings.Join(args, " ")))
			})
		})
	})
})

func generateTodos(num int) []api.Todo {
	var todos []api.Todo
	for i := 0; i < num; i++ {
		t := api.Todo{}
		t.ID = randomId()
		t.Title = randomTaskName()
		todos = append(todos, t)
	}
	return todos
}

func randomId() string {
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
