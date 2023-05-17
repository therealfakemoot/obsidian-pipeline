package obp

import (
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
		return fmt.Errorf("error validating target: %w", err)
	}

	return nil
}
