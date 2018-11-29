package flows

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
	RunSpecs(t, "Flows Suite")
}

var _ = BeforeSuite(func() {
	var exists bool
	HabiticaAPIURI, exists = os.LookupEnv("SERVER")
	Ω(exists).ShouldNot(BeFalse())
	Ω(HabiticaAPIURI).ShouldNot(BeEmpty())

	BuildVersion, exists = os.LookupEnv("BUILD_VERSION")
	Ω(exists).ShouldNot(BeFalse())
	Ω(BuildVersion).ShouldNot(BeEmpty())

	APIClient = api.NewHabiticaAPI(nil, HabiticaAPIURI, log.NewNullLogger("test"))
})
