// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceProfileResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDeviceProfileResourceConfig("deviceprofile-one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("chirpstack_device_profile.test", "id"),
					resource.TestCheckResourceAttr("chirpstack_device_profile.test", "name", "deviceprofile-one"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "chirpstack_device_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDeviceProfileResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("chirpstack_device_profile.test", "name", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDeviceProfileResourceConfig(deviceprofileName string) string {
	return fmt.Sprintf(`
resource "chirpstack_tenant" "test" {
  name = "test_tenant"
}
resource "chirpstack_device_profile" "test" {
  tenant_id                       = chirpstack_tenant.test.id
  name                            = %[1]q
  description                     = "test"
  region                          = "AU915"
  region_parameters_revision      = "A"
  mac_version                     = "LORAWAN_1_0_3"
  flush_queue_on_activate         = true
  allow_roaming                   = false
  expected_uplink_interval        = 3600
  device_status_request_frequency = 1
  device_supports_otaa            = true
  device_supports_class_b         = false
  device_supports_class_c         = false
}
`, deviceprofileName)
}
