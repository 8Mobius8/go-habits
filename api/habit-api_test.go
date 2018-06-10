package api_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/8Mobius8/go-habits/api"
)

var _ = Describe("Habitica API Router", func() {

	Context("when making regular requests", func() {

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

		Describe("when recieving okay headers from api", func() {
			okStatuses := []int{
				http.StatusOK, http.StatusCreated,
			}
			for _, status := range okStatuses {
				It(http.StatusText(status)+" will return with byte[] array of response from route", func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/v3/resource"),
							ghttp.RespondWith(status, `{"resource":"somebytes"}`),
						),
					)

					var aString struct {
						Resource string
					}
					err := habitapi.Get("/resource", &aString)
					Expect(err).ToNot(HaveOccurred())
					Expect(aString.Resource).To(Equal(`somebytes`))
				})
			}
		})

		Describe("when recieving errors from api", func() {
			errorStatuses := []int{
				http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound,
				http.StatusInternalServerError, http.StatusServiceUnavailable,
			}
			for _, errorStatus := range errorStatuses {

				It("will respond with "+http.StatusText(errorStatus)+" error when habitica error when api called failed", func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/v3"),
							ghttp.RespondWith(errorStatus, http.StatusText(errorStatus)),
						),
					)

					err := habitapi.Get("", nil)

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(http.StatusText(errorStatus)))
				})
			}
		})

		Describe("when authenicating user", func() {
			It("get api token", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v3/user/auth/local/login"),
						ghttp.RespondWith(http.StatusOK, `{"success": true,"data": {"id": "guid","apiToken": "token","newUser": false},"appVersion": "4.41.5"}`),
					),
				)

				creds := habitapi.Authenticate("bob", "p4ssw0rd")
				Expect(creds.ID).To(Equal("guid"))
				Expect(creds.APIToken).To(Equal("token"))
			})
		})

		Describe("when doing a request", func() {
			It("will parse response object", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v3/resource"),
						ghttp.RespondWith(http.StatusOK, []byte(`{"PropertyA":"A","PropertyB":10}`)),
					),
				)

				var aDataModel struct {
					PropertyA string
					PropertyB int
				}

				req, _ := http.NewRequest("GET", server.URL()+"/v3/resource", nil)

				err := habitapi.Do(req, &aDataModel)

				Expect(err).ToNot(HaveOccurred())
				Expect(aDataModel.PropertyA).To(Equal("A"))
				Expect(aDataModel.PropertyB).To(Equal(10))
			})
		})

		Describe("when doing a GET", func() {
			It("will parse response object", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v3/resource"),
						ghttp.RespondWith(http.StatusOK, []byte(`{"PropertyA":"A","PropertyB":10}`)),
					),
				)

				var aDataModel struct {
					PropertyA string
					PropertyB int
				}

				err := habitapi.Get("/resource", &aDataModel)

				Expect(err).ToNot(HaveOccurred())
				Expect(aDataModel.PropertyA).To(Equal("A"))
				Expect(aDataModel.PropertyB).To(Equal(10))
			})
		})

		Describe("when doing a POST", func() {
			It("will parse request object and return parsed response object", func() {
				type requestType struct {
					InputA string
					InputB int
				}
				requestBody := requestType{"penis", 25}

				body, _ := json.Marshal(requestBody)

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/v3/resource"),
						ghttp.VerifyBody(body),
						ghttp.RespondWith(http.StatusOK, []byte(`{"PropertyA":"A","PropertyB":10}`)),
					),
				)

				var aDataModel struct {
					PropertyA string
					PropertyB int
				}
				err := habitapi.Post("/resource", requestBody, &aDataModel)

				Expect(err).ToNot(HaveOccurred())
				Expect(aDataModel.PropertyA).To(Equal("A"))
				Expect(aDataModel.PropertyB).To(Equal(10))
			})
		})
	})
})
