package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits login", func() {
	It("exits with a zero", func() {
		session, stdin := GoHabitsWithStdin("login")
		defer stdin.Close()

		Eventually(session).Should(gbytes.Say(`Username:`))
		stdin.Write([]byte(userName + "\n"))

		Eventually(session).Should(gbytes.Say(`Password:`))
		stdin.Write([]byte(password + "\n"))

		Eventually(session).Should(gbytes.Say(`Didn't find config file. Create one at ~/.go-habits.yaml to save api key for later use.`))
		Eventually(session).Should(gexec.Exit(0))
	})
})
