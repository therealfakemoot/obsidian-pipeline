/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "fmt"

	"github.com/spf13/cobra"
)

var hugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "convert a set of Obsidian notes into a Hugo compatible directory structure",
	Long:  `generate hugo content from your vault`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// here is where I validate arguments, open and parse config files, etc
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	// hugoCmd.Flags().StringVar(&source, "source", "", "directory containing ready-to-publish posts")
	// hugoCmd.Flags().StringVar(&target, "target", "", "target Hugo directory (typically content/posts)")

	rootCmd.AddCommand(hugoCmd)
}
