// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	"deadline": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func TestAccQueueResource(t *testing.T) {
	testRoleARN := os.Getenv("TEST_DEADLINE_ROLE_ARN")
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
}

resource "deadline_fleet" "test" {
  farm_id = "${deadline_farm.test.id}"
	  display_name = %[1]q	
	  description = %[2]q
	role_arn = %[3]q
  min_worker_count = 0
  max_worker_count = 1
  configuration {
    mode                   = "aws_managed"
    ec2_instance_capabilities { 
    cpu_architecture       = "x86_64"
    min_cpu_count          = 1
    max_cpu_count          = 2
    memory_mib_range            {
    min = 1024
    max = 1024 * 4
}
    os_family              = "LINUX" // LINUX, WINDOWS
    root_ebs_volume {
      iops = 100
      size = 100
    }
}
  }
}

`, displayName, description, roleARN)
}
