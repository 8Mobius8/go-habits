package commands_test

import (
	"errors"
	"io"

	"github.com/8Mobius8/go-habits/api"
	cmd "github.com/8Mobius8/go-habits/cmd/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var StatusUpResponse = `{
		"success": true,
		"data": {
			"status": "up"
		},
		"appVersion": "4.58.0"
	}`

var _ = Describe("Status cmd", func() {
	Describe("StatusMessage", func() {
		It("returns the funnies when Habitica is reachable.", func() {
			var resp api.Status
			resp.Status = "up"

			Expect(cmd.StatusMessage(resp)).To(Equal("Habitica is reachable, GO catch all those pets!"))
		})

		It("returns the sad when Habitica is unreachable.", func() {
			var resp api.Status
			resp.Status = "down"

			Expect(cmd.StatusMessage(resp)).To(Equal(":( Habitica is unreachable."))
		})
	})
})

type MockStatusServer struct {
	GetServerStatusFunc func() (api.Status, error)
	GetHostURLFunc      func() string
}

func (m MockStatusServer) GetServerStatus() (api.Status, error) {
	return m.GetServerStatusFunc()
}

func (m MockStatusServer) GetHostURL() string {
	return m.GetHostURLFunc()
}

var _ = Describe("Command integration tests", func() {
	Describe("Status command", func() {
		var out io.Writer
		BeforeEach(func() {
			out = gbytes.NewBuffer()
		})
		It("will print to out to correct string when status is down", func() {
			mockServer := MockStatusServer{
				GetServerStatusFunc: func() (api.Status, error) {
					return api.Status{}, errors.New("Server is not available")
				},
				GetHostURLFunc: func() string {
					return "https://habitica.com/api"
				},
			}

			cmd.Status(out, mockServer)
			Eventually(out).Should(gbytes.Say("Habitica is unreachable."))
		})
		It("will return error with status code = 5 when status is down", func() {
			mockServer := MockStatusServer{
				GetServerStatusFunc: func() (api.Status, error) {
					return api.Status{}, api.NewGoHabitsError("Server is not available", 500, "/status")
				},
				GetHostURLFunc: func() string {
					return "https://habitica.com/api"
				},
			}

			err := cmd.Status(out, mockServer)
			ghe := err.(*api.GoHabitsError)
			Expect(ghe.StatusCode).To(Equal(5))
		})
		It("will print to out to correct string when status is up", func() {
			mockServer := MockStatusServer{
				GetServerStatusFunc: func() (api.Status, error) {
					return api.Status{Status: "up"}, nil
				},
				GetHostURLFunc: func() string {
					return "https://habitica.com/api"
				},
			}

			err := cmd.Status(out, mockServer)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(out).Should(gbytes.Say("Habitica is reachable."))
		})
	})
})
