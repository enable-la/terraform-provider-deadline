// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package associate_member_to_farm

import (
	"fmt"
	"github.com/enable-la/terraform-provider-aws-deadline/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
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

func TestAccAssociateMemberToFarmResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAssociateMemberToFarmResourceConfig("test", "this is a test farm", os.Getenv("TEST_PRINCIPAL_ID"), os.Getenv("TEST_IDENTITY_STORE_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("deadline_associate_member_to_farm.test", "principal_id", os.Getenv("TEST_PRINCIPAL_ID")),
					resource.TestCheckResourceAttr("deadline_associate_member_to_farm.test", "identity_store_id", os.Getenv("TEST_IDENTITY_STORE_ID")),
				),
			},
		},
	})
}

func testAccAssociateMemberToFarmResourceConfig(displayName string, description string, principalID string, identityStoreId string) string {
	return fmt.Sprintf(`
resource "deadline_farm" "test" {
  display_name = %[1]q
  description  = %[2]q
}

resource "deadline_associate_member_to_farm" "test" {
  farm_id = "${deadline_farm.test.id}"
  principal_id = %[3]q
  identity_store_id = %[4]q
}
`, displayName, description, principalID, identityStoreId)
}
