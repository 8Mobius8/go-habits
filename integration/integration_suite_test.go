package integration_test

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/8Mobius8/go-habits/api"
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
	GOHABITS = "go-habits"
)

var HABITICA_API string
var BUILD_VERSION string
var apiClient *api.HabiticaAPI
var apiToken string
var apiID string

var _ = Describe("go-habits", func() {
	var command *exec.Cmd

	BeforeSuite(func() {
		var exists bool
		HABITICA_API, exists = os.LookupEnv("SERVER")
		Ω(exists).ShouldNot(BeFalse())
		Ω(HABITICA_API).ShouldNot(BeEmpty())

		BUILD_VERSION, exists = os.LookupEnv("BUILD_VERSION")
		Ω(exists).ShouldNot(BeFalse())
		Ω(BUILD_VERSION).ShouldNot(BeEmpty())

		apiClient = api.NewHabiticaAPI(nil, HABITICA_API)
		RegisterUser(HABITICA_API, userName, password, email)
	})

	AfterSuite(func() {
		SaveAPIToken(HABITICA_API, userName, password)
		DeleteUser(HABITICA_API, userName, password, "go-habits integration test")
	})

	Describe("status command", func() {
		It("exits with a zero", func() {
			command = exec.Command(GOHABITS, "status")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gbytes.Say(`Habitica is reachable, GO catch all those pets!`))
			Eventually(session).Should(gexec.Exit(0))
		})
	})

	Describe("login command", func() {
		It("exits with a zero", func() {
			command = exec.Command(GOHABITS, "login")
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

	Describe("version command", func() {
		It("displays full version information", func() {
			command = exec.Command(GOHABITS, "version")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gbytes.Say(
				`go-habits version ` + BUILD_VERSION,
			))
			Eventually(session).Should(gexec.Exit(0))
		})
	})
})

func RegisterUser(serverUri, username, password, email string) {
	// payload := `{"username":"` + username + `","email":"` + email + `","password":"` + password + `","confirmPassword":"` + password + `"}`
	payload := struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
	}{
		username,
		password,
		email,
		password,
	}
	//resp, err := http.Post(serverUri+"/v3/user/auth/local/register", "application/json", bytes.NewBuffer([]byte(payload)))
	err := apiClient.Post("/user/auth/local/register", &payload, nil)
	Ω(err).ShouldNot(HaveOccurred())
}

func SaveAPIToken(serverUri, username, password string) {
	creds := apiClient.Authenticate(username, password)
	apiToken = creds.APIToken
	apiID = creds.ID
}

func DeleteUser(serverUri, username, password, feedback string) {
	payload := `{"password":"` + password + `","feedback":"` + feedback + `"}`
	req, err := http.NewRequest("DELETE", serverUri+"/v3/user", bytes.NewBuffer([]byte(payload)))
	Ω(err).ShouldNot(HaveOccurred())
	err = apiClient.Do(req, nil)
	Ω(err).ShouldNot(HaveOccurred())
}
