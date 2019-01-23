package flows

import (
	. "github.com/8Mobius8/go-habits/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var Say = gbytes.Say

const (
	flowsUsername = UserName + "flows"
)

var _ = Describe("go-habits add todos with tags", func() {
	BeforeEach(func() {
		RegisterUser(HabiticaAPIURI, flowsUsername, Password, Email)
		SaveAPIToken(HabiticaAPIURI, flowsUsername, Password)
		TouchConfigFile()
		ExpectSuccessfulLogin(flowsUsername, Password)
	})
	AfterEach(func() {
		DeleteUser(HabiticaAPIURI, flowsUsername, Password, "go-habits integration test")
	})

	var s *gexec.Session
	Context("have added a task with tags", func() {
		BeforeEach(func() {
			s = GoHabits("add", "a", "task", "with", "a", "#tag")
			Eventually(s).Should(gexec.Exit(0))
		})
		It("when listing tasks should also have tags", func() {
			s = GoHabits("list")
			Eventually(s).Should(Say("#tag"))
			Eventually(s).Should(gexec.Exit(0))
		})
	})
})
