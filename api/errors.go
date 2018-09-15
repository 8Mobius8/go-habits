package api

// GoHabitsError is simple error struct that holds an error code as int
type GoHabitsError struct {
	msg        string
	StatusCode int
	Path       string
}

func (err *GoHabitsError) Error() string {
	return err.msg
}

func NewGoHabitsError(message string, code int, path string) *GoHabitsError {
	return &GoHabitsError{message, code, path}
}
