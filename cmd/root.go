package cmd

import (
	"os"

	log "github.com/amoghe/distillog"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version of CLI, is set by go flags
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-habits",
	Short: "A Habitica CLI-interface written in go.",
	Long:  `A Habitica command line interface written in golang.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setupAPIClient()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = Version
	if err := rootCmd.Execute(); err != nil {
		exitCode := handleRootError(err)
		os.Exit(exitCode)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.SetOutput(os.Stdout)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.go-habits.yml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log", "", "logging level (default is none; can use ERROR,INFO,WARN,DEBUG)")

	rootCmd.PersistentFlags().StringVar(&habitsServerURL, "server", "http://habitica.com/api", "Set to '/api' uri of desired habits server.")
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
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
			log.Errorln(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-habits" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-habits")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Errorln("Error reading config file:", viper.ConfigFileUsed())
		log.Errorln(err)
	}
}
