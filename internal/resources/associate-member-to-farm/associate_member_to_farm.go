// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package associate_member_to_farm

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/deadline"
	dltypes "github.com/aws/aws-sdk-go-v2/service/deadline/types"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &AssociateMemberToFarmResource{}
var _ resource.ResourceWithImportState = &AssociateMemberToFarmResource{}

func New() resource.Resource {
	return &AssociateMemberToFarmResource{}
}

// AssociateMemberToFarmResource defines the resource implementation.
type AssociateMemberToFarmResource struct {
	client *deadline.Client
}

// AssociateMemberToFarmResourceModel describes the resource data model.
type AssociateMemberToFarmResourceModel struct {
	ID              types.String `tfsdk:"id"`
	FarmID          types.String `tfsdk:"farm_id"`
	IdentityStoreID types.String `tfsdk:"identity_store_id"`
	PrincipalID     types.String `tfsdk:"principal_id"`
	PrincipalType   types.String `tfsdk:"principal_type"`
	MembershipLevel types.String `tfsdk:"membership_level"`
}

func (r *AssociateMemberToFarmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_associate_member_to_farm"
}

func (r *AssociateMemberToFarmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Associate Member to Farm resource",
		Attributes: map[string]schema.Attribute{
			"farm_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the farm to associate the member to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"identity_store_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the identity store that the member belongs to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"principal_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the principal to associate to the farm",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"principal_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of principal to associate to the farm. Valid values are `USER` and `GROUP`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"membership_level": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The membership level of the principal to associate to the farm. Valid values are `VIEWER`, `CONTRIBUTOR`, `OWNER` and `MANAGER`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the associate_member_to_farm.",
			},
		},
	}
}

func (r *AssociateMemberToFarmResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AssociateMemberToFarmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AssociateMemberToFarmResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	request := &deadline.AssociateMemberToFarmInput{
		FarmId:          data.FarmID.ValueStringPointer(),
		PrincipalId:     data.PrincipalID.ValueStringPointer(),
		IdentityStoreId: data.IdentityStoreID.ValueStringPointer(),
	}
	if data.PrincipalType.ValueString() != "USER" && data.PrincipalType.ValueString() != "GROUP" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Invalid principal type: %s", data.PrincipalType.ValueString()))
		return
	}
	switch data.MembershipLevel.ValueString() {
	case "VIEWER":
		request.MembershipLevel = dltypes.MembershipLevelViewer
	case "CONTRIBUTOR":
		request.MembershipLevel = dltypes.MembershipLevelContributor
	case "OWNER":
		request.MembershipLevel = dltypes.MembershipLevelOwner
	case "MANAGER":
		request.MembershipLevel = dltypes.MembershipLevelManager
	}
	if data.PrincipalType.ValueString() == "USER" {
		request.PrincipalType = dltypes.DeadlinePrincipalTypeUser
	} else {
		request.PrincipalType = dltypes.DeadlinePrincipalTypeGroup
	}

	////
	// Does not return the ID of the created resource: https://docs.aws.amazon.com/deadline-cloud/latest/APIReference/API_AssociateMemberToFarm.html#API_AssociateMemberToFarm_RequestSyntax
	////
	_, err := r.client.AssociateMemberToFarm(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, %s, got error: %s", r.typeName(), "association", err))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%s-%s-%s", data.FarmID.ValueString(), data.PrincipalID.ValueString(), data.IdentityStoreID.ValueString()))
	tflog.Trace(ctx, fmt.Sprintf("created %s, id: %s", r.typeName(), data.ID.ValueString()))
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssociateMemberToFarmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AssociateMemberToFarmResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssociateMemberToFarmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AssociateMemberToFarmResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AssociateMemberToFarmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AssociateMemberToFarmResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	request := &deadline.DisassociateMemberFromFarmInput{
		FarmId:      data.ID.ValueStringPointer(),
		PrincipalId: data.PrincipalID.ValueStringPointer(),
	}
	_, err := r.client.DisassociateMemberFromFarm(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete %s, got error: %s", r.typeName(), err))
		return
	}
}

func (r *AssociateMemberToFarmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *AssociateMemberToFarmResource) typeName() string {
	return "deadline_associate_member_to_farm"
}
