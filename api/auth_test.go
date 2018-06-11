package api_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/8Mobius8/go-habits/api"
)

var _ = Describe("Habitica API Router", func() {

	var (
		server   *ghttp.Server
		habitapi *HabiticaAPI
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		habitapi = NewHabiticaAPI(nil, server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("when authenicating user", func() {
		It("get api token", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v3/user/auth/local/login"),
					ghttp.RespondWith(http.StatusOK, `{"success": true,"data":{"id": "guid","apiToken": "token","newUser": false},"appVersion": "4.41.5"}`),
				),
			)

			creds := habitapi.Authenticate("bob", "p4ssw0rd")
			Expect(creds.ID).To(Equal("guid"))
			Expect(creds.APIToken).To(Equal("token"))
		})

		It("will call request with authenticated headers", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v3/user/auth/local/login"),
					ghttp.RespondWith(http.StatusOK, `{"success": true,"data":{"id": "guid","apiToken": "token","newUser": false},"appVersion": "4.41.5"}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/empty"),
					ghttp.VerifyHeader(http.Header{
						"x-api-user": []string{"guid"},
					}),
					ghttp.VerifyHeader(http.Header{
						"x-api-key": []string{"token"},
					}),
				),
			)

			type empty struct{}
			var e empty
			habitapi.Authenticate("bob", "p4ssw0rd")
			req, _ := http.NewRequest("GET", server.URL()+"/v3/empty", nil)
			habitapi.Do(req, &e)
		})
	})

	Describe("when user has not been authenticated", func() {
		It("will call request without authenticated headers", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/empty"),
				),
			)

			type empty struct{}
			var e empty
			req, _ := http.NewRequest("GET", server.URL()+"/v3/empty", nil)
			habitapi.Do(req, &e)

			reqs := server.ReceivedRequests()
			Expect(len(reqs)).To(Equal(1))
			actualReq := reqs[0]
			Expect(actualReq.Header.Get("x-api-user")).To(Equal(""))
			Expect(actualReq.Header.Get("x-api-key")).To(Equal(""))
		})
	})
})
