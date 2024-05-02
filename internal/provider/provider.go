// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/halter-corp/terraform-provider-chirpstack/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ChirpstackProvider satisfies various provider interfaces.
var _ provider.Provider = &ChirpstackProvider{}
var _ provider.ProviderWithFunctions = &ChirpstackProvider{}

// ChirpstackProvider defines the provider implementation.
type ChirpstackProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ChirpstackProviderModel describes the provider data model.
type ChirpstackProviderModel struct {
	Host types.String `tfsdk:"host"`
	Port types.Int64  `tfsdk:"port"`
	Key  types.String `tfsdk:"key"`
}

func (p *ChirpstackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "chirpstack"
	resp.Version = p.version
}

func (p *ChirpstackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Chirpstack hostname",
				Optional:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Chirpstack port",
				Optional:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Chirpstack api key",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ChirpstackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ChirpstackProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	defaultClient := http.DefaultClient
	resp.DataSourceData = defaultClient

	host := data.Host.ValueString()
	if host == "" {
		host = os.Getenv("CHIRPSTACK_HOST")
	}
	port := int(data.Port.ValueInt64())
	if port == 0 {
		port, _ = strconv.Atoi(os.Getenv("CHIRPSTACK_PORT"))
	}
	key := data.Key.ValueString()
	if key == "" {
		key = os.Getenv("CHIRPSTACK_KEY")
	}

	conn, err := client.GetChirpstackConn(ctx, host, port, key)
	if err != nil {
		resp.Diagnostics.AddError("could not establish chirpstack connection", err.Error())
		return
	}

	resp.ResourceData = client.NewChirpstack(conn)
}

func (p *ChirpstackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
		NewTenantResource,
		NewApplicationResource,
		NewDeviceProfileResource,
	}
}

func (p *ChirpstackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func (p *ChirpstackProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ChirpstackProvider{
			version: version,
		}
	}
}
