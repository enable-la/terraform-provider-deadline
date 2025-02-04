---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "deadline_license_endpoint Resource - deadline"
subcategory: ""
description: |-
  LicenseEndpoint resource
---

# deadline_license_endpoint (Resource)

LicenseEndpoint resource

## Example Usage

```terraform
terraform {
  required_providers {
    awscc = {
      source  = "hashicorp/awscc"
      version = "1.2.0"
    }
  }
}


resource "awscc_deadline_license_endpoint" "example" {
  security_group_ids = ["sg-12345"]
  subnet_ids         = ["subnet-12345"]
  vpc_id             = "vpc-12345"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `security_group_ids` (List of String) The security groups that will be associated with the license endpoint
- `subnet_ids` (List of String) The subnet ids that will be associated to the license endpoint
- `vpc_id` (String) The VPC ID that the license endpoint is associated with

### Read-Only

- `id` (String) The ID of the licenseEndpoint.
