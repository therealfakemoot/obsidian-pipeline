/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"code.ndumas.com/ndumas/obsidian-pipeline"
)

var hugoBundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "convert a set of Obsidian notes into a Hugo compatible directory structure",
	Long:  `generate hugo content from your vault`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// here is where I validate arguments, open and parse config files, etc
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		source := viper.GetString("hugo.source")
		target := viper.GetString("hugo.target")

		err := obp.CopyPosts(source, target)
		if err != nil {
			return fmt.Errorf("error copying posts in %q: %w", source, err)
		}

		err = obp.Sanitize(source)
		if err != nil {
			return fmt.Errorf("error sanitizing posts in %q: %w", source, err)
		}

		err = obp.GatherMedia(source)
		if err != nil {
			return fmt.Errorf("error gathering media in %q: %w", source, err)
		}
		return nil
	},
}

func init() {
	hugoBundleCmd.Flags().StringP("source", "s", "", "path to vault directory containing hugo posts")
	err := viper.BindPFlag("hugo.source", hugoBundleCmd.Flags().Lookup("source"))
	if err != nil {
		log.Panicln("error binding viper to source flag:", err)
	}

	hugoBundleCmd.Flags().StringP("target", "t", "", "hugo content/ directory")
	err = viper.BindPFlag("hugo.target", hugoBundleCmd.Flags().Lookup("target"))
	if err != nil {
		log.Panicln("error binding viper to target flag:", err)
	}

	hugoBundleCmd.MarkFlagsRequiredTogether("source", "target")

	hugoCmd.AddCommand(hugoBundleCmd)
}
