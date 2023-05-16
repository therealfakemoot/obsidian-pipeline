/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "fmt"

	"github.com/spf13/cobra"
)

var (
	source, target string
)

// rootCmd represents the base command when called without any subcommands
var hugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "convert a set of Obsidian notes into a Hugo compatible directory structure",
	Long:  `long description`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// here is where I validate arguments, open and parse config files, etc
		return nil
	},
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	hugoCmd.PersistentFlags().StringVar(&source, "source", "", "directory containing ready-to-publish posts")
	hugoCmd.PersistentFlags().StringVar(&target, "target", "", "target Hugo directory (typically content/posts)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// rootCmd.SetHelpFunc(gloss.CharmHelp)
	// rootCmd.SetUsageFunc(gloss.CharmUsage)
	rootCmd.AddCommand(hugoCmd)
}
