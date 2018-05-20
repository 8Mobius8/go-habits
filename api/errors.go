package api

type APIError struct {
	msg  string
	code int
}

func (err *APIError) Error() string {
	return err.msg
}
