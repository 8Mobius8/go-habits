package cmd

import (
	"fmt"
	"os"

	cmds "github.com/8Mobius8/go-habits/cmd/commands"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// HabitsServerURL is the uri used to send API requests to
var HabitsServerURL string

// Version of CLI, is set by go flags
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-habits",
	Short: "A Habitica CLI-interface written in go.",
	Long:  `A Habitica command line interface written in golang.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = Version
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.go-habits.yml)")

	rootCmd.PersistentFlags().StringVarP(&HabitsServerURL, "server", "", "http://habitica.com/api", "Set to '/api' uri of desired habits server.")
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	addGoHabitCommands(rootCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-habits" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-habits")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func addGoHabitCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(versionCmd)
}

var addCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a todo to Habitica",
	Aliases: []string{"a", "a t"},
	Run:     cmds.Add,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List todos",
	Aliases: []string{"l", "l t"},
	Run:     cmds.List,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenicates with Habits server and saves api token in config file.",
	Long:  `Authenicates with Habits server and saves api token in config file.`,
	Run:   cmds.Login,
}

var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Check if Habitica api is reachable.",
	Aliases: []string{"s"},
	Run:     cmds.Status,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version of go-habits.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(rootCmd.Use + " version " + rootCmd.Version)
	},
}
