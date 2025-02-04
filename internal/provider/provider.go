// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/deadline"
	associatemembertofarm "github.com/enable-la/terraform-provider-aws-deadline/internal/resources/associate-member-to-farm"
	associatemembertofleet "github.com/enable-la/terraform-provider-aws-deadline/internal/resources/associate-member-to-fleet"
	associatequeuetofleet "github.com/enable-la/terraform-provider-aws-deadline/internal/resources/associate-queue-to-fleet"
	"github.com/enable-la/terraform-provider-aws-deadline/internal/resources/farm"
	"github.com/enable-la/terraform-provider-aws-deadline/internal/resources/fleet"
	licenseendpoint "github.com/enable-la/terraform-provider-aws-deadline/internal/resources/license-endpoint"
	"github.com/enable-la/terraform-provider-aws-deadline/internal/resources/queue"
	queueenvironment "github.com/enable-la/terraform-provider-aws-deadline/internal/resources/queue-environment"
	storageprofile "github.com/enable-la/terraform-provider-aws-deadline/internal/resources/storage-profile"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"log"
)

// Ensure AWSDeadlineProvider satisfies various provider interfaces.
var _ provider.Provider = &AWSDeadlineProvider{}
var _ provider.ProviderWithFunctions = &AWSDeadlineProvider{}

// AWSDeadlineProvider defines the provider implementation.
type AWSDeadlineProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AWSDeadlineProviderModel describes the provider data model.
type AWSDeadlineProviderModel struct {
	Region types.String `tfsdk:"region"`
}

func (p *AWSDeadlineProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "deadline"
	resp.Version = p.version
}

func (p *AWSDeadlineProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Description: "The AWS region to use for the Deadline API.",
				Optional:    true,
			},
		},
	}
}

func (p *AWSDeadlineProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AWSDeadlineProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	if data.Region.ValueString() != "" {
		cfg.Region = data.Region.ValueString()
	}
	svc := deadline.NewFromConfig(cfg)
	resp.DataSourceData = svc
	resp.ResourceData = svc
}

func (p *AWSDeadlineProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		farm.New,
		fleet.New,
		queue.New,
		queueenvironment.New,
		associatemembertofleet.New,
		associatemembertofarm.New,
		associatequeuetofleet.New,
		storageprofile.New,
		licenseendpoint.New,
	}
}

func (p *AWSDeadlineProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *AWSDeadlineProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AWSDeadlineProvider{
			version: version,
		}
	}
}
