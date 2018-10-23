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
	homedir "github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var (
	HABITICA_API  string
	BUILD_VERSION string
	ApiClient     *api.HabiticaAPI
	ApiToken      string
	ApiID         string
	command       *exec.Cmd
)

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

	res, err := http.Post(serverUri+"/v3/user/auth/local/register", "application/json", bytes.NewBuffer(data))
	Ω(err).ShouldNot(HaveOccurred())
	Ω(res.StatusCode).ShouldNot(BeNumerically(">=", 300))
}

// GetAPIToken returns habitica api token and id
func GetAPIToken(serverUri, username, password string) (token, id string) {
	creds := ApiClient.Authenticate(username, password)
	return creds.APIToken, creds.ID
}

// SaveAPIToken saves habitica api token and id to integration package
// variables
func SaveAPIToken(serverUri, username, password string) {
	token, id := GetAPIToken(serverUri, username, password)
	ApiToken = token
	ApiID = id
}

// DeleteUser removes a user from habitica using api
func DeleteUser(serverUri, username, password, feedback string) {
	payload := `{"password":"` + password + `","feedback":"` + feedback + `"}`
	req, err := http.NewRequest("DELETE", serverUri+"/v3/user", bytes.NewBuffer([]byte(payload)))
	Ω(err).ShouldNot(HaveOccurred())
	err = ApiClient.Do(req, nil)
	Ω(err).ShouldNot(HaveOccurred())
}

// ResetUser leaves auth and api token but removes all data from their account.
func ResetUser(userId, token string) {
	payload := ""
	req, err := http.NewRequest("POST", HABITICA_API+"/v3/user/reset", bytes.NewBuffer([]byte(payload)))
	Ω(err).ShouldNot(HaveOccurred())
	err = ApiClient.Do(req, nil)
	Ω(err).ShouldNot(HaveOccurred())
}

func EventuallyLogin(session *gexec.Session, in io.WriteCloser, username, password string) {
	Eventually(session).Should(gbytes.Say("Username:"))
	in.Write([]byte(username + "\n"))

	Eventually(session).Should(gbytes.Say("Password:"))
	in.Write([]byte(password + "\n"))
}

func TouchConfigFile() string {
	path := GetUserConfigPath()

	err := ioutil.WriteFile(path, []byte("test"), 0644)
	Expect(err).ShouldNot(HaveOccurred())

	return path
}

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

func GetUserConfigPath() string {
	userHomePath, err := homedir.Dir()
	Expect(err).ShouldNot(HaveOccurred())
	return userHomePath + "/.go-habits.yml"
}
