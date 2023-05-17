/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"github.com/spf13/cobra"

	"code.ndumas.com/ndumas/obsidian-pipeline"
)

// rootCmd represents the base command when called without any subcommands
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "loads a note and ensures its frontmatter follows the provided protobuf schema",
	Long: `Validate YAML frontmatter with jsonschema
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PreRunE: func(cmd *cobra.Command, args []string) error {
		schemaURL := cmd.Flag("schema").Value.String()
		if len(schemaURL) == 0 {
			return fmt.Errorf("Please profide a schema filename")
		}
		target := cmd.Flag("target").Value.String()
		if len(target) == 0 {
			return fmt.Errorf("Please profide a target filename")
		}
		root := os.DirFS(target)
		_, err := root.Open(".")
		if err != nil {
			return fmt.Errorf("cannot open target directory %q: %w", target, err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		schema := cmd.Flag("schema").Value.String()
		target := cmd.Flag("target").Value.String()
		root := os.DirFS(target)
		return fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			log.Printf("scanning %q\n", path)
			return obp.Validate(schema, target)
		})
		// return obp.Validate(schema, target)
	},
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	validateCmd.Flags().StringP("schema", "s", "base.schema", "path to protobuf file")
	validateCmd.Flags().StringP("target", "t", "", "directory containing validation targets")
	rootCmd.AddCommand(validateCmd)
}
