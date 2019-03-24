package commands

import (
	"os"
	"testing"

	"github.com/8Mobius8/go-habits/api"
	log "github.com/amoghe/distillog"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/8Mobius8/go-habits/integration"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commands Suite")
}

var _ = BeforeSuite(func() {
	HabiticaAPIURI = EnvironmentVariableShouldNotBeEmpty("SERVER")
	BuildVersion = EnvironmentVariableShouldNotBeEmpty("BUILD_VERSION")

	APIClient = api.NewHabiticaAPI(nil, HabiticaAPIURI, log.NewNullLogger("test"))
	RegisterUser(HabiticaAPIURI, UserName, Password, Email)
	SaveAPIToken(HabiticaAPIURI, UserName, Password)
})

var _ = AfterSuite(func() {
	DeleteUser(HabiticaAPIURI, UserName, Password, "go-habits integration test")
})

func EnvironmentVariableShouldNotBeEmpty(envVariable string) string {
	contents, exists := os.LookupEnv(envVariable)
	Ω(exists).ShouldNot(BeFalse(), "%s environment variable should exist", envVariable)
	Ω(contents).ShouldNot(BeEmpty(), "%s environment variable should not be empty", envVariable)
	return contents
}
