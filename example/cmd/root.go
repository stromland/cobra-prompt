package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cobraprompt "github.com/stromland/cobra-prompt"
)

const CallbackAnnotation = cobraprompt.CallbackAnnotation

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "cobra-prompt-example",
	Short: "CoPrompt example",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		//os.Exit(1) //Commented out, breaks the prompt loop
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra-prompt-example.yaml)")
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".cobra-prompt-example" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra-prompt-example")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
