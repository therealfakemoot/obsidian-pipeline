/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var hugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "manage your hugo blog using your vault as a source of truth",
	Long:  `manage your hugo blog using your vault as a source of truth`,
}

func init() {
	rootCmd.AddCommand(hugoCmd)
}
