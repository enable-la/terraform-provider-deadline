// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package farm

import (
	"fmt"
	"github.com/enable-la/terraform-provider-aws-deadline/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"deadline": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func TestAccFarmResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFarmResourceConfig("test", "this is a test farm"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("deadline_farm.test", "display_name", "test"),
					resource.TestCheckResourceAttr("deadline_farm.test", "description", "this is a test farm"),
				),
			},
		},
	})
}

func testAccFarmResourceConfig(displayName string, description string) string {
	return fmt.Sprintf(`
resource "deadline_farm" "test" {
  display_name = %[1]q
  description = %[2]q
}
`, displayName, description)
}
