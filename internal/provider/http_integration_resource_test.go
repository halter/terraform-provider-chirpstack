// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHttpIntegrationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccHttpIntegrationResourceConfig("http://localhost"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("chirpstack_http_integration.test", "id"),
					resource.TestCheckResourceAttrSet("chirpstack_http_integration.test", "application_id"),
					resource.TestCheckResourceAttr("chirpstack_http_integration.test", "event_endpoint_url", "http://localhost"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "chirpstack_http_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccHttpIntegrationResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("chirpstack_http_integration.test", "event_endpoint_url", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccHttpIntegrationResourceConfig(endpoint string) string {
	return fmt.Sprintf(`
resource "chirpstack_tenant" "test" {
  name = "test_tenant"
}
resource "chirpstack_application" "test" {
  tenant_id = chirpstack_tenant.test.id
  name = "test_app"
}
resource "chirpstack_http_integration" "test" {
  application_id = chirpstack_application.test.id
  encoding = "JSON"
  event_endpoint_url = %[1]q
}
`, endpoint)
}
