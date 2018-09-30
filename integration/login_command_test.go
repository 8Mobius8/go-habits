package integration

import (
	"io"
	"io/ioutil"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits login", func() {
	Context("given a configuration file does NOT exist", func() {
		It("exits with message about ", func() {
			s, in := GoHabitsWithStdin("login")
			defer in.Close()

			EventuallyLogin(s, in, userName, password)

			Eventually(s).Should(gbytes.Say(`Didn't find config file. Create one at ~/.go-habits.yml to save api key for later use.`))
			Eventually(s).Should(gexec.Exit(0))
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

			EventuallyLogin(s, in, userName, password)

			Eventually(s).Should(gbytes.Say("Updating config at"))
			Eventually(s).Should(gexec.Exit(0))

			data, err := ioutil.ReadFile(configPath)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(data).Should(MatchRegexp("apitoken: [0-9a-z-]+"))
			Expect(data).Should(MatchRegexp("id: [0-9a-z-]+"))
		})
	})
})

func EventuallyLogin(session *gexec.Session, in io.WriteCloser, username, password string) {
	Eventually(session).Should(gbytes.Say("Username:"))
	in.Write([]byte(username + "\n"))

	Eventually(session).Should(gbytes.Say("Password:"))
	in.Write([]byte(password + "\n"))
}

func TouchConfigFile() string {
	userHomePath, err := homedir.Dir()
	Expect(err).ShouldNot(HaveOccurred())
	path := userHomePath + "/.go-habits.yml"

	err = ioutil.WriteFile(path, []byte("test"), 0644)
	Expect(err).ShouldNot(HaveOccurred())

	return path
}

func RemoveConfigFile() {
	userHomePath, err := homedir.Dir()
	Expect(err).ShouldNot(HaveOccurred())
	path := userHomePath + "/.go-habits.yml"

	err = os.Remove(path)
	Expect(err).ShouldNot(HaveOccurred())
}
