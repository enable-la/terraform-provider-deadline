---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "deadline_queue Resource - deadline"
subcategory: ""
description: |-
  Queue resource
---

# deadline_queue (Resource)

Queue resource

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name of the queue.
- `farm_id` (String) The ID of the farm.

### Optional

- `allowed_storage_profile_ids` (List of String) The storage profile IDs to include in the queue.
- `default_budget_action` (String) The default budget action for the queue. Valid values are: 'NONE', 'STOP_SCHEDULING_AND_COMPLETE_TASKS', and 'STOP_SCHEDULING_AND_CANCEL_TASKS'.
- `description` (String) The description of the queue.
- `job_attachment_settings` (Block, Optional) (see [below for nested schema](#nestedblock--job_attachment_settings))
- `job_run_as_user` (Block, Optional) (see [below for nested schema](#nestedblock--job_run_as_user))
- `required_file_system_location_names` (List of String) The file system location name to include in the queue.
- `role_arn` (String) The IAM role ARN that workers will use while running jobs for this queue.
- `tags` (Map of Map of String) The tags to apply to the queue.

### Read-Only

- `id` (String) The ID of the queue.

<a id="nestedblock--job_attachment_settings"></a>
### Nested Schema for `job_attachment_settings`

Optional:

- `root_prefix` (String) The root prefix for the job attachment.
- `s3_bucket_name` (String) The S3 bucket name for the job attachment.


<a id="nestedblock--job_run_as_user"></a>
### Nested Schema for `job_run_as_user`

Optional:

- `posix_user` (Attributes) (see [below for nested schema](#nestedatt--job_run_as_user--posix_user))
- `run_as` (String) The user to run the job as. Either QUEUE_CONFIGURED_USER or WORKER_AGENT_USER.
- `windows_user` (Attributes) (see [below for nested schema](#nestedatt--job_run_as_user--windows_user))

<a id="nestedatt--job_run_as_user--posix_user"></a>
### Nested Schema for `job_run_as_user.posix_user`

Optional:

- `group` (String) The group to run the job as.
- `user` (String) The user to run the job as.


<a id="nestedatt--job_run_as_user--windows_user"></a>
### Nested Schema for `job_run_as_user.windows_user`

Optional:

- `password_arn` (String) The password ARN for the user to run the job as.
- `user` (String) The user to run the job as.
