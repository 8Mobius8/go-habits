package api_test

import (
	"math/rand"

	. "github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Status()", func() {
	Context("when given up status from server", func() {
		It("returns a Status as 'up'", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/status"),
					ghttp.RespondWith(200, UpServerStatus),
				),
			)

			s, err := habitapi.GetServerStatus()
			Expect(err).NotTo(HaveOccurred())
			Expect(s).To(BeEquivalentTo(Status{"up"}))
		})
	})

	Context("when given invalid status from server", func() {
		DescribeTable("returns a Status as 'down'",
			func(responseBody string, expectedStatus Status) {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v3/status"),
						ghttp.RespondWith(200, responseBody),
					),
				)

				s, err := habitapi.GetServerStatus()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(s).To(BeEquivalentTo(expectedStatus))
			},
			Entry("Server returns down", DownServerStatus, Status{"down"}),
			Entry("Server return something else", RandomServerStatus, Status{"down"}),
			Entry("Server doesn't return a `status` object", EmptyServerstatus, Status{"down"}),
			Entry("Server doesn't return a `status` object", EmptyServerstatus, Status{"down"}),
		)
	})

	Context("when given an error from server", func() {
		DescribeTable("returns a down status and error",
			func(statusCode int, responseBody string, expectedStatus Status) {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v3/status"),
						ghttp.RespondWith(statusCode, responseBody),
					),
				)

				s, err := habitapi.GetServerStatus()
				Expect(err).Should(HaveOccurred())
				Expect(s).To(BeEquivalentTo(expectedStatus))
			},
			Entry("Server returns down", 500, DownServerStatus, Status{"down"}),
			Entry("Server returns down", 500, ``, Status{"down"}),
		)
	})
})

var (
	UpServerStatus = `{
		"success": true,
		"data": {
			"status": "up"
		}
	}`
	DownServerStatus = `{
		"success": true,
		"data": {
			"status": "down"
		}
	}`
	RandomServerStatus = `{
		"success": true,
		"data": {
			"status": "` + string(rand.Int63n(5000)) + `"
		}
	}`
	EmptyServerstatus = `{
		"success": true,
		"data": {}
	}`
)
