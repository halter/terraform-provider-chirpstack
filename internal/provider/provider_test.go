// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"chirpstack": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	if os.Getenv("CHIRPSTACK_HOST") == "" {
		t.Skip("environment variable CHIRPSTACK_HOST is not set")
	}
	if os.Getenv("CHIRPSTACK_PORT") == "" {
		t.Skip("environment variable CHIRPSTACK_PORT is not set")
	}
	if os.Getenv("CHIRPSTACK_KEY") == "" {
		t.Skip("environment variable CHIRPSTACK_KEY is not set")
	}
}
