resource "deadline_farm" "test" {
  display_name = "test"
  description  = "this is a test farm"
}

resource "deadline_fleet" "test" {
  display_name = "test"
  farm_id      = deadline_farm.test.id
  description  = "this is a test farm"
}

resource "deadline_associate_member_to_fleet" "test" {
  farm_id           = deadline_farm.test.id
  fleet_id          = deadline_fleet.test.id
  member_id         = "test"
  identity_store_id = "example_identity_store"
  membership_level  = "VIEWER"
  principal_type    = "USER"
}
resource "deadline_fleet" "test" {
  farm_id          = deadline_farm.test.id
  display_name     = "test"
  description      = "This is a test fleet"
  role_arn         = "arn:aws:iam::123456789012:role/DeadlineWorkerRole"
  min_worker_count = 0
  max_worker_count = 1
  configuration {
    mode = "aws_managed"
    ec2_instance_capabilities {
      cpu_architecture = "x86_64"
      min_cpu_count    = 1
      max_cpu_count    = 2
      memory_mib_range {
        min = 1024
        max = 1024 * 4
      }
      os_family = "LINUX" // LINUX, WINDOWS
      root_ebs_volume {
        iops = 100
        size = 100
      }
    }
  }
}