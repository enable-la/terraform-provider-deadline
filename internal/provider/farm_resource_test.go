// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

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
