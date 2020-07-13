package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/matthiasng/terraform-provider-salt/salt"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: salt.Provider})
}
