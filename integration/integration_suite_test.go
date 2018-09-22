package integration

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"

	"github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var HABITICA_API string
var BUILD_VERSION string
var apiClient *api.HabiticaAPI
var apiToken string
var apiID string

// Global go-habits command
var command *exec.Cmd

var _ = BeforeSuite(func() {
	var exists bool
	HABITICA_API, exists = os.LookupEnv("SERVER")
	Ω(exists).ShouldNot(BeFalse())
	Ω(HABITICA_API).ShouldNot(BeEmpty())

	BUILD_VERSION, exists = os.LookupEnv("BUILD_VERSION")
	Ω(exists).ShouldNot(BeFalse())
	Ω(BUILD_VERSION).ShouldNot(BeEmpty())

	apiClient = api.NewHabiticaAPI(nil, HABITICA_API)
	RegisterUser(HABITICA_API, userName, password, email)
	SaveAPIToken(HABITICA_API, userName, password)
})

var _ = AfterSuite(func() {
	DeleteUser(HABITICA_API, userName, password, "go-habits integration test")
})

// GoHabitsWithStdin builds session and stdin writer for invoking
// commands go-habits. go-habits must be install in PATH
func GoHabitsWithStdin(args ...string) (*gexec.Session, io.WriteCloser) {
	command := exec.Command(GOHABITS, args...)
	stdin, err := command.StdinPipe()
	Ω(err).ShouldNot(HaveOccurred())
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	return session, stdin
}

// GoHabits builds session for invoking commands go-habits.
// The go-habits binary must be install in PATH.
func GoHabits(args ...string) *gexec.Session {
	command := exec.Command(GOHABITS, args...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	return session
}

// RegisterUser uses habitica api to register a new user
func RegisterUser(serverUri, username, password, email string) {
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

	err := apiClient.Post("/user/auth/local/register", &payload, nil)
	Ω(err).ShouldNot(HaveOccurred())
}

// SaveAPIToken saves habitica api token and id to integration package
// variables
func SaveAPIToken(serverUri, username, password string) {
	creds := apiClient.Authenticate(username, password)
	apiToken = creds.APIToken
	apiID = creds.ID
}

// DeleteUser removes a user from habitica using api
func DeleteUser(serverUri, username, password, feedback string) {
	payload := `{"password":"` + password + `","feedback":"` + feedback + `"}`
	req, err := http.NewRequest("DELETE", serverUri+"/v3/user", bytes.NewBuffer([]byte(payload)))
	Ω(err).ShouldNot(HaveOccurred())
	err = apiClient.Do(req, nil)
	Ω(err).ShouldNot(HaveOccurred())
}
