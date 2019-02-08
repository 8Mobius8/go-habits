package commands

import (
	"fmt"
	"io"

	api "github.com/8Mobius8/go-habits/api"
)

// StatusServer is an interface for habitica server that has
// an Server status
type StatusServer interface {
	GetServerStatus() (api.Status, error)
	GetHostURL() string
}

// Status or `go-habits status` allows habiters to check if the
// habitica server api is up and running.
func Status(out io.Writer, server StatusServer) error {
	fmt.Fprintln(out, "Using "+server.GetHostURL()+" as api server")
	res, err := server.GetServerStatus()
	if err != nil {
		fmt.Fprintln(out, err)
		ghe, ok := err.(*api.GoHabitsError)
		if ok {
			ghe.StatusCode = 5
			err = ghe
		}
	}

	fmt.Fprintln(out, StatusMessage(res))
	return err
}

// StatusMessage returns text based on Status message
func StatusMessage(resp api.Status) string {
	if resp.Status != "up" {
		return ":( Habitica is unreachable."
	}

	return "Habitica is reachable, GO catch all those pets!"
}
