package integration

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("go-habits status", func() {
	Context("when Habitica api is reachabled", func() {
		It("display Habitica is reachable status", func() {
			session := GoHabits("status")

			Eventually(session).Should(gbytes.Say(`Habitica is reachable, GO catch all those pets!`))
			Eventually(session).Should(gexec.Exit(0))
		})
		It("display Habitica is reachable status", func() {
			session := GoHabits("s")

			Eventually(session).Should(gbytes.Say(`Habitica is reachable, GO catch all those pets!`))
			Eventually(session).Should(gexec.Exit(0))
		})
	})
	Context("when Habitica api is NOT reachabled", func() {
		var serverUri string
		BeforeEach(func() {
			serverUri = os.Getenv("SERVER")
			os.Setenv("SERVER", "http://notaccessible/api")
		})
		AfterEach(func() {
			os.Setenv("SERVER", serverUri)
		})

		It("display Habitica is unreachable status", func() {
			session := GoHabits("status")

			Eventually(session).Should(gbytes.Say("Habitica is unreachable"))
			Eventually(session).Should(gexec.Exit(5))
		})
		It("display Habitica is unreachable status", func() {
			session := GoHabits("s")

			Eventually(session).Should(gbytes.Say("Habitica is unreachable"))
			Eventually(session).Should(gexec.Exit(5))
		})
	})
	Context("when --server option is used", func() {
		It("will change api server to given --server", func() {
			inlineServerURI := "http://inlineserverruri/api"
			session := GoHabits("status", "--server", inlineServerURI)

			Eventually(session).Should(gbytes.Say("Using " + inlineServerURI + " as api server"))
		})
	})
})
