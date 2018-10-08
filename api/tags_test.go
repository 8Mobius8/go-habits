package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Tags", func() {
	BeforeEach(func() {
		tagsCache = make(map[string]string)
	})

	Describe("getTag", func() {
		It("returns tasks with same tags, as words", func() {
			var tagID = "d268201e-c926-4a32-8ac1-7ca570c26b45"
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tags/d268201e-c926-4a32-8ac1-7ca570c26b45"),
					ghttp.RespondWith(200, ChoresTag),
				),
			)

			tag := habitapi.getTag(tagID)

			Expect(tag.ID).Should(Equal(tagID))
			Expect(tag.Name).Should(Equal("chores"))
		})
	})

	Describe("getTags", func() {
		It("returns tasks with same tags, as words", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tasks/user", "type=todos"),
					ghttp.RespondWith(200, ChoresTasks),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tags/d268201e-c926-4a32-8ac1-7ca570c26b45"),
					ghttp.RespondWith(200, ChoresTag),
				),
			)

			tasks := habitapi.GetTasks(TodoType)

			Expect(len(tasks)).Should(BeNumerically("==", 3))
			for _, task := range tasks {
				Expect(task.Tags).Should(ConsistOf("chores"))
			}
		})

		It("returns tasks with different tags, as words", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tasks/user", "type=todos"),
					ghttp.RespondWith(200, ChoresWorkTasks),
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

			tasks := habitapi.GetTasks(TodoType)

			Expect(len(tasks)).Should(BeNumerically("==", 6))
			Expect(tasks[0].Tags).Should(ConsistOf("chores"))
			Expect(tasks[1].Tags).Should(ConsistOf("chores"))
			Expect(tasks[2].Tags).Should(ConsistOf("chores"))
			Expect(tasks[3].Tags).Should(ConsistOf("work"))
			Expect(tasks[4].Tags).Should(ConsistOf("work"))
			Expect(tasks[5].Tags).Should(ConsistOf("work"))
		})
	})
})
