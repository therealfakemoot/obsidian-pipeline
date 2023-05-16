/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		schemaFilename := cmd.Flag("schema").Value.String()
		if len(schemaFilename) == 0 {
			return fmt.Errorf("Please profide a schema filename")
		}
		target := cmd.Flag("target").Value.String()
		if len(target) == 0 {
			return fmt.Errorf("Please profide a target filename")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var m interface{}
		var schemaBytes bytes.Buffer

		schemaFilename := cmd.Flag("schema").Value.String()
		schemaFile, err := os.Open(schemaFilename)
		if err != nil {
			log.Fatalf("could not open schema file: %s\n", err)
		}
		io.Copy(&schemaBytes, schemaFile)

		// err := yaml.Unmarshal([]byte(yamlText), &m)
		targetFilename := cmd.Flag("target").Value.String()
		target, err := os.Open(targetFilename)
		if err != nil {
			log.Fatalf("could not open target file: %s\n", err)
		}
		dec := yaml.NewDecoder(target)
		err = dec.Decode(&m)
		if err != nil {
			log.Fatalf("error decoding YAML: %s\n", err)
		}

		compiler := jsonschema.NewCompiler()
		if err := compiler.AddResource(schemaFilename, strings.NewReader(schemaBytes.String())); err != nil {
			log.Fatalf("error adding resource to jsonschema compiler: %s\n", err)
		}

		schema, err := compiler.Compile(schemaFilename)
		if err != nil {
			log.Fatalf("error compiling schema: %s\n", err)
		}
		if err := schema.Validate(m); err != nil {
			log.Fatalf("error validating: %s\n", err)
		}
		fmt.Println("validation successfull")
	},
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	validateCmd.Flags().StringP("schema", "s", "base.schema", "path to protobuf file")
	validateCmd.Flags().StringP("target", "t", "", "list of filenames to validate")
	rootCmd.AddCommand(validateCmd)
}
