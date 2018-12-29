package cmd

import (
	"github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type MockAuthenticateServer struct {
	AuthenticateOut api.UserToken
	ErrorOut        error
}

func (ms MockAuthenticateServer) Authenticate(user, password string) (api.UserToken, error) {
	return ms.AuthenticateOut, ms.ErrorOut
}

var _ = Describe("Login cmd", func() {
	Describe("Login", func() {
		var in, out *gbytes.Buffer
		var mockServer MockAuthenticateServer

		BeforeEach(func() {
			in = gbytes.NewBuffer()
			out = gbytes.NewBuffer()
			mockServer = MockAuthenticateServer{}
		})
		AfterEach(func() {
			in.Close()
			out.Close()
		})
		It("should print prompts for username and password", func() {
			expectedCreds := api.UserToken{
				UserName: "username",
				Password: "password",
			}
			mockServer.AuthenticateOut = expectedCreds
			Login(in, out, mockServer, []string{})

			Eventually(out).Should(gbytes.Say("Username:"))
			in.Write([]byte(expectedCreds.UserName + "\n"))
			Eventually(out).Should(gbytes.Say("Password:"))
			in.Write([]byte(expectedCreds.Password + "\n"))
		})
	})
})
