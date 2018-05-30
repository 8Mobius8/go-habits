package api

// GoHabitsError is simple error struct that holds an error code as int
type GoHabitsError struct {
	msg  string
	code int
}

func (err *GoHabitsError) Error() string {
	return err.msg
}
