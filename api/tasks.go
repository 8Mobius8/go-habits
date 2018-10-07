package api

import "fmt"

// Task is a Habitica task.
type Task struct {
	Order int
	Title string   `json:"text"`
	Tags  []string `json:"tags"`
	ID    string   `json:"id"`
	Type  string   `json:"type"`
}

type TaskType int

const (
	Habit         TaskType = 0
	Daily         TaskType = 1
	Todo          TaskType = 2
	Reward        TaskType = 3
	CompletedTodo TaskType = 4
)

func (tt TaskType) String() string {
	taskTypes := []string{
		"habit",
		"daily",
		"todo",
		"reward",
		"completedTodo",
	}
	if tt < Habit || tt > CompletedTodo {
		return "Unknown"
	}

	return taskTypes[tt]
}

func (tt TaskType) AsUrlParam() string {
	return "type=" + tt.String() + "s"
}

// GetTasks will return tasks from Habitica as authenticated user.
// Should call api.Authenticate() before using this.
func (api *HabiticaAPI) GetTasks(tt TaskType) []Task {
	tasks := api.getTasks(tt)
	addOrder(tasks)
	for i := 0; i < len(tasks); i++ {
		tasks[i].Tags = api.getTags(tasks[i])
	}
	return tasks
}

func (api *HabiticaAPI) getTasks(tt TaskType) []Task {
	var tasks []Task
	url := "/tasks/user"
	if tt.String() != "Unknown" {
		url += "?" + tt.AsUrlParam()
	}
	err := api.Get(url, &tasks)
	if err != nil {
		fmt.Println(err)
	}

	return tasks
}

func addOrder(tasks []Task) {
	for i := 0; i < len(tasks); i++ {
		tasks[i].Order = i + 1
	}
}

// AddTask will create the task on the server using the
// task struct as input.
func (api *HabiticaAPI) AddTask(t Task) (Task, error) {
	isOk, err := isValidTask(t)
	if !isOk {
		return Task{}, err
	}

	task, err := api.addTask(t)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func isValidTask(t Task) (bool, error) {
	if t.Title == "" {
		return false, NewGoHabitsError("task is missing text.", 1, "")
	}
	if t.ID != "" {
		return false, NewGoHabitsError("ID is set. You cannot create a new task with an id.", 1, "")
	}
	return true, nil
}

func (api *HabiticaAPI) addTask(t Task) (Task, error) {
	fmt.Println(t)
	err := api.Post("/tasks/user", t, &t)
	return t, err
}

func NewTask(tt TaskType) Task {
	return Task{
		Type: tt.String(),
	}
}
