package api_test

import (
	"net/http"
	"time"

	"github.com/8Mobius8/go-habits/api"
	. "github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
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
	CompletedTask = `{"success":"true","notifications":[],
		"data":{
				"text":"Completed Todo Title",
				"type":"todo",
				"tags":["valid", "test"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"completed": true,
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"someId"
		}}`
	// Task with due date set to 2019-02-19 0:00:00.000
	ValidTaskWithDate = `{"success":"true","notifications":[],
		"data":{
				"text":"Valid Todo Title",
				"type":"todo",
				"tags":["valid", "test"],
				"value":10,
				"priority":1,
				"attribute":"str",
				"date":"2019-02-19T00:00:00.000Z",
				"createdAt":"2017-01-07T17:52:09.121Z",
				"updatedAt":"2017-01-11T14:25:32.504Z",
				"id":"chore-id-1"
		}}`
)

var _ = Describe("Tasks", func() {

	DescribeTable("NewTask will create task with corresponding title and string representaiton of type",
		func(title string, tt TaskType, taskTypeString string) {
			task := NewTask(title, tt)
			Expect(task.Title).Should(BeEquivalentTo(title))
			Expect(task.Type).Should(Equal(taskTypeString))
		},
		Entry("Todo task type", "this is a todo", TodoType, "todo"),
		Entry("Habit task type", "this is a habit", HabitType, "habit"),
		Entry("Daily task type", "this is a daily", DailyType, "daily"),
		Entry("Reward task type", "this is a reward", RewardType, "reward"),
		Entry("Completed task type", "this is a completed todo", CompletedTodoType, "completedTodo"),
	)

	Describe("AddTask", func() {
		Context("given valid task", func() {
			var t Task
			BeforeEach(func() {
				t = Task{}
				t.Title = "Valid Todo Title"
				t.Tags = []string{"valid", "test"}

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v3/tags"),
						ghttp.RespondWith(200, CollectionTagsJSON([]string{})),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v3/tags"),
						ghttp.RespondWith(200, ValidTagJSON(t.Tags[0])),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v3/tags"),
						ghttp.RespondWith(200, CollectionTagsJSON([]string{})),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v3/tags"),
						ghttp.RespondWith(200, ValidTagJSON(t.Tags[1])),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v3/tasks/user"),
						ghttp.RespondWith(201, ValidTask),
					),
				)
			})
			AfterEach(func() {
				habitapi.ClearTagCache()
			})
			It("will return a task with an new id", func() {
				task, err := habitapi.AddTask(t)
				Expect(err).ToNot(HaveOccurred())
				Expect(task.ID).ShouldNot(BeEmpty())
			})
			It("will return a task with same title", func() {
				task, err := habitapi.AddTask(t)
				Expect(err).ToNot(HaveOccurred())
				Expect(task.Title).Should(Equal(t.Title))
			})
			It("will return a task with tags names", func() {
				task, err := habitapi.AddTask(t)
				Expect(err).ToNot(HaveOccurred())
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

	Describe("ScoreTaskUp", func() {
		Context("given a task with an id", func() {
			var task Task
			BeforeEach(func() {
				task = Task{ID: "someId"}
			})
			AfterEach(func() {
				Expect(len(server.ReceivedRequests())).Should(BeNumerically(">", 0))
			})
			Context("and id exists on server", func() {
				BeforeEach(func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("POST", "/v3/tasks/"+task.ID+"/score/up"),
						),
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/v3/tasks/"+task.ID),
							ghttp.RespondWith(200, CompletedTask),
						),
					)
				})
				It("will call the server with correct uri", func() {
					err := habitapi.ScoreTaskUp(task)
					Expect(err).ShouldNot(HaveOccurred())
				})
				It("the next get on task will return with completed true", func() {
					err := habitapi.ScoreTaskUp(task)
					Expect(err).ShouldNot(HaveOccurred())

					err = habitapi.Get("/tasks/"+task.ID, &task)
					Expect(err).ShouldNot(HaveOccurred())

					Expect(task.Completed).Should(Equal(true))
				})
			})
			Context("and id does not exist on server", func() {
				It("return with 404 error", func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("POST", "/v3/tasks/"+task.ID+"/score/up"),
							ghttp.RespondWith(404, `{
							"success": false,
							"error": "NotFound",
							"message": "Task not found."
						}`),
						),
					)

					err := habitapi.ScoreTaskUp(task)
					Expect(err).Should(HaveOccurred())
					goHabitsError := err.(*GoHabitsError)
					Expect(goHabitsError.Path).Should(Equal("/v3/tasks/" + task.ID + "/score/up"))
					Expect(goHabitsError.StatusCode).Should(Equal(404))
					Expect(goHabitsError.Error()).Should(ContainSubstring("Task not found."))
				})
			})
		})

		Context("given a task without an id", func() {
			It("will return with error saying no id was given", func() {
				task := Task{}
				err := habitapi.ScoreTaskUp(task)
				Expect(err).Should(HaveOccurred())
				goHabitsError := err.(*GoHabitsError)
				Expect(goHabitsError.Error()).Should(ContainSubstring("Task id is empty"))
			})
		})
	})

	Describe("DeleteTask", func() {
		Context("given the task does not exist on server", func() {
			It("will return with error saying that the task ID not exist", func() {
				task := Task{ID: "someid"}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", "/v3/tasks/"+task.ID),
						ghttp.RespondWith(404, `{
						"success": false,
						"error": "NotFound",
						"message": "Task not found."
						}`),
					),
				)

				err := habitapi.DeleteTask(task)
				Expect(err).Should(HaveOccurred())
				goHabitsError := err.(*GoHabitsError)
				Expect(goHabitsError.Error()).Should(ContainSubstring("Task not found"))
			})
			It("will return with error saying that the task needs an ID in order to delete", func() {
				task := Task{}

				err := habitapi.DeleteTask(task)
				Expect(err).Should(HaveOccurred())
				goHabitsError := err.(*GoHabitsError)
				Expect(goHabitsError.Error()).Should(ContainSubstring("Task id is empty"))
			})
		})
		Context("given the task exists on server", func() {
			It("will return nil on delete", func() {
				task := api.Task{ID: "id"}
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", "/v3/tasks/"+task.ID),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"data": {},
							"notifications": []
						}`),
					),
				)

				err := habitapi.DeleteTask(task)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("given the task exists on the server but not-authorized", func() {
			It("will return with unauthorized error", func() {
				task := api.Task{ID: "id"}
				message := randomString(10)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("DELETE", "/v3/tasks/"+task.ID),
						ghttp.RespondWith(http.StatusUnauthorized, `{
							"success": false,
							"error":"NotAuthorized",
							"message":"`+message+`"
						}`),
					),
				)

				err := habitapi.DeleteTask(task)
				Expect(err).Should(HaveOccurred())
				goHabitsError := err.(*GoHabitsError)
				Expect(goHabitsError.StatusCode).Should(Equal(401))
				Expect(goHabitsError.Error()).Should(Equal(message))
			})
		})
	})

	Describe("SetDueDate", func() {
		Context("given a task has already been created", func() {
			var task api.Task
			var date time.Time
			var dateExample = "2019-02-15T00:54:00.000Z"
			It("will succesfully update tasks due dates via API", func() {
				task = api.Task{ID: "id"}
				date = time.Now()

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT", "/v3/tasks/"+task.ID),
						ghttp.VerifyBody([]byte(`{"date":"`+date.Format(dateExample)+`"}`)),
						ghttp.RespondWith(http.StatusOK, ValidTask),
					),
				)

				_, err := habitapi.SetDueDate(task, time.Now())
				Expect(err).ToNot(HaveOccurred())
			})
			It("will return a task that has the due date set", func() {
				task = api.Task{ID: "id"}
				date = time.Now()

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("PUT", "/v3/tasks/"+task.ID),
						ghttp.VerifyBody([]byte(`{"date":"`+date.Format(dateExample)+`"}`)),
						// Date with
						ghttp.RespondWith(http.StatusOK, ValidTaskWithDate),
					),
				)

				actualTask, err := habitapi.SetDueDate(task, date)
				Expect(err).ToNot(HaveOccurred())
				Expect(actualTask.DueDate).To(BeEquivalentTo(date))
			})
		})
	})
})
