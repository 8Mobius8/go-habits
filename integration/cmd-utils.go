package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// Global variables for all integration tests
var (
	HabiticaAPIURI string
	BuildVersion   string
	APIClient      *api.HabiticaAPI
	APIToken       string
	APIID          string
	command        *exec.Cmd
)

// GoHabitsWithStdin builds session and stdin writer for invoking
// commands go-habits. go-habits must be install in PATH
func GoHabitsWithStdin(args ...string) (*gexec.Session, io.WriteCloser) {
	args = prependConfigArg(args...)
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
	args = prependConfigArg(args...)
	command := exec.Command(GOHABITS, args...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Ω(err).ShouldNot(HaveOccurred())
	return session
}

// RegisterUser uses habitica api to register a new user
func RegisterUser(serverURI, username, password, email string) {
	registerUserBody := struct {
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
	data, merr := json.Marshal(registerUserBody)
	if merr != nil {
		Ω(merr).ShouldNot(HaveOccurred())
	}

	res, err := http.Post(serverURI+"/v3/user/auth/local/register", "application/json", bytes.NewBuffer(data))
	Ω(err).ShouldNot(HaveOccurred())
	Ω(res.StatusCode).ShouldNot(BeNumerically(">=", 300))
}

// GetAPIToken returns habitica api token and id
func GetAPIToken(serverURI, username, password string) (token, id string) {
	creds := APIClient.Authenticate(username, password)
	return creds.APIToken, creds.ID
}

// SaveAPIToken saves habitica api token and id to integration package
// variables
func SaveAPIToken(serverURI, username, password string) {
	token, id := GetAPIToken(serverURI, username, password)
	APIToken = token
	APIID = id
}

// DeleteUser removes a user from habitica using api
func DeleteUser(serverURI, username, password, feedback string) {
	payload := `{"password":"` + password + `","feedback":"` + feedback + `"}`
	req, err := http.NewRequest("DELETE", serverURI+"/v3/user", bytes.NewBuffer([]byte(payload)))
	Ω(err).ShouldNot(HaveOccurred())
	err = APIClient.Do(req, nil)
	Ω(err).ShouldNot(HaveOccurred())
}

// ResetUser leaves auth and api token but removes all data from their account.
func ResetUser() {
	payload := ""
	req, err := http.NewRequest("POST", HabiticaAPIURI+"/v3/user/reset", bytes.NewBuffer([]byte(payload)))
	Ω(err).ShouldNot(HaveOccurred())
	err = APIClient.Do(req, nil)
	Ω(err).ShouldNot(HaveOccurred())
}

// EventuallyLogin takes an existing session, waits to be prompted username
// and pasword, and writes the username and password to stdin.
func EventuallyLogin(session *gexec.Session, in io.WriteCloser, username, password string) {
	Eventually(session).Should(gbytes.Say("Username:"))
	in.Write([]byte(username + "\n"))

	Eventually(session).Should(gbytes.Say("Password:"))
	in.Write([]byte(password + "\n"))
}

// TouchConfigFile gets the User's config file path and writes
// a new line to it.
func TouchConfigFile() string {
	path := GetUserConfigPath()

	err := ioutil.WriteFile(path, []byte("\n"), 0644)
	Expect(err).ShouldNot(HaveOccurred())

	return path
}

// RemoveConfigFile tries to remove the User's config file.
// It returns safely of the file doesn't exist.
func RemoveConfigFile() {
	path := GetUserConfigPath()

	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
	}
	Expect(err).ShouldNot(HaveOccurred())
}

// GetUserConfigPath returns the default path for a users config file.
func GetUserConfigPath() string {
	return ".go-habits.yml"
}

func prependConfigArg(args ...string) []string {
	configArgs := []string{"--config", GetUserConfigPath()}
	return append(configArgs, args...)
}
