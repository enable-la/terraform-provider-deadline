// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQueueResource(t *testing.T) {
	testRoleARN := os.Getenv("TEST_ROLE_ARN")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccQueueResourceConfig("test", "this is a test", testRoleARN),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("deadline_queue.test", "display_name", "test"),
					resource.TestCheckResourceAttr("deadline_queue.test", "description", "this is a test"),
					resource.TestCheckResourceAttr("deadline_queue.test", "role_arn", testRoleARN),
				),
			},
		},
	})
}

func testAccQueueResourceConfig(displayName string, description string, roleARN string) string {
	return fmt.Sprintf(`
resource "deadline_farm" "test" {
	display_name = %[1]q
    description  = "this is a farm"
}
resource "deadline_queue" "test" {
  farm_id = "${deadline_farm.test.id}"
  display_name = %[1]q
  description = %[2]q
  role_arn = %[3]q
  allowed_storage_profile_ids = ["storage_profile_id"]
}
`, displayName, description, roleARN)
}
