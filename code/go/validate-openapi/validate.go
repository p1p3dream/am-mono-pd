package main

import (
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

const (
	ExitSuccess = 0
	ExitFailure = 1
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("::error::Usage: validate-openapi <spec.yaml>")
		os.Exit(ExitFailure)
	}

	filename := os.Args[1]
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(filename)
	if err != nil {
		fmt.Printf("::error::Failed to load spec: %v\n", err)
		os.Exit(ExitFailure)
	}

	if err = doc.Validate(loader.Context); err != nil {
		fmt.Printf("::error::Spec is invalid: %v\n", err)
		os.Exit(ExitFailure)
	}

	fmt.Println("::notice::OpenAPI specification is valid!")
	os.Exit(ExitSuccess)
} 