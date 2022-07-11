package main

import (
	"context"
	"flag"
	"log"


	"github.com/Flagsmith/terraform-provider-flagsmith/flagsmith"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs


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
		Address: "registry.terraform.io/flagsmith/flagsmith",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), flagsmith.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
