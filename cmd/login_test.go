package cmd

import (
	"github.com/8Mobius8/go-habits/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

type mockAuthenticateServer struct {
	AuthenticateFunc func(user, password string) (api.UserToken, error)
}

func (m mockAuthenticateServer) Authenticate(user, password string) (api.UserToken, error) {
	return m.AuthenticateFunc(user, password)
}

var _ = Describe("Login cmd", func() {
	Describe("Login", func() {
		var in, out *gbytes.Buffer
		BeforeEach(func() {
			in = gbytes.NewBuffer()
			out = gbytes.NewBuffer()
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
			mockServer := mockAuthenticateServer{
				func(user, password string) (api.UserToken, error) {
					return expectedCreds, nil
				},
			}

			Login(in, out, mockServer, []string{})

			Eventually(out).Should(gbytes.Say("Username:"))
			in.Write([]byte(expectedCreds.UserName + "\n"))
			Eventually(out).Should(gbytes.Say("Password:"))
			in.Write([]byte(expectedCreds.Password + "\n"))
		})
	})
})
