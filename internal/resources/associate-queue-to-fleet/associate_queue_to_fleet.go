// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package associate_queue_to_fleet

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/deadline"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AssociateQueueToFleetResource{}
var _ resource.ResourceWithImportState = &AssociateQueueToFleetResource{}

func NewAssociateQueueToFleetResource() resource.Resource {
	return &AssociateQueueToFleetResource{}
}

// AssociateQueueToFleetResource defines the resource implementation.
type AssociateQueueToFleetResource struct {
	client *deadline.Client
}

// AssociateQueueToFleetResourceModel describes the resource data model.
type AssociateQueueToFleetResourceModel struct {
	ID      types.String `tfsdk:"id"`
	FarmID  types.String `tfsdk:"farm_id"`
	FleetID types.String `tfsdk:"fleet_id"`
	QueueID types.String `tfsdk:"queue_id"`
}

func (r *AssociateQueueToFleetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_associate_queue_to_fleet"
}

func (r *AssociateQueueToFleetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Associate Member to fleet resource",
		Attributes: map[string]schema.Attribute{
			"queue_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the farm to associate the member to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"fleet_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the fleet to associate the member to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"farm_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the farm to associate the member to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the associate_queue_to_fleet.",
			},
		},
	}
}

func (r *AssociateQueueToFleetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AssociateQueueToFleetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AssociateQueueToFleetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	request := &deadline.CreateQueueFleetAssociationInput{
		FarmId:  data.FarmID.ValueStringPointer(),
		FleetId: data.FleetID.ValueStringPointer(),
		QueueId: data.QueueID.ValueStringPointer(),
	}
	////
	// Does not return the ID of the created resource: https://docs.aws.amazon.com/deadline-cloud/latest/APIReference/API_AssociateMemberToFarm.html#API_AssociateMemberToFarm_RequestSyntax
	////
	_, err := r.client.CreateQueueFleetAssociation(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, %s, got error: %s", r.typeName(), "association", err))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%s-%s-%s", data.FarmID.ValueString(), data.FleetID.ValueString(), data.QueueID.ValueString()))
	tflog.Trace(ctx, fmt.Sprintf("created %s, id: %s", r.typeName(), data.ID.ValueString()))
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssociateQueueToFleetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AssociateQueueToFleetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssociateQueueToFleetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AssociateQueueToFleetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssociateQueueToFleetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AssociateQueueToFleetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := &deadline.DeleteQueueFleetAssociationInput{
		FarmId:  data.FarmID.ValueStringPointer(),
		FleetId: data.FleetID.ValueStringPointer(),
		QueueId: data.QueueID.ValueStringPointer(),
	}
	_, err := r.client.DeleteQueueFleetAssociation(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete %s, got error: %s", r.typeName(), err))
		return
	}
}

func (r *AssociateQueueToFleetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AssociateQueueToFleetResource) typeName() string {
	return "deadline_associate_queue_to_fleet"
}
