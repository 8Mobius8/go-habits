package cmd

import (
	api "github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List command", func() {
	Describe("formatTask", func() {
		Context("given a task without tags", func() {
			It("will return a formated string with the values filled in", func() {
				t := api.Task{}
				t.ID = randomID()
				t.Order = 1
				t.Title = " a simple task I need to complete"

				Expect(formatTask(t)).Should(MatchRegexp(`[0-9]+\[ \] [\s\w]+`))
			})
		})
		Context("given a task with tags", func() {
			It("will return a formated string with the values filled in", func() {
				t := api.Task{}
				t.ID = randomID()
				t.Order = 1
				t.Title = " a simple task I need to complete"
				t.Tags = []string{"misc", "uncategorised"}

				Expect(formatTask(t)).Should(MatchRegexp(`[0-9]+\[ \] [\s\w]+ (#[\w]+ )+`))
			})
		})

	})
})
