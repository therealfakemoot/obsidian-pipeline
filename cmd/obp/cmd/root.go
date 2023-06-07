/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// "code.ndumas.com/ndumas/obsidian-pipeline/gloss"
)

var (
	vault   string
	cfgFile string
	format  string
)

var rootCmd = &cobra.Command{
	TraverseChildren: true,
	Use:              "obp",
	Short:            "obp is a toolkit for managing your vault in headless contexts",
	Long:             `a suite of tools for managing your obsidian vault`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// here is where I validate arguments, open and parse config files, etc
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "~/.obp.toml", "config file")

	rootCmd.PersistentFlags().StringVar(&vault, "vault", "", "vault root directory")

	err := rootCmd.MarkPersistentFlagRequired("vault")

	if err != nil {
		log.Panicln("error setting vault flag as required")
	}

	rootCmd.PersistentFlags().StringVar(&format, "format", "markdown", "output format [markdown, json, csv]")

	err = rootCmd.MarkPersistentFlagRequired("format")
	if err != nil {
		log.Panicln("error setting format flag as required")
	}

	err = viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
	if err != nil {
		log.Panicln("error binding viper to format flag")
	}

	err = viper.BindPFlag("vault", rootCmd.PersistentFlags().Lookup("vault"))
	if err != nil {
		log.Panicln("error binding viper to vault flag")
	}

	// rootCmd.SetHelpFunc(gloss.CharmHelp)
	// rootCmd.SetUsageFunc(gloss.CharmUsage)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cmd" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".obp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
