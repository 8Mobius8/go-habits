package api

import (
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	DishesTask = `{"success":"true","notifications":[],
		"data":[
			{
				"text":"Clean the dishes",
				"type":"todo",
				"tags":["d268201e-c926-4a32-8ac1-7ca570c26b45"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"8a9d461b-f5eb-4a16-97d3-c03380c422a3"
				}]
		}`
	ChoresTasks = `{"success":"true","notifications":[],
		"data":[
			{
				"text":"Clean the dishes",
				"type":"todo",
				"tags":["d268201e-c926-4a32-8ac1-7ca570c26b45"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-1"
			},{
				"text":"Clean Bathroom",
				"type":"todo",
				"tags":["d268201e-c926-4a32-8ac1-7ca570c26b45"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-2"
			},{
				"text":"Laundry",
				"type":"todo",
				"tags":["d268201e-c926-4a32-8ac1-7ca570c26b45"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-3"
			}]
		}`
	ChoresTag = `{"success": true,"notifications": [],
			"data": {
				"id": "d268201e-c926-4a32-8ac1-7ca570c26b45",
				"name": "chores"
			}
		}`
	WorkTag = `{"success": true,"notifications": [],
		"data": {
			"id": "d268201e-26c9-32a4-a81c-570c7ca4526b",
			"name": "work"
		}
		}`
	ChoresWorkTasks = `{"success":"true","notifications":[],
		"data":[
			{
				"text":"Clean the dishes",
				"type":"todo",
				"tags":["d268201e-c926-4a32-8ac1-7ca570c26b45"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-1"
			},{
				"text":"Clean Bathroom",
				"type":"todo",
				"tags":["d268201e-c926-4a32-8ac1-7ca570c26b45"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-2"
			},{
				"text":"Laundry",
				"type":"todo",
				"tags":["d268201e-c926-4a32-8ac1-7ca570c26b45"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-3"
			},{
				"text":"Schedule Backlog grooming",
				"type":"todo",
				"tags":["d268201e-26c9-32a4-a81c-570c7ca4526b"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-4"
			},{
				"text":"Email Infrastruce team about expiring SSL certs",
				"type":"todo",
				"tags":["d268201e-26c9-32a4-a81c-570c7ca4526b"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-5"
			},{
				"text":"Email Infrastruce team about expired certs",
				"type":"todo",
				"tags":["d268201e-26c9-32a4-a81c-570c7ca4526b"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-6"
			}]
		}`

	ValidTask = `{"success":"true","notifications":[],
		"data":{
				"text":"Valid Todo Title",
				"type":"todo",
				"tags":["valid", "test"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-1"
		}}`
)

var _ = Describe("Tasks", func() {

	Describe("GetTasks", func() {
		Context("when given 'todo' type parameter", func() {
			It("calls server for todo type tasks for user", func() {
				server.AppendHandlers(
					ghttp.VerifyRequest("GET", "/v3/tasks/user", "type=todos"),
				)
				habitapi.GetTasks(TodoType)
			})
		})
	})

	Describe("addOrder", func() {
		It("returns tasks with Order as given", func() {
			var tasks = []Task{
				{0, "Clean Dishes", []string{"tag-guid-1"}, "", "todo"},
				{0, "Laundry", []string{"tag-guid-1"}, "", "todo"},
				{0, "Make bed", []string{"tag-guid-1"}, "", "todo"},
			}

			addOrder(tasks)

			Expect(len(tasks)).Should(BeNumerically("==", 3))
			for i, task := range tasks {
				Expect(task.Order).Should(BeNumerically("==", i+1))
			}
		})
	})

	Describe("AddTask", func() {
		Context("given valid task", func() {
			var t Task
			BeforeEach(func() {
				t = Task{}
				t.Title = "Valid Todo Title"
				t.Tags = []string{"valid", "test"}
			})
			It("will return a task with an new id", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v3/tasks/user"),
						ghttp.RespondWith(201, ValidTask),
					),
				)

				task, _ := habitapi.AddTask(t)
				Expect(task.ID).ShouldNot(BeEmpty())
			})
			It("will return a task with same title", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v3/tasks/user"),
						ghttp.RespondWith(201, ValidTask),
					),
				)

				task, _ := habitapi.AddTask(t)
				Expect(task.Title).Should(Equal(t.Title))
			})
			It("will return a task with tags names", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v3/tasks/user"),
						ghttp.RespondWith(201, ValidTask),
					),
				)

				task, _ := habitapi.AddTask(t)
				Expect(task.Tags).Should(Equal([]string{"valid", "test"}))
			})
		})

		Context("given a task that is invalid", func() {
			DescribeTable("will return error and empty task",
				func(t Task) {
					t, err := habitapi.AddTask(t)
					Expect(err).Should(HaveOccurred())
					Expect(t).Should(Equal(Task{}))
				},
				Entry("a task with out a title", Task{Title: ""}),
				Entry("a task with an id", Task{ID: "something"}),
			)
		})
	})
})
