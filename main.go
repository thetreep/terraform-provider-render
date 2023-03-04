package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/jackall3n/terraform-provider-render/render"
	"log"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	err := providerserver.Serve(context.Background(), render.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/jackall3n/render",
	})

	if err != nil {
		log.Fatalf("unable to serve provider: %s", err)
	}
}
