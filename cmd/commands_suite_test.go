package cmd

import (
	"testing"

	"github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var server *ghttp.Server

func TestCmd(t *testing.T) {
	BeforeEach(func() {
		server = ghttp.NewServer()
		habitsServer = api.NewHabiticaAPI(nil, server.URL())
		habitsServerURL = server.URL()
	})

	AfterEach(func() {
		server.Close()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "commands package unit tests")
}
