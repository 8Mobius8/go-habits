package cmd

import (
	api "github.com/8Mobius8/go-habits/api"
	homedir "github.com/mitchellh/go-homedir"
)

var cfgFile string
var habitsServer *api.HabiticaAPI
var habitsServerURL string

var DefaultGoHabitsConfigPath string

func init() {
	userHomePath, _ := homedir.Dir()
	DefaultGoHabitsConfigPath = userHomePath + "/.go-habits.yml"
}
