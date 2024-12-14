// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fleet

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
var _ resource.Resource = &FleetResource{}
var _ resource.ResourceWithImportState = &FleetResource{}

func New() resource.Resource {
	return &FleetResource{}
}

// FleetResource defines the resource implementation.
type FleetResource struct {
	client *deadline.Client
}

type FleetResourceConfigurationModel struct {
	Mode                    types.String                               `tfsdk:"mode"`
	Ec2MarketType           types.String                               `tfsdk:"ec2_market_type"`
	Ec2InstanceCapabilities *FleetResourceEc2InstanceCapabilitiesModel `tfsdk:"ec2_instance_capabilities"`
}
type FleetResourceEc2InstanceCapabilitiesMemoryRangeeeModel struct {
	Min types.Int32 `tfsdk:"min"`
	Max types.Int32 `tfsdk:"max"`
}
type FleetResourceEc2InstanceCapabilitiesModel struct {
	CpuArchitecture         types.String                                                      `tfsdk:"cpu_architecture"`
	MinCpuCount             types.Int32                                                       `tfsdk:"min_cpu_count"`
	MaxCpuCount             types.Int32                                                       `tfsdk:"max_cpu_count"`
	MemoryMibRange          *FleetResourceEc2InstanceCapabilitiesMemoryRangeeeModel           `tfsdk:"memory_mib_range"`
	OsFamily                types.String                                                      `tfsdk:"os_family"`
	AllowedInstanceType     types.List                                                        `tfsdk:"allowed_instance_types"`
	ExcludeInstanceType     types.List                                                        `tfsdk:"exclude_instance_types"`
	AcceleratorCapabilities *FleetResourceEc2InstanceCapabilitiesAcceleratorCapabilitiesModel `tfsdk:"accelerator_capabilities"`
	RootEBSVolume           *FleetResourceEc2InstanceCapabilitiesRootEBSVolumeModel           `tfsdk:"root_ebs_volume"`
}

type FleetResourceEc2InstanceCapabilitiesRootEBSVolumeModel struct {
	IOPs       types.Int32 `tfsdk:"iops"`
	Size       types.Int32 `tfsdk:"size"`
	Throughput types.Int32 `tfsdk:"throughput"`
}

type FleetResourceEc2InstanceCapabilitiesAcceleratorCapabilitiesModel struct {
	Selections types.ListType `tfsdk:"selections"`
	Count      types.Int32    `tfsdk:"count"`
}

// FleetResourceModel describes the resource data model.
type FleetResourceModel struct {
	DisplayName    types.String                     `tfsdk:"display_name"`
	Description    types.String                     `tfsdk:"description"`
	FarmId         types.String                     `tfsdk:"farm_id"`
	MinWorkerCount types.Int32                      `tfsdk:"min_worker_count"`
	MaxWorkerCount types.Int32                      `tfsdk:"max_worker_count"`
	RoleArn        types.String                     `tfsdk:"role_arn"`
	ID             types.String                     `tfsdk:"id"`
	Configuration  *FleetResourceConfigurationModel `tfsdk:"configuration"`
}

func (r *FleetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fleet"
}

func (r *FleetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fleet resource",
		Blocks: map[string]schema.Block{
			"configuration": schema.SingleNestedBlock{
				Blocks: map[string]schema.Block{
					"ec2_instance_capabilities": schema.SingleNestedBlock{
						Description: "The capabilities of the EC2 instance. Only required when the mode is 'aws_managed'.",
						Blocks: map[string]schema.Block{
							"accelerator_capabilities": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"selections": schema.ListNestedAttribute{
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{Required: true},
												"runtime": schema.StringAttribute{
													Optional: true,
												},
											},
										},
										Optional: true,
									},
									"count": schema.Int32Attribute{
										Optional:    true,
										Description: "The minimum number of accelerators that can be attached to the instance.If you set the value to 0, a worker will still have 1 GPU.",
									},
								},
							},
							"memory_mib_range": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"max": schema.Int32Attribute{
										Optional:    true,
										Description: "The number of IOPS for the root EBS volume. Only required when the mode is 'aws_managed'.",
									},
									"min": schema.Int32Attribute{
										Optional:    true,
										Description: "The size of the root EBS volume in GiB.",
									},
								},
							},
							"root_ebs_volume": schema.SingleNestedBlock{
								Attributes: map[string]schema.Attribute{
									"iops": schema.Int32Attribute{
										Optional:    true,
										Description: "The number of IOPS for the root EBS volume. Only required when the mode is 'aws_managed'.",
									},
									"size": schema.Int32Attribute{
										Optional:    true,
										Description: "The size of the root EBS volume in GiB.",
									},
									"throughput": schema.Int32Attribute{
										Optional:    true,
										Description: "The throughput of the root EBS volume in MiB/s.",
									},
								},
							},
						},
						Attributes: map[string]schema.Attribute{
							"cpu_architecture": schema.StringAttribute{
								Optional: true,
							},
							"min_cpu_count": schema.Int32Attribute{Optional: true},
							"max_cpu_count": schema.Int32Attribute{Optional: true},

							"os_family": schema.StringAttribute{
								Optional: true,
							},
							"allowed_instance_types": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
							},
							"exclude_instance_types": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
							},
						},
					},
				},
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						Optional:    true,
						Description: "The mode of the fleet configuration. It can either be 'aws_managed' or 'customer_managed'.",
					},
					"ec2_market_type": schema.StringAttribute{
						Optional:    true,
						Description: "The market type of the EC2 instance. It can either be 'spot' or 'on-demand'. Only required when the mode is 'aws_managed'.",
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the fleet.",
				Required:            true,
			},
			"role_arn": schema.StringAttribute{
				MarkdownDescription: "The ARN of the role that the fleet assumes.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the fleet.",
				Optional:            true,
			},
			"min_worker_count": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "The minimum number of workers that can be started in the fleet.",
			},
			"max_worker_count": schema.Int32Attribute{
				Required:            true,
				MarkdownDescription: "The maximum number of workers that can be started in the fleet.",
			},
			"farm_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the farm.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the fleet.",
			},
		},
	}
}

func (r *FleetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FleetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FleetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var configurationType dltypes.FleetConfiguration
	if data.Configuration != nil {
		if data.Configuration.Mode.ValueString() == "customer_managed" {
			configurationType = &dltypes.FleetConfigurationMemberCustomerManaged{
				Value: dltypes.CustomerManagedFleetConfiguration{
					Mode: dltypes.AutoScalingModeEventBasedAutoScaling,
				},
			}
		} else {
			archType := dltypes.CpuArchitectureTypeX8664
			archTypeSelector := dltypes.CpuArchitectureType(data.Configuration.Ec2InstanceCapabilities.CpuArchitecture.ValueString())
			if archTypeSelector == "arm64" {
				archType = dltypes.CpuArchitectureTypeArm64
			}
			osFamily := dltypes.ServiceManagedFleetOperatingSystemFamilyWindows
			if data.Configuration.Ec2InstanceCapabilities.OsFamily.ValueString() == "linux" {
				osFamily = dltypes.ServiceManagedFleetOperatingSystemFamilyLinux
			} else if data.Configuration.Ec2InstanceCapabilities.OsFamily.ValueString() == "windows" {
				osFamily = dltypes.ServiceManagedFleetOperatingSystemFamilyWindows
			}
			marketType := dltypes.Ec2MarketTypeOnDemand
			if data.Configuration.Ec2MarketType.ValueString() == "spot" {
				marketType = dltypes.Ec2MarketTypeSpot
			} else if data.Configuration.Ec2MarketType.ValueString() == "on-demand" {
				marketType = dltypes.Ec2MarketTypeOnDemand
			}
			configurationType = &dltypes.FleetConfigurationMemberServiceManagedEc2{
				Value: dltypes.ServiceManagedEc2FleetConfiguration{
					InstanceCapabilities: &dltypes.ServiceManagedEc2InstanceCapabilities{
						CpuArchitectureType: archType,
						OsFamily:            osFamily,
						MemoryMiB: &dltypes.MemoryMiBRange{
							Min: data.Configuration.Ec2InstanceCapabilities.MemoryMibRange.Min.ValueInt32Pointer(),
							Max: data.Configuration.Ec2InstanceCapabilities.MemoryMibRange.Max.ValueInt32Pointer(),
						},
						VCpuCount: &dltypes.VCpuCountRange{
							Min: data.Configuration.Ec2InstanceCapabilities.MinCpuCount.ValueInt32Pointer(),
							Max: data.Configuration.Ec2InstanceCapabilities.MaxCpuCount.ValueInt32Pointer(),
						},
					},
					InstanceMarketOptions: &dltypes.ServiceManagedEc2InstanceMarketOptions{
						Type: marketType,
					},
				},
			}
		}
	} else {
		resp.Diagnostics.AddError("Client Error", "Configuration is required")
		return
	}
	createRequest := deadline.CreateFleetInput{
		FarmId:         data.FarmId.ValueStringPointer(),
		MinWorkerCount: data.MinWorkerCount.ValueInt32(),
		MaxWorkerCount: data.MaxWorkerCount.ValueInt32Pointer(),
		DisplayName:    data.DisplayName.ValueStringPointer(),
		Description:    data.Description.ValueStringPointer(),
		RoleArn:        data.RoleArn.ValueStringPointer(),
		Configuration:  configurationType,
	}
	createOutputRaw, err := r.client.CreateFleet(ctx, &createRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, %s, got error: %s", r.typeName(), data.DisplayName.String(), err))
		return
	}
	if createOutputRaw == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, %s, got error: %s", r.typeName(), data.DisplayName.String(), "empty response"))
		return
	}
	createOutput := *createOutputRaw
	if createOutput.FleetId == nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create %s, %s, got error: %s", r.typeName(), data.DisplayName.String(), "empty Fleet"))
		return
	}
	data.ID = types.StringValue(*createOutput.FleetId)
	tflog.Trace(ctx, fmt.Sprintf("created %s", r.typeName()))
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FleetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FleetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FleetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FleetResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	request := &deadline.UpdateFleetInput{
		FarmId:      data.FarmId.ValueStringPointer(),
		FleetId:     data.ID.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		DisplayName: data.DisplayName.ValueStringPointer(),
	}
	var configurationType dltypes.FleetConfiguration
	if data.Configuration.Mode.ValueString() == "customer_managed" {
		configurationType = &dltypes.FleetConfigurationMemberCustomerManaged{
			Value: dltypes.CustomerManagedFleetConfiguration{
				Mode: dltypes.AutoScalingModeEventBasedAutoScaling,
			},
		}
	} else {
		archType := dltypes.CpuArchitectureTypeX8664
		archTypeSelector := dltypes.CpuArchitectureType(data.Configuration.Ec2InstanceCapabilities.CpuArchitecture.ValueString())
		if archTypeSelector == "arm64" {
			archType = dltypes.CpuArchitectureTypeArm64
		}
		configurationType = &dltypes.FleetConfigurationMemberServiceManagedEc2{
			Value: dltypes.ServiceManagedEc2FleetConfiguration{
				InstanceCapabilities: &dltypes.ServiceManagedEc2InstanceCapabilities{
					CpuArchitectureType: archType,
					VCpuCount: &dltypes.VCpuCountRange{
						Min: data.Configuration.Ec2InstanceCapabilities.MinCpuCount.ValueInt32Pointer(),
						Max: data.Configuration.Ec2InstanceCapabilities.MaxCpuCount.ValueInt32Pointer(),
					},
				},
				InstanceMarketOptions: &dltypes.ServiceManagedEc2InstanceMarketOptions{},
			},
		}
	}
	request.Configuration = configurationType
	_, err := r.client.UpdateFleet(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update %s, got error: %s", r.typeName(), err))
		return
	}
	data.Description = types.StringValue(*request.Description)
	data.DisplayName = types.StringValue(*request.DisplayName)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FleetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FleetResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceRequest := &deadline.DeleteFleetInput{
		FarmId:  data.FarmId.ValueStringPointer(),
		FleetId: data.ID.ValueStringPointer(),
	}
	_, err := r.client.DeleteFleet(ctx, deleteResourceRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete %s, got error: %s", r.typeName(), err))
		return
	}
}

func (r *FleetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *FleetResource) typeName() string {
	return "deadline_fleet"
}
