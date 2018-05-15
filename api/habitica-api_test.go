package api_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/8Mobius8/go-habits/api"
)

var _ = Describe("habitica api router", func() {
	var (
		server   *ghttp.Server
		habitapi *HabiticaAPI
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		habitapi = NewHabiticaApi(nil, server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("when receving errors from api", func() {
		errorStatuses := []int{
			http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound,
			http.StatusInternalServerError, http.StatusServiceUnavailable,
		}
		for _, errorStatus := range errorStatuses {

			It("will respond with "+http.StatusText(errorStatus)+" error when habitica error when api called failed", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v3/status"),
						ghttp.RespondWith(errorStatus, http.StatusText(errorStatus)),
					),
				)

				_, err := habitapi.Status()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(http.StatusText(errorStatus)))
			})
		}
	})
})
