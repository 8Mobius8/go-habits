package api

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var (
	server   *ghttp.Server
	habitapi *HabiticaAPI
)

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeEach(func() {
		server = ghttp.NewServer()
		habitapi = NewHabiticaAPI(nil, server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	RunSpecs(t, "api package unit tests")
}
