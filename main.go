package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/camptocamp/terraform-provider-geoserver/geoserver"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: geoserver.Provider,
	})
}
