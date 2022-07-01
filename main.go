package main

import (
	"context"
	"flag"
	"log"


	"github.com/Flagsmith/terraform-provider-flagsmith/flagsmith"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		Address: "registry.terraform.io/hashicorp/scaffolding",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), flagsmith.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
