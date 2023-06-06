package obp

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"gopkg.in/yaml.v3"
)

type PrettyDetailFormat int

const (
	JSON = iota
	Markdown
	CSV
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

func recurseDetails(detailed jsonschema.Detailed, acc []jsonschema.Detailed) []jsonschema.Detailed {
	acc = append(acc, detailed)
	for _, e := range detailed.Errors {
		acc = append(acc, recurseDetails(e, acc)...)
	}

	return acc
}

func PrettyDetails(w io.Writer, format PrettyDetailFormat, details jsonschema.Detailed) error {
	acc := make([]jsonschema.Detailed, 0)
	errors := recurseDetails(details, acc)
	switch format {
	case JSON:
		enc := json.NewEncoder(w)
		err := enc.Encode(errors)
		if err != nil {
			return fmt.Errorf("error writing JSON payload to provided writer: %w", err)
		}
	default:
		return fmt.Errorf("unknown format")

	}

	return nil
}
