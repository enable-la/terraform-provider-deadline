// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storageprofile

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/deadline"
	dltypes "github.com/aws/aws-sdk-go-v2/service/deadline/types"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &StorageProfileResource{}
var _ resource.ResourceWithImportState = &StorageProfileResource{}

func New() resource.Resource {
	return &StorageProfileResource{}
}

// StorageProfileResource defines the resource implementation.
type StorageProfileResource struct {
	client *deadline.Client
}

type StorageProfileFileSystemLocations struct {
	Name types.String `tfschema:"name"`
	Path types.String `tfschema:"path"`
	Type types.String `tfschema:"type"`
}

// StorageProfileResourceModel describes the resource data model.
type StorageProfileResourceModel struct {
	DisplayName         types.String                         `tfsdk:"display_name"`
	FarmId              types.String                         `tfsdk:"farm_id"`
	OSFamily            types.String                         `tfsdk:"os_family"`
	ID                  types.String                         `tfsdk:"id"`
	FileSystemLocations []*StorageProfileFileSystemLocations `tfsdk:"file_system_location"`
}

func (r *StorageProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_profile"
}

func (r *StorageProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "StorageProfile resource",
		Blocks: map[string]schema.Block{
			"file_system_location": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.RequiresReplace(),
					},
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the file system location",
							Required:            true,
						},
						"path": schema.StringAttribute{
							MarkdownDescription: "Path of the file system location",
							Required:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the file system location. Can be either local, or shared",
							Required:            true,
						},
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the storage profile.",
				Required:            true,
			},
			"farm_id": schema.StringAttribute{
				MarkdownDescription: "The deadline farm associated with the storage profile.",
				Required:            true,
			},
			"os_family": schema.StringAttribute{
				MarkdownDescription: "The OS family of the storage profile. Can be: windows, linux or macos",
				Required:            true,
			},

			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the storage profile.",
			},
		},
	}
}

func (r *StorageProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func determineOsProfile(inputOS string) dltypes.StorageProfileOperatingSystemFamily {
	osFamily := dltypes.StorageProfileOperatingSystemFamilyWindows
	switch inputOS {
	case "macos":
		osFamily = dltypes.StorageProfileOperatingSystemFamilyMacos
	case "windows":
		osFamily = dltypes.StorageProfileOperatingSystemFamilyWindows
	case "linux":
		osFamily = dltypes.StorageProfileOperatingSystemFamilyLinux
	}
	return osFamily
}

func getFileSystemLocations(diags diag.Diagnostics, data StorageProfileResourceModel) []dltypes.FileSystemLocation {
	var locations []dltypes.FileSystemLocation
	if len(data.FileSystemLocations) > 0 {
		for _, loc := range data.FileSystemLocations {
			cl := dltypes.FileSystemLocation{}
			switch cl.Type {
			case "local":
				cl.Type = dltypes.FileSystemLocationTypeLocal
			case "shared":
				cl.Type = dltypes.FileSystemLocationTypeShared
			}
			cl.Name = loc.Name.ValueStringPointer()
			cl.Path = loc.Path.ValueStringPointer()
			locations = append(locations, cl)
		}
	} else {
		diags.AddError("Unable to determine file system locations", "There is no file system locations available please add one or more file system locations")
	}
	return locations
}

func (r *StorageProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data StorageProfileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	osFamily := determineOsProfile(data.OSFamily.String())
	fSystemLocations := getFileSystemLocations(resp.Diagnostics, data)
	if resp.Diagnostics.HasError() {
		return
	}
	storageprofileRequest := deadline.CreateStorageProfileInput{
		DisplayName:         data.DisplayName.ValueStringPointer(),
		FarmId:              data.FarmId.ValueStringPointer(),
		OsFamily:            osFamily,
		FileSystemLocations: fSystemLocations,
	}
	storageprofileOutput, err := r.client.CreateStorageProfile(ctx, &storageprofileRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, %s, got error: %s", r.typeName(), data.DisplayName.String(), err))
		return
	}
	data.ID = types.StringValue(*storageprofileOutput.StorageProfileId)
	tflog.Trace(ctx, fmt.Sprintf("created %s", r.typeName()))
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StorageProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	storageprofileResponse, err := r.client.GetStorageProfile(ctx, &deadline.GetStorageProfileInput{
		StorageProfileId: data.ID.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read %s, got error: %s", r.typeName(), err))
		return
	}
	data.ID = types.StringValue(*storageprofileResponse.StorageProfileId)
	if len(storageprofileResponse.OsFamily.Values()) > 0 {
		data.OSFamily = types.StringValue(fmt.Sprintf("%v", storageprofileResponse.OsFamily))
	}
	if storageprofileResponse.DisplayName != nil {
		data.DisplayName = types.StringValue(*storageprofileResponse.DisplayName)
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data StorageProfileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	osFamily := determineOsProfile(data.OSFamily.String())
	updateRequest := deadline.UpdateStorageProfileInput{
		StorageProfileId: data.ID.ValueStringPointer(),
		FarmId:           data.FarmId.ValueStringPointer(),
		DisplayName:      data.DisplayName.ValueStringPointer(),
		OsFamily:         osFamily,
	}
	_, err := r.client.UpdateStorageProfile(ctx, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update %s, got error: %s", r.typeName(), err))
		return
	}
	data.DisplayName = types.StringValue(*updateRequest.DisplayName)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StorageProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceRequest := &deadline.DeleteStorageProfileInput{
		StorageProfileId: data.ID.ValueStringPointer(),
	}
	_, err := r.client.DeleteStorageProfile(ctx, deleteResourceRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete %s, got error: %s", r.typeName(), err))
		return
	}
}

func (r *StorageProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *StorageProfileResource) typeName() string {
	return "deadline_storage_profile"
}
