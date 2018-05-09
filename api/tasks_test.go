package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/8Mobius8/go-habits/api"
)

var _ = Describe("'go-habits tasks'", func() {
	It("list current tasks", func() {
		resp := TasksResponse{
			Success: true,
			Data: []Task{
				Task{
					Text: "Get stuff done!",
				},
			},
		}

		tasks := GetTasks(resp)
		task := tasks[0]
		Expect((task.Text)).To(Equal("Get stuff done!"))
	})
})
