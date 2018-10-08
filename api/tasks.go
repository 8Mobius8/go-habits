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

// TaskType is a int representation of the different types of
// tasks that Habitica has in it's api. Use the standard
// String() method to have a string version.
type TaskType int

// Constant types that should be used when creating or getting
// tasks from Habitica server.
const (
	HabitType         TaskType = 0
	DailyType         TaskType = 1
	TodoType          TaskType = 2
	RewardType        TaskType = 3
	CompletedTodoType TaskType = 4
)

func (tt TaskType) String() string {
	taskTypes := []string{
		"habit",
		"daily",
		"todo",
		"reward",
		"completedTodo",
	}
	if tt < HabitType || tt > CompletedTodoType {
		return "Unknown"
	}

	return taskTypes[tt]
}

func (tt TaskType) asURLParam() string {
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
		url += "?" + tt.asURLParam()
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
// task struct as input. Any new task must have a title
// and type.
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
	err := api.Post("/tasks/user", t, &t)
	return t, err
}

// NewTask creates a new task of a particular Task type.
func NewTask(title string, tt TaskType) Task {
	return Task{
		Title: title,
		Type:  tt.String(),
	}
}
