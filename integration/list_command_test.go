package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits list", func() {
	It("actually exists", func() {
		session := GoHabits("list")
		Eventually(session).Should(gexec.Exit(0))
	})
})
