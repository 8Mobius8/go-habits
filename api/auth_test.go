package api_test

import (
	"net/http"

	. "github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Habitica API Router", func() {

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
				VerifyAuthHeaders(expectedCreds.ID, expectedCreds.APIToken),
			)

			type empty struct{}
			var e empty
			habitapi.Authenticate("bob", "p4ssw0rd")
			habitapi.Get("/empty", &e)
			habitapi.Post("/empty", &e, &e)
		})

		Context("when using Do will do any type of request with authenticated headers", func() {
			methods := []string{
				http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodPost, http.MethodPatch,
			}
			for _, method := range methods {
				It("will do any type of request with authenticated headers", func() {
					habitapi.Authenticate("bob", "p4ssw0rd")
					req, _ := http.NewRequest(method, "/echo", nil)

					type empty struct{}
					var e empty
					habitapi.Do(req, e)
				})
			}

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
			Expect(actualReq.Header).ShouldNot(HaveKey("x-api-user"))
			Expect(actualReq.Header).ShouldNot(HaveKey("x-api-key"))
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
