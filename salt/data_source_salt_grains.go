package salt

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceSaltGrains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSaltGrainsRead,
		Schema: map[string]*schema.Schema{
			// Filters
			"target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "glob",
			},

			// Out parameters
			"grains": {
				Computed: true,
				Type:     schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSaltGrainsRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)

	target := d.Get("target").(string)
	targetType := d.Get("target_type").(string)

	result, err := client.Run(target, targetType, "grains.items")
	if err != nil {
		return err
	}

	grains := map[string]string{}
	for name, entry := range result {
		s, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		grains[name] = string(s)
	}

	d.SetId(fmt.Sprintf("target-%s-%s", targetType, target))
	d.Set("grains", grains)

	return nil
}
