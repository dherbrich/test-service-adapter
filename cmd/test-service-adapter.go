package main

import (
	"log"
	"os"

	"github.com/dherbric/test-service-adapter/adapter"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)


func main() {
	logger := log.New(os.Stderr, "test-service-adapter", log.LstdFlags)


	manifestGenerator := adapter.TestServiceManifestGenerator{
		Logger: logger,
	}

	handler := serviceadapter.CommandLineHandler{
		ManifestGenerator:     manifestGenerator,
	}
	serviceadapter.HandleCLI(os.Args, handler)
}