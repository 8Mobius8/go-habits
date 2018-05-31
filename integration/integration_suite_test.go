package integration_test

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")

}

var _ = Describe("go-habits", func() {
	var goHabitsCLIPath string
	var command *exec.Cmd

	BeforeSuite(func() {
		var err error
		goHabitsCLIPath, err = gexec.Build("github.com/8mobius8/go-habits")
		Ω(err).ShouldNot(HaveOccurred())
	})

	Describe("isup command", func() {
		It("exits with a zero", func() {
			command = exec.Command(goHabitsCLIPath, "isup")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			Eventually(session).Should(gbytes.Say(`Habitica is reachable, GO catch all those pets\!`))
			Eventually(session).Should(gexec.Exit(0))
		})
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})
})
