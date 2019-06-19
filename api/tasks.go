package api

import "time"

// Task is a Habitica task.
type Task struct {
	Order     int
	Title     string   `json:"text"`
	Tags      []string `json:"tags"`
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Completed bool     `json:"completed"`
	date      string   `json:"date"`
	DueDate   time.Time
}

// NewTask creates a new task of a particular Task type.
func NewTask(title string, tt TaskType) Task {
	return Task{
		Title: title,
		Type:  tt.String(),
	}
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
		tasks[i].Tags = api.getTagsByTask(tasks[i])
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
		api.logger.Errorln(err)
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

	tagIDs := []string{}
	for _, tagName := range t.Tags {
		tag, err := api.GetTag(tagName)
		if err != nil {
			return Task{}, err
		}
		if tag.ID == "" {
			_, err = api.AddTag(tagName)
			if err != nil {
				return Task{}, err
			}
		}
		tag, err = api.GetTag(tagName)
		if err != nil {
			return Task{}, err
		}
		tagIDs = append(tagIDs, tag.ID)
	}
	t.Tags = tagIDs

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

type ScoreUpDelta struct {
	Mp  float64 `json:"mp"`
	Exp float64 `json:"exp"`
	Gp  float64 `json:"gp"`
}

// ScoreTaskUp calls api to score a task up. Equvilant to marking the task as
// completed. This results in a experience, gold, and other reward gain.
func (api *HabiticaAPI) ScoreTaskUp(t Task) (ScoreUpDelta, error) {
	resp := ScoreUpDelta{}
	if t.ID == "" {
		return resp, NewGoHabitsError("Task id is empty", 1, "")
	}
	err := api.Post("/tasks/"+t.ID+"/score/up", struct{}{}, &resp)
	return resp, err
}

// DeleteTask removes a task from the habits server. Does not complete, cause
// rewards for removing.
func (api *HabiticaAPI) DeleteTask(t Task) error {
	if t.ID == "" {
		return NewGoHabitsError("Task id is empty", 1, "")
	}
	return api.Delete("/tasks/" + t.ID)
}

// HabiticaDateExample is to be used date.Format to be able to format dates to and from
// Habitica API timestamps and golang native times.
var HabiticaDateExample = "2019-02-15T00:54:00.000Z"

// SetDueDate sets the due date for a task using the given date as a time struct
func (api *HabiticaAPI) SetDueDate(t Task, date time.Time) (Task, error) {
	t.DueDate = date
	taskUpdate := struct {
		Date string `json:"date"`
	}{
		date.Format(HabiticaDateExample),
	}
	err := api.Put("/tasks/"+t.ID, taskUpdate, t)
	if err != nil {
		return Task{}, err
	}
	return t, nil
}
