// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package queue_environment

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
var _ resource.Resource = &QueueEnvironmentResource{}
var _ resource.ResourceWithImportState = &QueueEnvironmentResource{}

func New() resource.Resource {
	return &QueueEnvironmentResource{}
}

// QueueEnvironmentResource defines the resource implementation.
type QueueEnvironmentResource struct {
	client *deadline.Client
}

// QueueEnvironmentResourceModel describes the resource data model.
type QueueEnvironmentResourceModel struct {
	QueueId      types.String `tfsdk:"queue_id"`
	FarmId       types.String `tfsdk:"farm_id	"`
	Priority     types.Int32  `tfsdk:"priority"`
	TemplateType types.String `tfsdk:"template_type"`
	Template     types.String `tfsdk:"template"`
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
}

func (r *QueueEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_queue_environment"
}

func (r *QueueEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "QueueEnvironment resource",
		Attributes: map[string]schema.Attribute{
			"farm_id": schema.StringAttribute{
				MarkdownDescription: "The display name of the queueEnvironment.",
				Required:            true,
			},
			"priority": schema.Int32Attribute{
				MarkdownDescription: "sets the priority of the environments in the queue from 0 to 10,000, where 0 is the highest priority. If two environments share the same priority value, the environment created first takes higher priority.",
				Required:            true,
			},
			"queue_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the queue.",
				Required:            true,
			},
			"template": schema.StringAttribute{
				Required:    true,
				Description: "The environment template to use in the queue. See examples here: https://github.com/aws-deadline/deadline-cloud-samples/blob/mainline/README.md",
			},
			"template_type": schema.StringAttribute{
				Required:    true,
				Description: "The environment template to use in the queue. Can be either json or yaml",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the queueEnvironment.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the QueueEnvironment.",
			},
		},
	}
}

func (r *QueueEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *QueueEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data QueueEnvironmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	queueEnvironmentRequest := deadline.CreateQueueEnvironmentInput{}
	queueEnvironmentOutput, err := r.client.CreateQueueEnvironment(ctx, &queueEnvironmentRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, got error: %s", r.typeName(), err))
		return
	}
	data.ID = types.StringValue(*queueEnvironmentOutput.QueueEnvironmentId)
	tflog.Trace(ctx, fmt.Sprintf("created %s", r.typeName()))
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueueEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data QueueEnvironmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	queueEnvironmentResponse, err := r.client.GetQueueEnvironment(ctx, &deadline.GetQueueEnvironmentInput{
		QueueEnvironmentId: data.ID.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s, got error: %s", r.typeName(), err))
		return
	}
	data.ID = types.StringValue(*queueEnvironmentResponse.QueueEnvironmentId)
	data.Name = types.StringValue(*queueEnvironmentResponse.Name)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueueEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data QueueEnvironmentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	templateType := dltypes.EnvironmentTemplateTypeJson
	switch data.TemplateType.ValueString() {
	case "json":
		templateType = dltypes.EnvironmentTemplateTypeJson
	case "yaml":
		templateType = dltypes.EnvironmentTemplateTypeYaml
	}

	updateRequest := deadline.UpdateQueueEnvironmentInput{
		QueueEnvironmentId: data.ID.ValueStringPointer(),
		FarmId:             data.FarmId.ValueStringPointer(),
		Template:           data.Template.ValueStringPointer(),
		TemplateType:       templateType,
	}
	_, err := r.client.UpdateQueueEnvironment(ctx, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update %s, got error: %s", r.typeName(), err))
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QueueEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data QueueEnvironmentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceRequest := &deadline.DeleteQueueEnvironmentInput{
		QueueEnvironmentId: data.ID.ValueStringPointer(),
	}
	_, err := r.client.DeleteQueueEnvironment(ctx, deleteResourceRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete %s, got error: %s", r.typeName(), err))
		return
	}
}

func (r *QueueEnvironmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *QueueEnvironmentResource) typeName() string {
	return "deadline_queue_environment"
}
