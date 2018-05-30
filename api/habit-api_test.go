package api_test

import (
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
							ghttp.VerifyRequest("GET", "/v3"),
							ghttp.RespondWith(status, "somebytes"),
						),
					)

					res, err := habitapi.Get("")
					Expect(err).ToNot(HaveOccurred())
					Expect(res).To(Equal([]byte(`somebytes`)))
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

					_, err := habitapi.Get("")

					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(http.StatusText(errorStatus)))
				})
			}
		})
	})

	Describe("when parsing a http response", func() {
		var habitapi *HabiticaAPI

		BeforeEach(func() {
			habitapi = NewHabiticaAPI(nil, "")
		})

		It("Will return an struct with fields filled in", func() {
			var aDataModel struct {
				PropertyA string
				PropertyB int
			}
			var stringDataModel = []byte(`{"PropertyA":"A","PropertyB":10}`)

			habitapi.ParseResponse(stringDataModel, &aDataModel)

			Expect(aDataModel.PropertyA).To(Equal("A"))
			Expect(aDataModel.PropertyB).To(Equal(10))
		})
	})
})
