package cmd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/8Mobius8/go-habits/cmd"
)

var _ = Describe("go-habits isup", func() {
	It("returns reachable message", func() {
		message := PrintHabiticaAPIStatus()

		Expect(message).To(Equal("Habitica is reachable, GO catch all those pets!"))
	})
})
