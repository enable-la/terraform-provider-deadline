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