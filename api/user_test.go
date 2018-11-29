package api_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("API User", func() {
	Describe("when getting a user's information", func() {
		It("returns character statistics for a user.", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v3/user/anonymized"),
					ghttp.RespondWith(http.StatusOK,
						`{"data":{"user":{"stats":{"hp": 50,"mp": 30,"exp": 12,"gp": 9001,"lvl": 1,"class": "warrior","points": 5,"str": 1,"con": 1,"int": 1,"per": 1,"toNextLevel": 150,"maxHealth": 50,"maxMP": 30}}}}`),
				),
			)

			stats, _ := habitapi.Stats()
			Expect(stats.Level).To(Equal(1))
			Expect(stats.Health).To(Equal(50))
			Expect(stats.Mana).To(Equal(30))
			Expect(stats.Experience).To(Equal(12))
			Expect(stats.Gold).To(Equal(9001))
			Expect(stats.Class).To(Equal("warrior"))
			Expect(stats.Points).To(Equal(5))
			Expect(stats.Strength).To(Equal(1))
			Expect(stats.Constitution).To(Equal(1))
			Expect(stats.Intelligence).To(Equal(1))
			Expect(stats.Perception).To(Equal(1))
			Expect(stats.ExperienceToNextLevel).To(Equal(150))
			Expect(stats.MaxHealth).To(Equal(50))
			Expect(stats.MaxMana).To(Equal(30))
		})
	})
})
