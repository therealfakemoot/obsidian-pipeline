package obp

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"gopkg.in/yaml.v3"
)

func Validate(schemaURL string, r io.Reader) error {
	var m interface{}

	dec := yaml.NewDecoder(r)
	err := dec.Decode(&m)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		return fmt.Errorf("error compiling schema: %w", err)
	}
	if err := schema.Validate(m); err != nil {
		return err
	}

	return nil
}

func flattenDetails(detailed jsonschema.Detailed, acc []jsonschema.Detailed) []jsonschema.Detailed {
	if detailed.Error != "" {
		acc = append(acc, detailed)
	}
	for _, e := range detailed.Errors {
		acc = append(acc, flattenDetails(e, acc)...)
	}

	return acc
}

func PrettyDetails(w io.Writer, format string, details jsonschema.Detailed, filename string) error {
	acc := make([]jsonschema.Detailed, 0)
	errors := flattenDetails(details, acc)
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		err := enc.Encode(details)
		if err != nil {
			return fmt.Errorf("error writing JSON payload to provided writer: %w", err)
		}
	case "markdown":
		fmt.Fprintf(w, "# Validation Errors for %q\n", filename)
		fmt.Fprintf(w, "eyword Location|Instance Location|Error\n")
		fmt.Fprintf(w, "--|---|---\n")
		for _, e := range errors {
			fmt.Fprintf(w, "%s|%s|%s\n", e.KeywordLocation, e.InstanceLocation, e.Error)
		}
	default:
		return fmt.Errorf("unknown format")

	}

	return nil
}
