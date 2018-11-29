package cmd

import (
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/8Mobius8/go-habits/api"
	log "github.com/amoghe/distillog"
)

var server *ghttp.Server

func TestCmd(t *testing.T) {
	BeforeEach(func() {
		server = ghttp.NewServer()
		habitsServer = api.NewHabiticaAPI(nil, server.URL(), log.NewNullLogger("test"))
		habitsServerURL = server.URL()
	})

	AfterEach(func() {
		server.Close()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "commands package unit tests")
}

type noopCloser struct {
	io.Writer
}

func (n noopCloser) Close() error { return nil }
