package api

import "fmt"

// Todo is a Habitica task or todo.
type Todo struct {
	Order int
	Title string   `json:"text"`
	Tags  []string `json:"tags"`
}

func (api *HabiticaAPI) GetTodos() []Todo {
	todos := api.getTodos()
	addOrder(todos)
	for i := 0; i < len(todos); i++ {
		todos[i].Tags = api.getTags(todos[i])
	}
	return todos
}

func (api *HabiticaAPI) getTodos() []Todo {
	var todos []Todo
	err := api.Get("/tasks/user", &todos)
	if err != nil {
		fmt.Println(err)
	}

	return todos
}

func addOrder(todos []Todo) {
	for i := 0; i < len(todos); i++ {
		todos[i].Order = i + 1
	}
}
