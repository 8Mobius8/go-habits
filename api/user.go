package api

// Stats returns a Stats object with the currently authenticated
// Habitica user's basic stats.
func (api *HabiticaAPI) Stats() (Stats, error) {
	body, err := api.Get("/user/anonymized")

	var user userResponse
	if err == nil {
		api.ParseResponse(body, &user)
	}

	charData := user.Data.User.Stats
	var stats Stats
	stats.Class = charData.Class
	stats.Constitution = charData.Con
	stats.Experience = charData.Exp
	stats.ExperienceToNextLevel = charData.ToNextLevel
	stats.Gold = charData.Gp
	stats.Health = charData.Hp
	stats.Intelligence = charData.Int
	stats.Level = charData.Lvl
	stats.Mana = charData.Mp
	stats.MaxHealth = charData.MaxHealth
	stats.MaxMana = charData.MaxMP
	stats.Perception = charData.Per
	stats.Points = charData.Points
	stats.Strength = charData.Str

	return stats, nil
}

type userResponse struct {
	Data struct {
		User struct {
			Stats struct {
				Hp          int
				Mp          int
				Exp         int
				Gp          int
				Lvl         int
				Class       string
				Points      int
				Str         int
				Con         int
				Int         int
				Per         int
				ToNextLevel int
				MaxHealth   int
				MaxMP       int
			}
		}
	}
}

// Stats struct with User's current stats
type Stats struct {
	Level                 int
	Health                int
	Mana                  int
	Experience            int
	Gold                  int
	Class                 string
	Points                int
	Strength              int
	Constitution          int
	Intelligence          int
	Perception            int
	ExperienceToNextLevel int
	MaxHealth             int
	MaxMana               int
}
