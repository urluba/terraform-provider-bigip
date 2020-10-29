---
layout: "bigip"
page_title: "BIG-IP: bigip_event_service_discovery"
sidebar_current: "docs-bigip-resource-event_service_discovery-x"
description: |-
   Provides details about bigip_event_service_discovery resource
---

# bigip_event_service_discovery

`bigip_event_service_discovery` Terraform resource to update pool members based on event driven Service Discovery.

The API endpoint for Service discovery tasks should be available before using the resource and with this resource,we will be able to connect to a specific endpoint related to event based service discovery that will allow us to update the list of pool members
 

## Example Usage


```hcl
resource "bigip_event_service_discovery" "test" {
  tenant_name      = "Sample_event_sd"
  application_name = "My_app"
  pool_name        = "My_pool"
  node_list {
    id   = "newNode1"
    ip   = "192.168.2.3"
    port = 8080
  }
  node_list {
    id   = "newNode2"
    ip   = "192.168.2.4"
    port = 8080
  }
}

```      

## Argument Reference

* `tenant_name` - (Required) Name of the Partition (tenant)

* `application_name` - (Required) Name of the application

* `pool_name` - (Required) Name of the pool where nodes reside

* `node_list` - (Required) Map of node which will be added to pool which will be having node name(id),node address(ip) and node port(port)

For more information, please refer below document
https://clouddocs.f5.com/products/extensions/f5-appsvcs-extension/latest/declarations/discovery.html?highlight=service%20discovery#event-driven-service-discovery
