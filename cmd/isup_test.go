package cmd_test

import (
	api "github.com/8Mobius8/go-habits/api"
	. "github.com/8Mobius8/go-habits/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Isup", func() {
	It("returns the funnies when Habitica is reachable.", func() {
		resp := api.StatusResponse{
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
		resp := api.StatusResponse{
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
		resp := api.StatusResponse{
			Success: false,
		}

		Expect(IsUpMessage(resp)).To(Equal(":( Habitica is unreachable."))
	})
})
