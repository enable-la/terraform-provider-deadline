// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package license_endpoint

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/deadline"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &LicenseEndpointResource{}
var _ resource.ResourceWithImportState = &LicenseEndpointResource{}

func New() resource.Resource {
	return &LicenseEndpointResource{}
}

// LicenseEndpointResource defines the resource implementation.
type LicenseEndpointResource struct {
	client *deadline.Client
}

// LicenseEndpointResourceModel describes the resource data model.
type LicenseEndpointResourceModel struct {
	SecurityGroupIds []types.String `tfsdk:"security_group_ids"`
	SubnetIds        []types.String `tfsdk:"subnet_ids"`
	VpcId            types.String   `tfsdk:"vpc_id"`
	ID               types.String   `tfsdk:"id"`
}

func (r *LicenseEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license_endpoint"
}

func (r *LicenseEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "LicenseEndpoint resource",
		Attributes: map[string]schema.Attribute{
			"security_group_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "The security groups that will be associated with the license endpoint",
			},
			"subnet_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The subnet ids that will be associated to the license endpoint",
				Required:            true,
			},
			"vpc_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The VPC ID that the license endpoint is associated with",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the licenseEndpoint.",
			},
		},
	}
}

func (r *LicenseEndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LicenseEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LicenseEndpointResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	subnets := []string{}
	for _, subnet := range data.SubnetIds {
		subnets = append(subnets, subnet.String())
	}
	sgIds := []string{}
	for _, sgId := range data.SecurityGroupIds {
		sgIds = append(sgIds, sgId.String())
	}
	licenseEndpointRequest := deadline.CreateLicenseEndpointInput{
		VpcId:            data.VpcId.ValueStringPointer(),
		SubnetIds:        subnets,
		SecurityGroupIds: sgIds,
	}
	licenseEndpointOutput, err := r.client.CreateLicenseEndpoint(ctx, &licenseEndpointRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, got error: %s", r.typeName(), err))
		return
	}
	data.ID = types.StringValue(*licenseEndpointOutput.LicenseEndpointId)
	tflog.Trace(ctx, fmt.Sprintf("created %s", r.typeName()))
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LicenseEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LicenseEndpointResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	licenseEndpointResponse, err := r.client.GetLicenseEndpoint(ctx, &deadline.GetLicenseEndpointInput{
		LicenseEndpointId: data.ID.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s, got error: %s", r.typeName(), err))
		return
	}
	data.ID = types.StringValue(*licenseEndpointResponse.LicenseEndpointId)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LicenseEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LicenseEndpointResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LicenseEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LicenseEndpointResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceRequest := &deadline.DeleteLicenseEndpointInput{
		LicenseEndpointId: data.ID.ValueStringPointer(),
	}
	_, err := r.client.DeleteLicenseEndpoint(ctx, deleteResourceRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete %s, got error: %s", r.typeName(), err))
		return
	}
}

func (r *LicenseEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *LicenseEndpointResource) typeName() string {
	return "deadline_license_endpoint"
}
