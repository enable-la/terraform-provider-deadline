resource "deadline_farm" "test" {
  display_name = "test"
  description  = "this is a test farm"
}

resource "deadline_queue" "test" {
  farm_id      = deadline_farm.test.id
  display_name = "test queue"
  description  = "This is a test queue"
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

resource "deadline_queue_environment" "test" {
  farm_id  = deadline_farm.test.id
  queue_id = deadline_queue.test.id
  //
  // https://github.com/aws-deadline/deadline-cloud-samples/blob/mainline/README.md
  //
  template      = file("${path.module}/environment.yaml")
  template_type = "yaml"
}