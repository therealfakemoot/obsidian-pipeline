/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version, Tag string
)

// rootCmd represents the base command when called without any subcommands
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints obp version info",
	Long: `displays git tag and sha
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", Version, Tag)
	},
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(versionCmd)

}
