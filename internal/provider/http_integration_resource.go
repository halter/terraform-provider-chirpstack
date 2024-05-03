// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/halter-corp/terraform-provider-chirpstack/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &HttpIntegrationResource{}
var _ resource.ResourceWithImportState = &HttpIntegrationResource{}

func NewHttpIntegrationResource() resource.Resource {
	return &HttpIntegrationResource{}
}

// HttpIntegrationResource defines the resource implementation.
type HttpIntegrationResource struct {
	chirpstack client.Chirpstack
}

// HttpIntegrationResourceModel describes the resource data model.
type HttpIntegrationResourceModel struct {
	Id               types.String `tfsdk:"id"`
	ApplicationId    types.String `tfsdk:"application_id"`
	Headers          types.Map    `tfsdk:"headers"`
	Encoding         types.String `tfsdk:"encoding"`
	EventEndpointUrl types.String `tfsdk:"event_endpoint_url"`
}

func (r *HttpIntegrationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_integration"
}

func (r *HttpIntegrationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Http Integration resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Http Integration identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: "Application ID",
				Required:            true,
			},
			"encoding": schema.StringAttribute{
				MarkdownDescription: "Http Integration encoding. JSON or PROTOBUF.",
				Required:            true,
			},
			"event_endpoint_url": schema.StringAttribute{
				MarkdownDescription: "Http Integration URL",
				Required:            true,
			},
			"headers": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Http Integration headers",
				Optional:            true,
			},
		},
	}
}

func (r *HttpIntegrationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	chirpstack, ok := req.ProviderData.(client.Chirpstack)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.chirpstack = chirpstack
}

func httpIntegrationFromData(data *HttpIntegrationResourceModel) *api.HttpIntegration {
	httpIntegration := &api.HttpIntegration{
		ApplicationId:    data.ApplicationId.ValueString(),
		Encoding:         api.Encoding(api.Encoding_value[data.Encoding.ValueString()]),
		EventEndpointUrl: data.EventEndpointUrl.ValueString(),
	}

	if !data.Headers.IsNull() {
		httpIntegration.Headers = map[string]string{}
		for k, v := range data.Headers.Elements() {
			httpIntegration.Headers[k] = v.String()
		}
	}

	return httpIntegration
}

func httpIntegrationToData(httpIntegration *api.HttpIntegration, data *HttpIntegrationResourceModel) {
	data.ApplicationId = types.StringValue(httpIntegration.ApplicationId)
	data.Encoding = types.StringValue(httpIntegration.Encoding.String())
	data.EventEndpointUrl = types.StringValue(httpIntegration.EventEndpointUrl)
	if len(httpIntegration.Headers) > 0 {
		headers := map[string]attr.Value{}
		for k, v := range httpIntegration.Headers {
			headers[k] = types.StringValue(v)
		}
		data.Headers = types.MapValueMust(types.StringType, headers)
	}
}

func (r *HttpIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HttpIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create httpintegration, got error: %s", err))
	//     return
	// }
	httpIntegration := httpIntegrationFromData(&data)
	fmt.Printf("create: %+v\n", httpIntegration)
	err := r.chirpstack.CreateHttpIntegration(ctx, httpIntegration)
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to create httpintegration, got error: %s", err))
		return
	}

	// For the purposes of this httpintegration code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(httpIntegration.ApplicationId)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HttpIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data HttpIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read httpintegration, got error: %s", err))
	//     return
	// }
	httpIntegration, err := r.chirpstack.GetHttpIntegration(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to read httpintegration, got error: %s", err))
		return
	}

	httpIntegrationToData(httpIntegration, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HttpIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data HttpIntegrationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update httpintegration, got error: %s", err))
	//     return
	// }
	httpIntegration := httpIntegrationFromData(&data)
	err := r.chirpstack.UpdateHttpIntegration(ctx, httpIntegration)
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to update httpintegration, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HttpIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data HttpIntegrationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete httpintegration, got error: %s", err))
	//     return
	// }
	err := r.chirpstack.DeleteHttpIntegration(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to delete httpintegration, got error: %s", err))
		return
	}
}

func (r *HttpIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
