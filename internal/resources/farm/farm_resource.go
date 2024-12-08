package farm

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
var _ resource.Resource = &FarmResource{}
var _ resource.ResourceWithImportState = &FarmResource{}

func NewFarmResource() resource.Resource {
	return &FarmResource{}
}

// FarmResource defines the resource implementation.
type FarmResource struct {
	client *deadline.Client
}

// FarmResourceModel describes the resource data model.
type FarmResourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Description types.String `tfsdk:"description"`
	ID          types.String `tfsdk:"id"`
}

func (r *FarmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_farm"
}

func (r *FarmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Farm resource",
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the farm.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the farm.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the farm.",
			},
		},
	}
}

func (r *FarmResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FarmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FarmResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	farmRequest := deadline.CreateFarmInput{
		DisplayName: data.DisplayName.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
	}
	farmOutput, err := r.client.CreateFarm(ctx, &farmRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, %s, got error: %s", r.typeName(), data.DisplayName.String(), err))
		return
	}
	data.ID = types.StringValue(*farmOutput.FarmId)
	tflog.Trace(ctx, fmt.Sprintf("created %s", r.typeName()))
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FarmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FarmResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	farmResponse, err := r.client.GetFarm(ctx, &deadline.GetFarmInput{
		FarmId: data.ID.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s, got error: %s", r.typeName(), err))
		return
	}
	data.Description = types.StringValue(*farmResponse.Description)
	data.DisplayName = types.StringValue(*farmResponse.DisplayName)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FarmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FarmResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	updateRequest := deadline.UpdateFarmInput{
		FarmId:      data.ID.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		DisplayName: data.DisplayName.ValueStringPointer(),
	}
	_, err := r.client.UpdateFarm(ctx, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update %s, got error: %s", r.typeName(), err))
		return
	}
	data.Description = types.StringValue(*updateRequest.Description)
	data.DisplayName = types.StringValue(*updateRequest.DisplayName)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FarmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FarmResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceRequest := &deadline.DeleteFarmInput{
		FarmId: data.ID.ValueStringPointer(),
	}
	_, err := r.client.DeleteFarm(ctx, deleteResourceRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete %s, got error: %s", r.typeName(), err))
		return
	}
}

func (r *FarmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *FarmResource) typeName() string {
	return "deadline_farm"
}
