package obp

import (
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"gopkg.in/yaml.v3"
)

func Validate(schemaURL, filename string) error {
	var m interface{}

	target, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open target file: %w", err)
	}
	dec := yaml.NewDecoder(target)
	err = dec.Decode(&m)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}

	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		return fmt.Errorf("error compiling schema: %w", err)
	}
	if err := schema.Validate(m); err != nil {
		return fmt.Errorf("error validating target %q: %w", filename, err)
	}

	return nil
}
