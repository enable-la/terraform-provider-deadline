// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFleetResource(t *testing.T) {
	testRoleARN := os.Getenv("TEST_ROLE_ARN")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFleetResourceConfig("test", "this is a test farm", testRoleARN),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("deadline_fleet.test", "display_name", "test"),
					resource.TestCheckResourceAttr("deadline_fleet.test", "description", "this is a test farm"),
				),
			},
		},
	})
}

func testAccFleetResourceConfig(displayName string, description string, roleARN string) string {
	return fmt.Sprintf(`
resource "deadline_farm" "test" {
	display_name = %[1]q
    description  = "this is a farm"
}
resource "deadline_fleet" "test" {
  farm_id = "${deadline_farm.test.id}"
  display_name = %[1]q
  description = %[2]q
  role_arn = %[3]q
  configuration = "aws_managed"
  min_worker_count = "0"
  max_worker_count = "1"
}
`, displayName, description, roleARN)
}
