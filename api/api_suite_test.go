package api_test

import (
	"io"
	"testing"

	. "github.com/8Mobius8/go-habits/api"
	log "github.com/amoghe/distillog"
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
		ginkgoLogger := log.NewStreamLogger("GinkgoLog", noopCloser{GinkgoWriter})
		habitapi = NewHabiticaAPI(nil, server.URL(), ginkgoLogger)

	})

	AfterEach(func() {
		server.Close()
	})

	RunSpecs(t, "api package unit tests")
}

type noopCloser struct {
	io.Writer
}

func (n noopCloser) Close() error { return nil }
