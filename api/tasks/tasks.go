package tasks

type TasksResponse struct {
	Success bool
	Data    []Task
}

type Task struct {
	Text string
}

func GetTasks(resp TasksResponse) []Task {
	return resp.Data
}
