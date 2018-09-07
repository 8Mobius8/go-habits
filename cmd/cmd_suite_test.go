package cmd

import (
	"testing"

	"github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cmd package unit tests")
}

var _ = Describe("Status", func() {
	It("returns the funnies when Habitica is reachable.", func() {
		var resp api.Status
		resp.Status = "up"

		Expect(StatusMessage(resp)).To(Equal("Habitica is reachable, GO catch all those pets!"))
	})

	It("returns the sad when Habitica is unreachable.", func() {
		var resp api.Status
		resp.Status = "down"

		Expect(StatusMessage(resp)).To(Equal(":( Habitica is unreachable."))
	})
})
