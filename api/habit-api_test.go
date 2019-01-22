package api_test

import (
	"encoding/json"
	"math/rand"
	"net/http"

	. "github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Habitica API Router", func() {

	Describe("when creating a new Habitica Client", func() {
		Context("And no hostUrl or client is given", func() {
			It("will set the client to use `https://habitica.com/api` as api route", func() {
				client := NewHabiticaAPI(nil, "", nil)
				Expect(client.GetHostURL()).To(Equal("https://habitica.com/api"))
			})
		})
		Context("And a hostURL is given", func() {
			It("will set the client to use the given string as api route", func() {
				randomAPIRoute := randomString(10)
				client := NewHabiticaAPI(nil, randomAPIRoute, nil)
				Expect(client.GetHostURL()).To(Equal(randomAPIRoute))
			})
		})
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
						ghttp.RespondWith(status, `{"data":{"resource":"somebytes"}}`),
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

	Describe("when doing a request", func() {
		It("will parse response object", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/resource"),
					ghttp.RespondWith(http.StatusOK, []byte(`{"data":{"PropertyA":"A","PropertyB":10}}`)),
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
		It("and a dial error occurs returns gohabits error ", func() {
			req, _ := http.NewRequest("GET", "http://notaccessibl.io/v3/resource", nil)
			var e struct{}
			err := habitapi.Do(req, &e)
			Expect(err).To(HaveOccurred())
			ghe := err.(*GoHabitsError)
			Expect(ghe.StatusCode).To(Equal(1))
		})
	})

	Describe("when doing a GET", func() {
		It("will parse response object", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/resource"),
					ghttp.RespondWith(http.StatusOK, []byte(`{"data":{"PropertyA":"A","PropertyB":10}}`)),
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
			requestBody := requestType{"test", 25}

			body, _ := json.Marshal(requestBody)

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v3/resource"),
					ghttp.VerifyBody(body),
					ghttp.RespondWith(http.StatusOK, []byte(`{"data":{"PropertyA":"A","PropertyB":10}}`)),
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

	Describe("when doing a DELETE", func() {
		It("will parase request object and return the parse response object", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/v3/resource"),
					ghttp.RespondWith(http.StatusOK, nil),
				),
			)

			err := habitapi.Delete("/resource")

			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when recieving errors from API", func() {
		Describe("and errors are HTTP status errors", func() {
			errorStatuses := []int{
				http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound,
				http.StatusInternalServerError, http.StatusServiceUnavailable,
			}
			for _, errorStatus := range errorStatuses {

				It("will respond with go-habits code error when recieving '"+http.StatusText(errorStatus)+"'", func() {
					server.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", "/v3"),
							ghttp.RespondWith(errorStatus, ``),
						),
					)

					err := habitapi.Get("", nil)
					habitErr := err.(*GoHabitsError)
					Expect(err).Should(HaveOccurred())
					Expect(habitErr.StatusCode).To(Equal(errorStatus))
				})
			}
		})

		Describe("and errors are Habitica API errors", func() {
			It("will response with Habitica error message Go-Habits Status code", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						// https://habitica.com/apidoc/#api-User-UserGet
						ghttp.VerifyRequest("GET", "/v3/user"),
						ghttp.RespondWith(http.StatusUnauthorized, `{"success":false,"error":"NotAuthorized","message":"Missing authentication headers."}`),
					),
				)
				type empty struct{}
				var e empty
				err := habitapi.Get("/user", e)
				habitErr := err.(*GoHabitsError)

				Expect(habitErr).To(HaveOccurred())
				Expect(habitErr.StatusCode).To(Equal(401))
				Expect(habitErr.Error()).To(Equal("Missing authentication headers."))
				Expect(habitErr.Path).To(Equal("/v3/user"))
			})
		})
	})

	Context("when recieving unmarshable content from API", func() {
		It("the Do function will return an error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/something"),
					ghttp.RespondWith(http.StatusOK, `{"data":{"resource":"somebrokenbytes}`),
				),
			)
			req, _ := http.NewRequest("GET", server.URL()+"/v3/something", nil)
			err := habitapi.Do(req, nil)

			Expect(err).To(HaveOccurred())
			habitErr := err.(*GoHabitsError)
			Expect(habitErr.Error()).To(And(ContainSubstring("Unmarshal"), ContainSubstring("failed")))
			Expect(habitErr.Path).To(Equal("/v3/something"))
		})
	})
})

const pool = "0987654321abcdefghijklmnopqrstuvwxyz"

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}
	return string(bytes)
}
