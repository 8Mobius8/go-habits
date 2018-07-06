package integration_test

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")

}

const (
	userName = "test"
	password = "test"
	email    = "test@efungus.io"
)

var serverUri string

var _ = Describe("go-habits", func() {
	var (
		goHabitsCLIPath string
		command         *exec.Cmd
	)

	BeforeSuite(func() {
		serverUri = os.Getenv("SERVER")
		var err error
		goHabitsCLIPath, err = gexec.Build("github.com/8mobius8/go-habits")
		Ω(err).ShouldNot(HaveOccurred())

		RegisterUser(userName, password, email)
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	Describe("isup command", func() {
		It("exits with a zero", func() {
			command = exec.Command(goHabitsCLIPath, "isup")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gbytes.Say(`Habitica is reachable, GO catch all those pets!`))
			Eventually(session).Should(gexec.Exit(0))
		})
	})

	Describe("login command", func() {
		It("exits with a zero", func() {
			command = exec.Command(goHabitsCLIPath, "login")
			stdin, err := command.StdinPipe()
			Ω(err).ShouldNot(HaveOccurred())
			defer stdin.Close()

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gbytes.Say(`Username:`))
			stdin.Write([]byte(userName + "\n"))

			Eventually(session).Should(gbytes.Say(`Password:`))
			stdin.Write([]byte(password + "\n"))

			Eventually(session).Should(gbytes.Say(`Didn't find config file. Create one at ~/.go-habits.yaml to save api key for later use.`))
			Eventually(session).Should(gexec.Exit(0))
		})
	})
})

func RegisterUser(username string, password string, email string) {

	payload := `{"username":"` + username + `","email":"` + email + `","password":"` + password + `","confirmPassword":"` + password + `"}`
	resp, err := http.Post(serverUri+"/v3/user/auth/local/register", "application/json", bytes.NewBuffer([]byte(payload)))
	Ω(err).ShouldNot(HaveOccurred())
	Ω(resp.StatusCode).ShouldNot(BeNumerically(">=", 400))
}
