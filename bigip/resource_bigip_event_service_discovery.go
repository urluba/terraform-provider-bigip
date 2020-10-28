/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"regexp"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceServiceDiscovery() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceDiscoveryCreate,
		Read:   resourceServiceDiscoveryRead,
		Update: resourceServiceDiscoveryUpdate,
		Delete: resourceServiceDiscoveryDelete,
		Exists: resourceServiceDiscoveryExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"partition": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the partition/tenant",
				ForceNew:     true,
			},

			"application": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Application",
				ForceNew:    true,
			},
			"pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "pool name",
				ForceNew:    true,
			},
			"node_list": {
				Type:     schema.TypeList,
				Optional: true,
				//MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "name of node",
						},
						"ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ip of nonde",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "port",
						},
					    },
					},
				},
			},
}

func resourceServiceDiscoveryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	return nil
}

func resourceServiceDiscoveryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()


	return nil
}

func resourceServiceDiscoveryExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
}

func resourceServiceDiscoveryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
}

func resourceServiceDiscoveryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	return nil
}
