package api

import "fmt"

// Todo is a Habitica task or todo.
type Todo struct {
	Order int
	Title string   `json:"text"`
	Tags  []string `json:"tags"`
	ID    string   `json:"id"`
}

// GetTodos will return todos from Habitica as authenticated user.
// Should call api.Authenticate() before using this.
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

func (api *HabiticaAPI) AddTodo(t Todo) Todo {
	task := api.addTask(t)
	return buildTodo(task)
}

func (api *HabiticaAPI) addTask(t Todo) task {
	todoTask := task{"", t.Title, "todo"}
	err := api.Post("/tasks/user", todoTask, &todoTask)
	if err != nil {
		fmt.Println(err)
	}
	return todoTask
}

func buildTodo(t task) Todo {
	todo := Todo{}
	todo.ID = t.ID
	todo.Title = t.Text
	return todo
}

type task struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Type string `json:"type"`
}
