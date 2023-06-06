/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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

// rootCmd represents the base command when called without any subcommands
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
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "~/.obp.toml", "config file")
	rootCmd.PersistentFlags().StringVar(&vault, "vault", "", "vault root directory")
	rootCmd.MarkPersistentFlagRequired("vault")
	rootCmd.PersistentFlags().StringVar(&format, "format", "markdown", "output format [markdown, json, csv]")
	rootCmd.MarkPersistentFlagRequired("format")

	viper.BindPFlag("format", validateCmd.Flags().Lookup("format"))
	viper.BindPFlag("vault", validateCmd.Flags().Lookup("vault"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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
