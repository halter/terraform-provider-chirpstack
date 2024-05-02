// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApplicationResourceConfig("application-one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("chirpstack_application.test", "id"),
					resource.TestCheckResourceAttr("chirpstack_application.test", "name", "application-one"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "chirpstack_application.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccApplicationResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("chirpstack_application.test", "name", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApplicationResourceConfig(applicationName string) string {
	return fmt.Sprintf(`
resource "chirpstack_tenant" "test" {
  name = "test_tenant"
}
resource "chirpstack_application" "test" {
  tenant_id = chirpstack_tenant.test.id
  name      = %[1]q
}
`, applicationName)
}
