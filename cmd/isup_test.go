package cmd_test

import (
	api "github.com/8Mobius8/go-habits/api"
	. "github.com/8Mobius8/go-habits/cmd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Isup", func() {
	It("returns the funnies when Habitica is reachable.", func() {
		var resp api.Status
		resp = "up"

		Expect(IsUpMessage(resp)).To(Equal("Habitica is reachable, GO catch all those pets!"))
	})

	It("returns the sad when Habitica is unreachable.", func() {
		var resp api.Status
		resp = "down"

		Expect(IsUpMessage(resp)).To(Equal(":( Habitica is unreachable."))
	})
})
