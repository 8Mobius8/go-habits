package api

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
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

	Describe("when authenicating with user", func() {
		var expectedCreds UserToken

		BeforeEach(func() {
			expectedCreds = defaultUserCredientials()
			server.AppendHandlers(
				defaultAuthHandlers(expectedCreds.ID, expectedCreds.APIToken),
			)
		})

		It("get api token", func() {
			creds := habitapi.Authenticate("bob", "p4ssw0rd")
			Expect(creds.ID).To(Equal(expectedCreds.ID))
			Expect(creds.APIToken).To(Equal(expectedCreds.APIToken))
		})

		It("will call request with authenticated headers", func() {
			server.AppendHandlers(
				VerifyAuthHeaders(expectedCreds.ID, expectedCreds.APIToken),
			)

			type empty struct{}
			var e empty
			habitapi.Authenticate("bob", "p4ssw0rd")
			habitapi.Get("/empty", &e)
		})
	})

	Describe("when user has not been authenticated", func() {
		It("will call request without authenticated headers", func() {
			server.AppendHandlers(
				ghttp.VerifyRequest("GET", "/v3/empty"),
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

func defaultUserCredientials() UserToken {
	return UserToken{
		ID:       "id",
		APIToken: "token",
	}
}

func defaultAuthHandlers(id string, token string) http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyRequest("POST", "/v3/user/auth/local/login"),
		RespondAuthWithOk(id, token),
	)
}

func RespondAuthWithOk(id string, token string) http.HandlerFunc {
	return ghttp.RespondWith(http.StatusOK, `{ "success": true,"data":{"id": "`+id+`","apiToken": "`+token+`","newUser": false},"appVersion": "4.41.5"}`)
}

func VerifyAuthHeaders(id string, token string) http.HandlerFunc {
	return ghttp.CombineHandlers(
		ghttp.VerifyHeader(http.Header{
			"x-api-user": []string{id},
		}),
		ghttp.VerifyHeader(http.Header{
			"x-api-key": []string{token},
		}),
	)
}
