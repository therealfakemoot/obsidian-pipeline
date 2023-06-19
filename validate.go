package obp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/santhosh-tekuri/jsonschema/v5"
	// allow the jsonschema validator to auto-download http-hosted schemas.
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"gopkg.in/yaml.v3"
)

var ErrUnsupportedOutputFormat = errors.New("unspported output format")

// Validate accepts a Markdown file as input via the Reader
// and parses the frontmatter present, if any. It then
// applies the schema fetched from schemaURL against the
// decoded YAML.
func Validate(schemaURL string, r io.Reader) error {
	var frontmatter interface{}

	dec := yaml.NewDecoder(r)

	err := dec.Decode(&frontmatter)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}

	compiler := jsonschema.NewCompiler()

	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		return fmt.Errorf("error compiling schema: %w", err)
	}

	if err != nil {
		return fmt.Errorf("frontmatter failed validation: %w", schema.Validate(frontmatter))
	}

	return nil
}

func recurseDetails(detailed jsonschema.Detailed, acc map[string]jsonschema.Detailed) map[string]jsonschema.Detailed {
	if detailed.Error != "" {
		acc[detailed.AbsoluteKeywordLocation] = detailed
	}

	for _, e := range detailed.Errors {
		acc = recurseDetails(e, acc)
	}

	return acc
}

// PrettyDetails takes error output from jsonschema.Validate
// and pretty-prints it to stdout.
//
// Supported formats are: JSON, Markdown.
func PrettyDetails(writer io.Writer, format string, details jsonschema.Detailed, filename string) error {
	// acc := make([]jsonschema.Detailed, 0)
	acc := make(map[string]jsonschema.Detailed)
	errors := recurseDetails(details, acc)

	switch format {
	case "json":
		enc := json.NewEncoder(writer)

		err := enc.Encode(details)
		if err != nil {
			return fmt.Errorf("error writing JSON payload to provided writer: %w", err)
		}
	case "markdown":
		fmt.Fprintf(writer, "# Validation Errors for %q\n", filename)
		fmt.Fprintf(writer, "Validation Rule|Failing Property|Error\n")
		fmt.Fprintf(writer, "--|---|---\n")

		for _, e := range errors {
			fmt.Fprintf(writer, "%s|%s|%s\n", e.KeywordLocation, e.InstanceLocation, e.Error)
		}
	default:
		return ErrUnsupportedOutputFormat
	}

	return nil
}
