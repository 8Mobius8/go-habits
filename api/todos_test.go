package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

const (
	DishesTodo = `{"success":"true","notifications":[],
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
	ChoresTodos = `{"success":"true","notifications":[],
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
	ChoresWorkTodos = `{"success":"true","notifications":[],
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

	NewTodo = `{"success":"true","notifications":[],
		"data":{
				"text":"New Todo",
				"type":"todo",
				"tags":[],
				"value":10,
				"priority":1,
				"attribute":"str",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-1"
		}}`
)

var _ = Describe("Todos", func() {

	Describe("getTodos", func() {
		It("returns at least one valid todo", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/tasks/user"),
					ghttp.RespondWith(200, DishesTodo),
				),
			)

			todos := habitapi.getTodos()

			Expect(len(todos)).Should(BeNumerically("==", 1))
			todo := todos[0]
			Expect(todo.Title).Should(Equal("Clean the dishes"))
		})
	})

	Describe("addOrder", func() {
		It("returns todos with Order as given", func() {
			var todos = []Todo{
				{0, "Clean Dishes", []string{"tag-guid-1"}, ""},
				{0, "Laundry", []string{"tag-guid-1"}, ""},
				{0, "Make bed", []string{"tag-guid-1"}, ""},
			}

			addOrder(todos)

			Expect(len(todos)).Should(BeNumerically("==", 3))
			for i, todo := range todos {
				Expect(todo.Order).Should(BeNumerically("==", i+1))
			}
		})
	})

	Describe("AddTodo", func() {
		It("will return a todo with an new id", func() {
			t := Todo{}
			t.Title = "New Todo"

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v3/tasks/user"),
					ghttp.RespondWith(201, NewTodo),
				),
			)

			t = habitapi.AddTodo(t)
			Expect(t.ID).ShouldNot(BeEmpty())
		})
	})
})
