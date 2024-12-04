provider "aws" {

}

provider "deadline" {

}

resource "deadline_farm" "test" {
  display_name = "example-farm-01"
  description  = "this is a test farm"
}

data "aws_iam_role" "fleet_role" {
  name = "DeadlineFleetRole"
}

resource "deadline_monitor" "test" {
  display_name                 = "example-monitor-01"
  subdomain                    = "example-monitor-01" // subdomain.region.deadlinecloud.amazonaws.com
  role_arn                     = data.aws_iam_role.fleet_role.arn
  identity_center_instance_arn = "arn:aws:identitycenter:us-west-2:123456789012:instance/12345678-1234-1234-1234-123456789012"
}


resource "deadline_fleet" "test" {
  farm_id          = deadline_farm.test.id
  display_name     = "example-fleet-01"
  description      = "This would be an example fleet"
  role_arn         = data.aws_iam_role.fleet_role.arn
  min_worker_count = 0
  max_worker_count = 1
  configuration {
    mode                   = "aws_managed"
    cpu_architecture       = "x86_64"
    cpu_count              = 4
    memory_mib             = 1024 * 4
    os_family              = "LINUX" // LINUX, WINDOWS
    allowed_instance_types = ["c5.large"]
    acclerator_capabilities {
      selections = ["t4"] // t4, a10g, l4, l40s
      count      = 1
    }
    root_ebs_volume {
      iops = 100
      size = 100
    }
  }
}

resource "deadlne_queue" "test" {
  display_name = "example-queue-01"
  description  = "This would be an example queue"
  farm_id      = deadline_farm.test.id
  fleet_id     = deadline_fleet.test.id
  min_slots    = 1
  max_slots    = 1
  priority     = 100
}

resource "deadline_associate_queue_to_fleet" "test" {
  farm_id  = deadline_farm.test.id
  fleet_id = deadline_fleet.test.id
  queue_id = deadline_queue.test.id
}


resource "deadline_associate_member_to_farm" "jdoe" {
  farm_id           = deadline_farm.test.id
  member_id         = "jdoe"
  identity_store_id = "example_identity_store"
  membership_level  = "VIEWER" // VIEWER, CONTRIBUTOR, MANAGER, OWNER
  principal_type    = "USER"   // USER, GROUP
}

resource "deadline_associate_member_to_fleet" "jdoe" {
  farm_id          = deadline_farm.test.id
  fleet_id         = deadline_fleet.test.id
  member_id        = "jdoe"
  membership_level = "VIEWER" // VIEWER, CONTRIBUTOR, MANAGER, OWNER
  principal_type   = "USER"   // USER, GROUP
}