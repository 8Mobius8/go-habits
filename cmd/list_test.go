package cmd

import (
	api "github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List command", func() {
	Describe("formatTodo", func() {
		Context("given a todo without tags", func() {
			It("will return a formated string with the values filled in", func() {
				t := api.Todo{}
				t.ID = randomId()
				t.Order = 1
				t.Title = " a simple todo I need to complete"

				Expect(formatTodo(t)).Should(MatchRegexp(`[0-9]+\[ \] [\s\w]+`))
			})
		})
		Context("given a todo with tags", func() {
			It("will return a formated string with the values filled in", func() {
				t := api.Todo{}
				t.ID = randomId()
				t.Order = 1
				t.Title = " a simple todo I need to complete"
				t.Tags = []string{"misc", "uncategorised"}

				Expect(formatTodo(t)).Should(MatchRegexp(`[0-9]+\[ \] [\s\w]+ (#[\w]+ )+`))
			})
		})

	})
})
