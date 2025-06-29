package main

import (
	"fmt"
	"log"
	"os"

	"github.com/act3-ai/gitoci/pkg/apis"

	"git.act3-ace.com/ace/go-common/pkg/genschema"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Must specify a target directory for schema generation.")
	}

	scheme := apis.NewScheme()

	// Generate JSON Schema definitions
	if err := genschema.GenerateGroupSchemas(
		os.Args[1],
		scheme,
		[]string{"gitoci.act3-ai.io"},
		"github.com/act3-ai/gitoci",
	); err != nil {
		log.Fatal(fmt.Errorf("JSON Schema generation failed: %w", err))
	}
}
