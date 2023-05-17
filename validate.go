package obp

import (
	"fmt"
	"log"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"gopkg.in/yaml.v3"
)

func Validate(schemaURL, filename string) error {
	var m interface{}

	target, err := os.Open(filename)
	if err != nil {
		log.Fatalf("could not open target file: %s\n", err)
	}
	dec := yaml.NewDecoder(target)
	err = dec.Decode(&m)
	if err != nil {
		log.Fatalf("error decoding YAML: %s\n", err)
	}

	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		log.Fatalf("error compiling schema: %s\n", err)
	}
	if err := schema.Validate(m); err != nil {
		log.Fatalf("error validating: %#v\n", err)
	}
	fmt.Println("validation successfull")

	return nil
}
