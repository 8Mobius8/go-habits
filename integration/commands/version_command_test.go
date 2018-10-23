package commands

import (
	. "github.com/8Mobius8/go-habits/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits version", func() {
	It("displays full version information", func() {
		session := GoHabits("version")

		Eventually(session).Should(gbytes.Say(
			`go-habits version ` + BuildVersion,
		))
		Eventually(session).Should(gexec.Exit(0))
	})
	It("displays full version information using --verison option", func() {
		session := GoHabits("--version")

		Eventually(session).Should(gbytes.Say(
			`go-habits version ` + BuildVersion,
		))
		Eventually(session).Should(gexec.Exit(0))
	})
})
