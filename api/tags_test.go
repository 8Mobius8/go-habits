package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Todos", func() {
	BeforeEach(func() {
		tagsCache = make(map[string]string)
	})

	Describe("getTag", func() {
		It("returns todos with same tags, as words", func() {
			var tagID = "d268201e-c926-4a32-8ac1-7ca570c26b45"
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tags/d268201e-c926-4a32-8ac1-7ca570c26b45"),
					ghttp.RespondWith(200, ChoresTag),
				),
			)

			tag := habitapi.getTag(tagID)

			Expect(tag.Id).Should(Equal(tagID))
			Expect(tag.Name).Should(Equal("chores"))
		})
	})

	Describe("getTags", func() {
		It("returns todos with same tags, as words", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tasks/user"),
					ghttp.RespondWith(200, ChoresTodos),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tags/d268201e-c926-4a32-8ac1-7ca570c26b45"),
					ghttp.RespondWith(200, ChoresTag),
				),
			)

			todos := habitapi.GetTodos()

			Expect(len(todos)).Should(BeNumerically("==", 3))
			for _, todo := range todos {
				Expect(todo.Tags).Should(ConsistOf("chores"))
			}
		})

		It("returns todos with different tags, as words", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tasks/user"),
					ghttp.RespondWith(200, ChoresWorkTodos),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tags/d268201e-c926-4a32-8ac1-7ca570c26b45"),
					ghttp.RespondWith(200, ChoresTag),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tags/d268201e-26c9-32a4-a81c-570c7ca4526b"),
					ghttp.RespondWith(200, WorkTag),
				),
			)

			todos := habitapi.GetTodos()

			Expect(len(todos)).Should(BeNumerically("==", 6))
			Expect(todos[0].Tags).Should(ConsistOf("chores"))
			Expect(todos[1].Tags).Should(ConsistOf("chores"))
			Expect(todos[2].Tags).Should(ConsistOf("chores"))
			Expect(todos[3].Tags).Should(ConsistOf("work"))
			Expect(todos[4].Tags).Should(ConsistOf("work"))
			Expect(todos[5].Tags).Should(ConsistOf("work"))
		})
	})
})
