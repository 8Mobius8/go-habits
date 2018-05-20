package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/8Mobius8/go-habits/api/status"
	. "github.com/8Mobius8/go-habits/cmd"
)

var _ = Describe("Isup", func() {
	It("returns the funnies when Habitica is reachable.", func() {
		resp := status.StatusResponse{
			Success: true,
			Data: struct {
				Status string
			}{
				Status: "up",
			},
		}

		Expect(IsUpMessage(resp)).To(Equal("Habitica is reachable, GO catch all those pets!"))
	})

	It("returns the sad when Habitica is unreachable.", func() {
		resp := status.StatusResponse{
			Success: true,
			Data: struct {
				Status string
			}{
				Status: "down",
			},
		}

		Expect(IsUpMessage(resp)).To(Equal(":( Habitica is unreachable."))
	})

	It("returns the sad when Habitica is unreachable.", func() {
		resp := status.StatusResponse{
			Success: false,
		}

		Expect(IsUpMessage(resp)).To(Equal(":( Habitica is unreachable."))
	})
})
