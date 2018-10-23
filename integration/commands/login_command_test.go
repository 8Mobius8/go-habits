package commands

import (
	"fmt"
	"io/ioutil"

	. "github.com/8Mobius8/go-habits/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits login", func() {
	Context("given a configuration file does NOT exist", func() {
		It("exits with message saying its creating a new file", func() {
			s, in := GoHabitsWithStdin("login")
			defer in.Close()

			EventuallyLogin(s, in, UserName, Password)

			Eventually(s).Should(gbytes.Say(`Didn't find config file. Creating a new config file at .*\.go-habits.yml`))
			Eventually(s).Should(gexec.Exit(0))
		})
		It("has created a new config file with keys in it", func() {
			s, in := GoHabitsWithStdin("login")
			defer in.Close()
			EventuallyLogin(s, in, UserName, Password)
			Eventually(s).Should(gexec.Exit(0))
			b, err := ioutil.ReadFile("./login-test-config.yml")
			Expect(err).ToNot(HaveOccurred())
			fmt.Fprintln(GinkgoWriter, "\nConfig file contents:")
			fmt.Fprintln(GinkgoWriter, string(b))
			Expect(string(b)).Should(MatchRegexp("apitoken: [a-z0-9-]*"))
			Expect(string(b)).Should(MatchRegexp("id: [a-z0-9-]*"))
		})
	})

	Context("given a configuration file does exist", func() {
		var configPath string

		BeforeEach(func() {
			configPath = TouchConfigFile()
		})
		AfterEach(func() {
			RemoveConfigFile()
		})
		It("updates tokens to file", func() {
			s, in := GoHabitsWithStdin("login")
			defer in.Close()

			EventuallyLogin(s, in, UserName, Password)

			Eventually(s).Should(gbytes.Say("Updating config at"))
			Eventually(s).Should(gexec.Exit(0))

			data, err := ioutil.ReadFile(configPath)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(MatchRegexp("apitoken: [0-9a-z-]+"))
			Expect(data).Should(MatchRegexp("id: [0-9a-z-]+"))
		})
	})
})
