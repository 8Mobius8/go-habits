package api

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func ValidTagJSON(tagName string) string {
	return `{
		"success": true,
		"data": {
			"name": "` + tagName + `",
			"id": "b2131ed2-2e5d-4ad1-affb-db59342ad227"
		}
	}`
}

func ValidTagWithIdJSON(tagName string, id string) string {
	return `{
		"success": true,
		"data": {
			"name": "` + tagName + `",
			"id": "` + id + `"
		}
	}`
}

func CollectionTagsJSON(tagNames []string) string {

	dataField := fmt.Sprintln("[")
	for i, tagName := range tagNames {
		dataField += fmt.Sprintln("{")
		dataField += fmt.Sprintf("\"id\":\"guid-%d\",\n", i)
		dataField += fmt.Sprintln(`"name":"` + tagName + `"`)
		dataField += fmt.Sprintln("}")
		if i < len(tagNames)-1 {
			dataField += ","
		}
	}
	dataField += fmt.Sprintln("]")
	return `{
		"success": true,
		"data": ` + dataField + `
		}`
}

func ServerErrorJSON(errorShort, message string) string {
	return `{
		"success": false,
		"error": "` + errorShort + `",
		"message": "` + message + `"
	}`
}

var _ = Describe("Tags", func() {
	BeforeEach(func() {
		tagsCache = make(map[string]*Tag)
	})

	Describe("GetTagById", func() {
		It("returns tasks with same tags, as words", func() {
			var tagID = "d268201e-c926-4a32-8ac1-7ca570c26b45"
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tags/d268201e-c926-4a32-8ac1-7ca570c26b45"),
					ghttp.RespondWith(200, ChoresTag),
				),
			)

			tag := habitapi.GetTagByID(tagID)

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

	Describe("AddTag", func() {
		Context("given a tag name", func() {
			var tagName string
			BeforeEach(func() {
				tagName = "tag-name"
			})
			It("will return with a valid id", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v3/tags"),
						ghttp.RespondWith(201, ValidTagJSON(tagName)),
					),
				)

				tag, _ := habitapi.AddTag(tagName)
				Expect(tag.ID).ShouldNot(BeEmpty())
				Expect(tag.Name).Should(Equal(tagName))
			})
			Context("given an error is returned by server", func() {
				It("will return with an err", func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("POST", "/v3/tags"),
							ghttp.RespondWith(201, ServerErrorJSON("BadRequest", "User Validation failed")),
						),
					)

					_, err := habitapi.AddTag(tagName)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})

	Describe("GetTag", func() {
		Context("given a tag name", func() {
			var tagName string
			BeforeEach(func() {
				tagName = "tag-name"
			})
			Context("tag exists on server", func() {
				BeforeEach(func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/v3/tags"),
							ghttp.RespondWith(201, CollectionTagsJSON([]string{tagName})),
						),
					)
				})

				It("will return with a valid id and same name", func() {
					tag, _ := habitapi.GetTag(tagName)
					Expect(tag.ID).ShouldNot(BeEmpty())
					Expect(tag.Name).Should(Equal(tagName))
				})

				It("will use cache when tagName is called twice", func() {
					habitapi.GetTag(tagName)
					tag, _ := habitapi.GetTag(tagName)
					Expect(tag.ID).ShouldNot(BeEmpty())
					Expect(tag.Name).Should(Equal(tagName))
				})
			})
			Context("tag does not exist on server", func() {
				BeforeEach(func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/v3/tags"),
							ghttp.RespondWith(404, CollectionTagsJSON([]string{})),
						),
					)
				})

				It("will return with a empty tag", func() {
					tag, _ := habitapi.GetTag(tagName)
					Expect(tag.ID).Should(BeEmpty())
					Expect(tag.Name).Should(BeEmpty())
				})
			})
		})
	})

	Describe("GetTags", func() {
		Context("two tags exists on server", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v3/tags"),
						ghttp.RespondWith(200, CollectionTagsJSON([]string{"tag1", "tag2"})),
					),
				)
			})
			It("will return a list of corresponding tags", func() {
				tags, _ := habitapi.GetTags()
				Expect(len(tags)).To(BeNumerically("==", 2))
				Expect(tags[0]).To(BeEquivalentTo(Tag{Name: "tag1", ID: "guid-0"}))
				Expect(tags[1]).To(BeEquivalentTo(Tag{Name: "tag2", ID: "guid-1"}))
			})
			It("will use cache when GetTag is called after GetTags", func() {
				habitapi.GetTags()
				tag, _ := habitapi.GetTag("tag1")
				Expect(tag).To(BeEquivalentTo(Tag{Name: "tag1", ID: "guid-0"}))
			})
		})
	})
})
