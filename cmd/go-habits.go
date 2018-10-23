package cmd

import (
	api "github.com/8Mobius8/go-habits/api"
	homedir "github.com/mitchellh/go-homedir"
)

var cfgFile string
var habitsServer *api.HabiticaAPI
var habitsServerURL string

var defaultGoHabitsConfigPath string

func init() {
	userHomePath, _ := homedir.Dir()
	defaultGoHabitsConfigPath = userHomePath + "/.go-habits.yml"
}
