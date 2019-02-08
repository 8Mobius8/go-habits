package commands_test

import (
	"io"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var server *ghttp.Server

func TestCmd(t *testing.T) {
	BeforeEach(func() {
		server = ghttp.NewServer()
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
