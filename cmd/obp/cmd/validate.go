/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"code.ndumas.com/ndumas/obsidian-pipeline"
)

// rootCmd represents the base command when called without any subcommands
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "loads a note and ensures its frontmatter follows the provided protobuf schema",
	Long: `Validate YAML frontmatter with jsonschema
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		schema := viper.GetString("schema")
		target := viper.GetString("target")
		if target == "" {
			return fmt.Errorf("target flag must not be empty")
		}
		root := os.DirFS(target)

		err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("could not walk %q: %w", path, err)
			}

			if d.IsDir() {
				return nil
			}

			absPath, err := filepath.Abs(filepath.Join(target, path))
			if err != nil {
				return fmt.Errorf("error generating absolute path for %q", target)
			}
			file, err := os.Open(absPath)
			if err != nil {
				return fmt.Errorf("could not open target file: %w", err)
			}
			defer file.Close()
			err = obp.Validate(schema, file)
			if err != nil {
				details, ok := err.(*jsonschema.ValidationError)
				if !ok {
					return fmt.Errorf("eror validating %q: %w", path, err)
				}
				obp.PrettyDetails(cmd.OutOrStdout(), viper.GetString("format"), details.DetailedOutput(), absPath)
			}
			return nil
		})

		return fmt.Errorf("validate command failed: %w", err)
	},
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	validateCmd.Flags().StringP("schema", "s", "base.schema", "path to protobuf file")
	validateCmd.Flags().StringP("target", "t", "", "directory containing validation targets")
	validateCmd.MarkFlagsRequiredTogether("schema", "target")
	validateCmd.PersistentFlags().StringVar(&format, "format", "markdown", "output format [markdown, json, csv]")
	rootCmd.AddCommand(validateCmd)

	err := viper.BindPFlag("schema", validateCmd.Flags().Lookup("schema"))
	if err != nil {
		log.Panicln("error binding viper to schema flag:", err)
	}

	err = viper.BindPFlag("target", validateCmd.Flags().Lookup("target"))
	if err != nil {
		log.Panicln("error binding viper to target flag:", err)
	}

	err = viper.BindPFlag("format", validateCmd.PersistentFlags().Lookup("format"))
	if err != nil {
		log.Panicln("error binding viper to format flag:", err)
	}
}
