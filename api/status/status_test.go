package status_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/8Mobius8/go-habits/api/status"
)

var _ = Describe("Status Habitica API", func() {
	It("returns the funnies when Habitica is reachable.", func() {
		resp := StatusResponse{
			Success: true,
			Data: struct {
				Status string
			}{
				Status: "up",
			},
		}

		Expect(HabiticaStatusMessage(resp)).To(Equal("Habitica is reachable, GO catch all those pets!"))
	})

	It("returns the sad when Habitica is unreachable.", func() {
		resp := StatusResponse{
			Success: true,
			Data: struct {
				Status string
			}{
				Status: "down",
			},
		}

		Expect(HabiticaStatusMessage(resp)).To(Equal(":( Habitica is unreachable."))
	})

	It("returns the sad when Habitica is unreachable.", func() {
		resp := StatusResponse{
			Success: false,
		}

		Expect(HabiticaStatusMessage(resp)).To(Equal(":( Habitica is unreachable."))
	})
})
