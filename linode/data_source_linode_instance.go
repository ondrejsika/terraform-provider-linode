package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/linode/linodego"
)

func dataSourceLinodeInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeInstanceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Description: "This Linode's Private IPv4 Address.  The regional private IP address range is 192.168.128/17 address shared by all Linode Instances in a region.",
				Computed:    true,
			},
		},
	}
}

func dataSourceLinodeInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	reqInstance, _ := strconv.Atoi(d.Get("id").(string))

	// if reqInstance == 0 {
	// 	return fmt.Errorf("Instance id is required")
	// }

	Instance, err := client.GetInstance(context.Background(), reqInstance)
	if err != nil {
		return fmt.Errorf("Error listing Instances: %s", err)
	}

	instanceNetwork, err := client.GetInstanceIPAddresses(context.Background(), int(reqInstance))

	if err != nil {
		return fmt.Errorf("Error getting the IPs for Linode instance %s: %s", d.Id(), err)
	}
	private := instanceNetwork.IPv4.Private

	if len(private) > 0 {
		// d.Set("private_ip", true)
		d.Set("private_ip_address", private[0].Address)
	} else {
		// d.Set("private_ip", false)
	}

	if Instance != nil {
		d.SetId(d.Get("id").(string))
		// d.Set("private_ip_address", Instance.IPv6)
		// d.Set("private_ip_address", "1.2.3.4")
		return nil
	}

	d.SetId("")

	return fmt.Errorf("Instance %d was not found", reqInstance)
}
