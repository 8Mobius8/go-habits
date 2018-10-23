package commands

import (
	"os"
	"testing"

	"github.com/8Mobius8/go-habits/api"
	. "github.com/8Mobius8/go-habits/integration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var _ = BeforeSuite(func() {
	var exists bool
	HABITICA_API, exists = os.LookupEnv("SERVER")
	Ω(exists).ShouldNot(BeFalse())
	Ω(HABITICA_API).ShouldNot(BeEmpty())

	BUILD_VERSION, exists = os.LookupEnv("BUILD_VERSION")
	Ω(exists).ShouldNot(BeFalse())
	Ω(BUILD_VERSION).ShouldNot(BeEmpty())

	ApiClient = api.NewHabiticaAPI(nil, HABITICA_API)
	RegisterUser(HABITICA_API, UserName, Password, Email)
	SaveAPIToken(HABITICA_API, UserName, Password)
})

var _ = AfterSuite(func() {
	DeleteUser(HABITICA_API, UserName, Password, "go-habits integration test")
})
