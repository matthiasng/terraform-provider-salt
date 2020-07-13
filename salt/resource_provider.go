package salt

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"SALT_ADDRESS",
				}, "localhost:8500"),
			},

			"eauth": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"SALT_EAUTH",
				}, ""),
			},

			"username": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"SALT_USERNAME",
				}, ""),
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"SALT_PASSWORD",
				}, ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"salt_grains": dataSourceSaltGrains(),
		},

		ResourcesMap: map[string]*schema.Resource{},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &client{
		address:  d.Get("address").(string),
		eauth:    d.Get("eauth").(string),
		username: d.Get("username").(string),
		password: d.Get("password").(string),
	}

	return client, nil
}

func getClient(meta interface{}) *client {
	return meta.(*client)
}
