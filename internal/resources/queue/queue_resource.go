package queue

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/deadline"
	dltypes "github.com/aws/aws-sdk-go-v2/service/deadline/types"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &QueueResource{}
var _ resource.ResourceWithImportState = &QueueResource{}

func NewQueueResource() resource.Resource {
	return &QueueResource{}
}

// QueueResource defines the resource implementation.
type QueueResource struct {
	client *deadline.Client
}

type QueueResourceConfigurationModel struct {
	Mode                    types.String                              `tfsdk:"mode"`
	Ec2MarketType           types.String                              `tfsdk:"ec2_market_type"`
	Ec2InstanceCapabilities QueueResourceEc2InstanceCapabilitiesModel `tfsdk:"ec2_instance_capabilities"`
}

type QueueResourceEc2InstanceCapabilitiesModel struct {
	CpuArchitecture         types.String `tfsdk:"cpu_architecture"`
	MinCpuCount             types.Int32  `tfsdk:"min_cpu_count"`
	MaxCpuCount             types.Int32  `tfsdk:"max_cpu_count"`
	MemoryMib               types.Int32  `tfsdk:"memory_mib"`
	OsFamily                types.String `tfsdk:"os_family"`
	AllowedInstanceType     []types.String
	ExcludeInstanceType     []types.String                                                   `tfsdk:"exclude_instance_types"`
	AcceleratorCapabilities QueueResourceEc2InstanceCapabilitiesAcceleratorCapabilitiesModel `tfsdk:"accelerator_capabilities"`
}

func (r *QueueResourceEc2InstanceCapabilitiesModel) Value() {

}

type QueueResourceEc2InstanceCapabilitiesAcceleratorCapabilitiesModel struct {
	Selections types.ListType `tfsdk:"selections"`
	Count      types.Int32    `tfsdk:"count"`
}
type QueueResourceJobAttachmentSettingsModel struct {
	RootPrefix   types.String `tfsdk:"root_prefix"`
	S3BucketName types.String `tfsdk:"s3_bucket_name"`
}
type QueueResourceJobRunAsUserPosixUserModel struct {
	Group types.String `tfsdk:"group"`
	User  types.String `tfsdk:"user"`
}
type QueueResourceJobRunAsUserWindowsUserModel struct {
	PasswordArn types.String `tfsdk:"password_arn"`
	User        types.String `tfsdk:"user"`
}
type QueueResourceJobRunAsUserModel struct {
	PosixUser   *QueueResourceJobRunAsUserPosixUserModel   `tfsdk:"posix_user"`
	WindowsUser *QueueResourceJobRunAsUserWindowsUserModel `tfsdk:"windows_user"`
	RunAs       types.String                               `tfsdk:"run_as"`
}

// QueueResourceModel describes the resource data model.
type QueueResourceModel struct {
	DisplayName                     types.String                             `tfsdk:"display_name"`
	Description                     types.String                             `tfsdk:"description"`
	FarmId                          types.String                             `tfsdk:"farm_id"`
	RoleArn                         types.String                             `tfsdk:"role_arn"`
	ID                              types.String                             `tfsdk:"id"`
	AllowedStorageProfileIds        []types.String                           `tfsdk:"allowed_storage_profile_ids"`
	DefaultBudgetAction             types.String                             `tfsdk:"default_budget_action"`
	JobAttachmentSettings           *QueueResourceJobAttachmentSettingsModel `tfsdk:"job_attachment_settings"`
	JobRunAsUser                    *QueueResourceJobRunAsUserModel          `tfsdk:"job_run_as_user"`
	RequiredFileSystemLocationNames []types.String                           `tfsdk:"required_file_system_location_names"`
	Tags                            map[string]types.String                  `tfsdk:"tags"`
}

func (r *QueueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_queue"
}

func (r *QueueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Queue resource",
		Blocks: map[string]schema.Block{
			"job_attachment_settings": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"root_prefix": schema.StringAttribute{
						Optional:    true,
						Description: "The root prefix for the job attachment.",
					},
					"s3_bucket_name": schema.StringAttribute{
						Optional:    true,
						Description: "The S3 bucket name for the job attachment.",
					},
				},
			},
			"job_run_as_user": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{
					"posix_user": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"group": schema.StringAttribute{
								Optional:    true,
								Description: "The group to run the job as.",
							},
							"user": schema.StringAttribute{
								Optional:    true,
								Description: "The user to run the job as.",
							},
						},
					},
					"windows_user": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"password_arn": schema.StringAttribute{
								Optional:    true,
								Description: "The password ARN for the user to run the job as.",
							},
							"user": schema.StringAttribute{
								Optional:    true,
								Description: "The user to run the job as.",
							},
						},
					},
					"run_as": schema.StringAttribute{
						Optional:    true,
						Description: "The user to run the job as. Either QUEUE_CONFIGURED_USER or WORKER_AGENT_USER.",
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the queue.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the queue.",
			},
			"farm_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the farm.",
			},
			"default_budget_action": schema.StringAttribute{
				Optional:    true,
				Description: "The default budget action for the queue. Valid values are: 'NONE', 'STOP_SCHEDULING_AND_COMPLETE_TASKS', and 'STOP_SCHEDULING_AND_CANCEL_TASKS'.",
			},
			"allowed_storage_profile_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "The storage profile IDs to include in the queue.",
			},
			"required_file_system_location_names": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "The file system location name to include in the queue.",
			},
			"role_arn": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The IAM role ARN that workers will use while running jobs for this queue.",
			},
			"tags": schema.MapAttribute{
				ElementType:         types.MapType{},
				Optional:            true,
				MarkdownDescription: "The tags to apply to the queue.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the queue.",
			},
		},
	}
}

func (r *QueueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*deadline.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *QueueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data QueueResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := &deadline.CreateQueueInput{
		FarmId:      data.FarmId.ValueStringPointer(),
		DisplayName: data.DisplayName.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		RoleArn:     data.RoleArn.ValueStringPointer(),
	}
	if len(data.AllowedStorageProfileIds) > 0 {
		allowedStorageProfileIds := make([]string, len(data.AllowedStorageProfileIds))
		for k, v := range data.AllowedStorageProfileIds {
			allowedStorageProfileIds[k] = v.String()
		}
		createRequest.AllowedStorageProfileIds = allowedStorageProfileIds
	}
	if data.JobAttachmentSettings.RootPrefix.ValueStringPointer() != nil && data.JobAttachmentSettings.S3BucketName.ValueStringPointer() != nil {
		if data.JobAttachmentSettings.RootPrefix.ValueString() == "" || data.JobAttachmentSettings.S3BucketName.ValueString() == "" {
			createRequest.JobAttachmentSettings = &dltypes.JobAttachmentSettings{
				RootPrefix:   data.JobAttachmentSettings.RootPrefix.ValueStringPointer(),
				S3BucketName: data.JobAttachmentSettings.S3BucketName.ValueStringPointer(),
			}
		}
	}
	if data.JobRunAsUser.PosixUser.User.ValueString() != "" {
		createRequest.JobRunAsUser = &dltypes.JobRunAsUser{
			Posix: &dltypes.PosixUser{
				Group: data.JobRunAsUser.PosixUser.Group.ValueStringPointer(),
				User:  data.JobRunAsUser.PosixUser.User.ValueStringPointer(),
			},
		}
	}
	if data.JobRunAsUser.WindowsUser.User.ValueString() != "" {
		createRequest.JobRunAsUser = &dltypes.JobRunAsUser{
			Windows: &dltypes.WindowsUser{
				PasswordArn: data.JobRunAsUser.WindowsUser.PasswordArn.ValueStringPointer(),
				User:        data.JobRunAsUser.WindowsUser.User.ValueStringPointer(),
			},
		}
	}
	if data.JobRunAsUser.RunAs.ValueString() != "" {
		if data.JobRunAsUser.RunAs.ValueString() == "QUEUE_CONFIGURED_USER" {
			createRequest.JobRunAsUser.RunAs = dltypes.RunAsQueueConfiguredUser
		} else if data.JobRunAsUser.RunAs.ValueString() == "WORKER_AGENT_USER" {
			createRequest.JobRunAsUser.RunAs = dltypes.RunAsWorkerAgentUser
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid value for run_as, got %s", data.JobRunAsUser.RunAs.ValueString()))
			return
		}
	}
	createOutput, err := r.client.CreateQueue(ctx, createRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, %s, got error: %s", r.typeName(), data.DisplayName.String(), err))
		return
	}
	data.ID = types.StringValue(*createOutput.QueueId)
	tflog.Trace(ctx, fmt.Sprintf("created %s", r.typeName()))
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data QueueResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	getResponse, err := r.client.GetQueue(ctx, &deadline.GetQueueInput{
		QueueId: data.ID.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s, got error: %s", r.typeName(), err))
		return
	}
	data.Description = types.StringValue(*getResponse.Description)
	data.DisplayName = types.StringValue(*getResponse.DisplayName)
	data.FarmId = types.StringValue(*getResponse.FarmId)
	data.RoleArn = types.StringValue(*getResponse.RoleArn)
	data.JobAttachmentSettings.RootPrefix = types.StringValue(*getResponse.JobAttachmentSettings.RootPrefix)
	data.JobAttachmentSettings.S3BucketName = types.StringValue(*getResponse.JobAttachmentSettings.S3BucketName)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data QueueResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	updateRequest := &deadline.UpdateQueueInput{
		FarmId:      data.FarmId.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		DisplayName: data.DisplayName.ValueStringPointer(),
	}
	if len(data.AllowedStorageProfileIds) > 0 {
		allowedStorageProfileIds := make([]string, len(data.AllowedStorageProfileIds))
		for k, v := range data.AllowedStorageProfileIds {
			allowedStorageProfileIds[k] = v.String()
		}
		updateRequest.AllowedStorageProfileIdsToAdd = allowedStorageProfileIds
	}
	if data.JobAttachmentSettings.RootPrefix.ValueStringPointer() != nil && data.JobAttachmentSettings.S3BucketName.ValueStringPointer() != nil {
		if data.JobAttachmentSettings.RootPrefix.ValueString() == "" || data.JobAttachmentSettings.S3BucketName.ValueString() == "" {
			updateRequest.JobAttachmentSettings = &dltypes.JobAttachmentSettings{
				RootPrefix:   data.JobAttachmentSettings.RootPrefix.ValueStringPointer(),
				S3BucketName: data.JobAttachmentSettings.S3BucketName.ValueStringPointer(),
			}
		}
	}
	if data.JobRunAsUser.PosixUser.User.ValueString() != "" {
		updateRequest.JobRunAsUser = &dltypes.JobRunAsUser{
			Posix: &dltypes.PosixUser{
				Group: data.JobRunAsUser.PosixUser.Group.ValueStringPointer(),
				User:  data.JobRunAsUser.PosixUser.User.ValueStringPointer(),
			},
		}
	}
	if data.JobRunAsUser.WindowsUser.User.ValueString() != "" {
		updateRequest.JobRunAsUser = &dltypes.JobRunAsUser{
			Windows: &dltypes.WindowsUser{
				PasswordArn: data.JobRunAsUser.WindowsUser.PasswordArn.ValueStringPointer(),
				User:        data.JobRunAsUser.WindowsUser.User.ValueStringPointer(),
			},
		}
	}
	if data.JobRunAsUser.RunAs.ValueString() != "" {
		if data.JobRunAsUser.RunAs.ValueString() == "QUEUE_CONFIGURED_USER" {
			updateRequest.JobRunAsUser.RunAs = dltypes.RunAsQueueConfiguredUser
		} else if data.JobRunAsUser.RunAs.ValueString() == "WORKER_AGENT_USER" {
			updateRequest.JobRunAsUser.RunAs = dltypes.RunAsWorkerAgentUser
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid value for run_as, got %s", data.JobRunAsUser.RunAs.ValueString()))
			return
		}
	}
	_, err := r.client.UpdateQueue(ctx, updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update %s, got error: %s", r.typeName(), err))
		return
	}
	data.Description = types.StringValue(*updateRequest.Description)
	data.DisplayName = types.StringValue(*updateRequest.DisplayName)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data QueueResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceRequest := &deadline.DeleteQueueInput{
		QueueId: data.ID.ValueStringPointer(),
	}
	_, err := r.client.DeleteQueue(ctx, deleteResourceRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete %s, got error: %s", r.typeName(), err))
		return
	}
}

func (r *QueueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *QueueResource) typeName() string {
	return "deadline_queue"
}
