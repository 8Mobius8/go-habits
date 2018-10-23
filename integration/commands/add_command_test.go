package commands

import (
	. "github.com/8Mobius8/go-habits/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits add command", func() {
	It("actually exists", func() {
		session := GoHabits("add", "-h")
		Eventually(session).Should(gexec.Exit(0))
	})

	Context("Given the user is already signed-in", func() {
		BeforeEach(func() {
			TouchConfigFile()
			expectSuccessfulLogin(UserName, Password)
		})
		AfterEach(func() {
			ResetUser(ApiID, ApiToken)
			RemoveConfigFile()
		})

		Describe("Creating a new todo", func() {
			It("will print newly created task as confirmation", func() {
				s := GoHabits("add", "clean my fishbowl")
				Eventually(s).Should(gbytes.Say("clean my fishbowl"))
				Eventually(s).Should(gexec.Exit(0))
			})

			It("will print newly created task and order as confirmation", func() {
				s := GoHabits("add", "clean my fishbowl")
				Eventually(s).Should(gexec.Exit(0))
				s = GoHabits("add", "make bed")
				Eventually(s).Should(gbytes.Say("1"))
				Eventually(s).Should(gbytes.Say("make bed"))
				Eventually(s).Should(gexec.Exit(0))
			})
		})
	})
})
