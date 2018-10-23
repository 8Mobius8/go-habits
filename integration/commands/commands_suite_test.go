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
	HabiticaAPIURI, exists = os.LookupEnv("SERVER")
	Ω(exists).ShouldNot(BeFalse())
	Ω(HabiticaAPIURI).ShouldNot(BeEmpty())

	BuildVersion, exists = os.LookupEnv("BUILD_VERSION")
	Ω(exists).ShouldNot(BeFalse())
	Ω(BuildVersion).ShouldNot(BeEmpty())

	APIClient = api.NewHabiticaAPI(nil, HabiticaAPIURI)
	RegisterUser(HabiticaAPIURI, UserName, Password, Email)
	SaveAPIToken(HabiticaAPIURI, UserName, Password)
})

var _ = AfterSuite(func() {
	DeleteUser(HabiticaAPIURI, UserName, Password, "go-habits integration test")
})
