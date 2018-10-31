package flows

import (
	"io"

	. "github.com/8Mobius8/go-habits/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var RegularFlow = struct {
	username string
	password string
	email    string
}{
	"gh-regularuser",
	"gh-regularpassword",
	"gh-regularuser@go-habits.io",
}

var Say = gbytes.Say

var _ = Describe("CLI User usage flows", func() {
	var s *gexec.Session
	var in io.WriteCloser
	BeforeEach(func() {
		RemoveConfigFile()
	})
	AfterEach(func() {
		SaveAPIToken(HabiticaAPIURI, RegularFlow.username, RegularFlow.password)
		DeleteUser(HabiticaAPIURI, RegularFlow.username, RegularFlow.password, "test")
		RemoveConfigFile()
	})
	It("A typical first usage", func() {
		GoHabitsEventuallyExitSafely := func(args ...string) {
			s := GoHabits(args...)
			Eventually(s).Should(gexec.Exit(0))
		}

		By("Reset envirnoment")
		RemoveConfigFile()

		By("first sign-in")
		// CLI user gets usage for `login` command.
		s = GoHabits("login", "-h")
		Eventually(s).Should(Say(`Authenicates with Habits server and saves api token in config file`))
		Eventually(s).Should(Say(`You will need create an account on https://habitica.com`))
		Eventually(s).Should(gexec.Exit(0))

		// CLI user will have made an account on habitica ui.
		RegisterUser(HabiticaAPIURI, RegularFlow.username, RegularFlow.password, RegularFlow.email)

		// CLI user will run login, pass-in their username and password
		// and CLI will say login successful and config file has been updated
		s, in = GoHabitsWithStdin("login")
		EventuallyLogin(s, in, RegularFlow.username, RegularFlow.password)
		Eventually(s).Should(Say(`Login Successful`))
		// CLI user will run login and see that a config file was created for them.
		Eventually(s).Should(Say(`Didn't find config file`))
		Eventually(s).Should(Say(`Created a new config file at .*\.go-habits`))
		Eventually(s).Should(gexec.Exit(0))
		in.Close()

		By("add tasks with and without tags and complete one")
		// CLI user adds tasks
		GoHabitsEventuallyExitSafely("add", "Buy groceries from Whole Paycheck")
		GoHabitsEventuallyExitSafely("add", "Write blog post on go-habits evolution")

		// CLI user adds tasks with tags
		GoHabitsEventuallyExitSafely("add", "Complain to TPM about standups being too big #work")
		GoHabitsEventuallyExitSafely("add", "Rewrite component in react #work")

		// CLI user lists tasks and sees all tasks with tags just created
		s = GoHabits("list")
		Eventually(s).Should(Say("1.*Rewrite component in react #work"))
		Eventually(s).Should(Say("2.*Complain to TPM about standups being too big #work"))
		Eventually(s).Should(Say("3.*Write blog post on go-habits evolution"))
		Eventually(s).Should(Say("4.*Buy groceries from Whole Paycheck"))
		Eventually(s).Should(gexec.Exit(0))

		// CLI user completes a task
		s = GoHabits("complete", "4")
		Eventually(s).Should(Say("4.*X.*Buy groceries from Whole Paycheck"))

		// CLI user lists their tasks to see what to work on next
		s = GoHabits("list")
		Eventually(s).Should(Say("1.*Rewrite component in react #work"))
		Eventually(s).Should(Say("2.*Complain to TPM about standups being too big #work"))
		Eventually(s).Should(Say("3.*Write blog post on go-habits evolution"))
		Eventually(s).Should(gexec.Exit(0))
	})
})
