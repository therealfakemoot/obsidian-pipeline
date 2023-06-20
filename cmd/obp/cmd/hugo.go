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

var hugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "convert a set of Obsidian notes into a Hugo compatible directory structure",
	Long:  `generate hugo content from your vault`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// here is where I validate arguments, open and parse config files, etc
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		target := viper.GetString("target")
		source := viper.GetString("source")
		log.Printf("provided source: %q\n", source)
		log.Printf("provided target: %q\n", target)
		return nil

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
	hugoCmd.Flags().StringP("source", "s", "", "path to vault directory containing hugo posts")

	err := viper.BindPFlag("source", hugoCmd.Flags().Lookup("source"))
	if err != nil {
		log.Panicln("error binding viper to source flag:", err)
	}

	hugoCmd.Flags().StringP("target", "t", "", "hugo content/ directory")
	err = viper.BindPFlag("target", hugoCmd.Flags().Lookup("target"))
	if err != nil {
		log.Panicln("error binding viper to target flag:", err)
	}

	hugoCmd.MarkFlagsRequiredTogether("source", "target")

	rootCmd.AddCommand(hugoCmd)
}
