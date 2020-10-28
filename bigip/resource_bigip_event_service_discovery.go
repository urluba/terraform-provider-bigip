/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceServiceDiscovery() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceDiscoveryCreate,
		Read:   resourceServiceDiscoveryRead,
		Update: resourceServiceDiscoveryUpdate,
		Delete: resourceServiceDiscoveryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the partition/tenant",
				ForceNew:    true,
			},
			"application_name": {
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
				Type:     schema.TypeSet,
				Optional: true,
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
}

func resourceServiceDiscoveryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO]: Client: %+v", client)
	tenantUri := fmt.Sprintf("~%s~%s~%s", d.Get("tenant_name").(string), d.Get("application_name").(string), d.Get("pool_name").(string))
	log.Printf("[INFO]: tenantUri: %+v", tenantUri)
	var nodeList []interface{}
	if m, ok := d.GetOk("node_list"); ok {
		for _, node := range m.(*schema.Set).List() {
			log.Printf("[INFO]: node Value: %+v", node)
			nodeList = append(nodeList, node)
		}
	}
	log.Printf("[INFO]: node Value: %+v", nodeList)
	err := client.AddServiceDiscoveryNodes(tenantUri, nodeList)
	if err != nil {
		return fmt.Errorf("Error modifying node %s: %v ", nodeList, err)
	}
	d.SetId(tenantUri)
	return resourceServiceDiscoveryRead(d, meta)
}

func resourceServiceDiscoveryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	tenantUri := d.Id()
	var nodeList []interface{}
	log.Printf("[INFO] Get Event driven service discovery nodes for Application:%+v", tenantUri)
	if m, ok := d.GetOk("node_list"); ok {
		for _, node := range m.(*schema.Set).List() {
			nodeList = append(nodeList, node)
		}
	}
	serviceDiscoveryResp, err := client.GetServiceDiscoveryNodes(tenantUri)
	if err != nil {
		return fmt.Errorf("Error Reading node : %v ", err)
	}
	//log.Printf("[INFO]: providerOptions: %+v", as3Resp.(map[string]interface{})["result"].(map[string]interface{})["providerOptions"].(map[string]interface{})["nodeList"])
	nodeList1 := serviceDiscoveryResp.(map[string]interface{})["result"].(map[string]interface{})["providerOptions"].(map[string]interface{})["nodeList"]
	nodeListCount := d.Get("node_list.#").(int)
	if len(nodeList1.([]interface{})) != nodeListCount {
		d.SetId("")
		return fmt.Errorf("[DEBUG] Get Node list failed for  (%s): %s", d.Id(), err)
	}
	if serviceDiscoveryResp == nil {
		d.SetId("")
		return fmt.Errorf("[DEBUG]serviceDiscoveryResp is : %s", serviceDiscoveryResp)
	}
	return nil
}

func resourceServiceDiscoveryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	tenantUri := d.Id()
	log.Printf("[INFO]: tenantUri: %+v", tenantUri)
	var nodeList []interface{}
	if m, ok := d.GetOk("node_list"); ok {
		for _, node := range m.(*schema.Set).List() {
			log.Printf("[INFO]: node Value: %+v", node)
			nodeList = append(nodeList, node)
		}
	}
	err := client.AddServiceDiscoveryNodes(tenantUri, nodeList)
	if err != nil {
		return fmt.Errorf("Error modifying node %s: %v ", nodeList, err)
	}
	return resourceServiceDiscoveryRead(d, meta)
}

func resourceServiceDiscoveryDelete(d *schema.ResourceData, meta interface{}) error {
	clientBigip := meta.(*bigip.BigIP)
	tenantUri := d.Id()
	url := clientBigip.Host + "/mgmt/shared/service-discovery/task/" + tenantUri + "/nodes/"
	payload := strings.NewReader("[ ]\n")
	log.Printf("[DEBUG] tenantUri Complete :%v", url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("Error while creating http request with AS3 json:%+v ", err)
	}
	req.SetBasicAuth(clientBigip.User, clientBigip.Password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	if resp.Status != "200 OK" || err != nil {
		defer resp.Body.Close()
		return fmt.Errorf("Error while Sending/Posting http request with AS3 json :%s  %v", bodyString, err)
	}
	defer resp.Body.Close()
	d.SetId("")
	return nil
}
