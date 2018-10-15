package cmd

import (
	"io"

	"github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/ghttp"
	"github.com/spf13/cobra"
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

			Expect(StatusMessage(resp)).To(Equal("Habitica is reachable, GO catch all those pets!"))
		})

		It("returns the sad when Habitica is unreachable.", func() {
			var resp api.Status
			resp.Status = "down"

			Expect(StatusMessage(resp)).To(Equal(":( Habitica is unreachable."))
		})
	})

	It("returns string with out habitica reachable.", func() {
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/v3/status"),
				ghttp.RespondWith(200, StatusUpResponse),
			),
		)
		out, err := Status()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(out).To(ContainSubstring("Habitica is reachable, GO catch all those pets!"))
	})

	It("returns err status code 5 when server return down", func() {
		server.SetAllowUnhandledRequests(true)
		out, err := Status()
		Expect(err).Should(HaveOccurred())
		Expect(out).To(ContainSubstring("Habitica is unreachable"))
		ghe := err.(*api.GoHabitsError)
		Expect(ghe.StatusCode).To(Equal(5))
	})
})

var _ = Describe("Command integration tests", func() {
	Describe("Status command", func() {
		var buf *gbytes.Buffer
		var mockCmd *cobra.Command
		var out io.Writer
		BeforeEach(func() {
			buf = gbytes.NewBuffer()
			mockCmd = &cobra.Command{}
			mockCmd.SetOutput(buf)
			out = mockCmd.OutOrStdout()
		})
		It("will print to out to correct string when status is down", func() {
			server.SetAllowUnhandledRequests(true)

			_ = statusCmd.RunE(mockCmd, []string{})
			Eventually(out).Should(gbytes.Say("Habitica is unreachable."))
		})
		It("will return error with status code = 5 when status is down", func() {
			server.SetAllowUnhandledRequests(true)

			err := statusCmd.RunE(mockCmd, []string{})
			ghe := err.(*api.GoHabitsError)
			Expect(ghe.StatusCode).To(Equal(5))
		})
		It("will print to out to correct string when status is up", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/status"),
					ghttp.RespondWith(200, StatusUpResponse),
				),
			)

			_ = statusCmd.RunE(mockCmd, []string{})
			Eventually(out).Should(gbytes.Say("Habitica is reachable."))
		})
	})
})
