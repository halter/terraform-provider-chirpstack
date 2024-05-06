// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/halter-corp/terraform-provider-chirpstack/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TenantResource{}
var _ resource.ResourceWithImportState = &TenantResource{}

func NewTenantResource() resource.Resource {
	return &TenantResource{}
}

// TenantResource defines the resource implementation.
type TenantResource struct {
	chirpstack client.Chirpstack
}

// TenantResourceModel describes the resource data model.
type TenantResourceModel struct {
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Id                  types.String `tfsdk:"id"`
	CanHaveGateways     types.Bool   `tfsdk:"can_have_gateways"`
	MaxGatewayCount     types.Int64  `tfsdk:"max_gateway_count"`
	MaxDeviceCount      types.Int64  `tfsdk:"max_device_count"`
	PrivateGatewaysUp   types.Bool   `tfsdk:"private_gateways_up"`
	PrivateGatewaysDown types.Bool   `tfsdk:"private_gateways_down"`
}

func (r *TenantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant"
}

func (r *TenantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Tenant resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Tenant identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Tenant name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Tenant description",
				Optional:            true,
			},
			"can_have_gateways": schema.BoolAttribute{
				MarkdownDescription: `Can the tenant create and "own" Gateways?`,
				Optional:            true,
				Computed:            true,
			},
			"max_gateway_count": schema.Int64Attribute{
				MarkdownDescription: "Max. gateway count for tenant. When set to 0, the tenant can have unlimited gateways.",
				Optional:            true,
				Computed:            true,
			},
			"max_device_count": schema.Int64Attribute{
				MarkdownDescription: "Max. device count for tenant. When set to 0, the tenant can have unlimited devices.",
				Optional:            true,
				Computed:            true,
			},
			"private_gateways_up": schema.BoolAttribute{
				MarkdownDescription: "Private gateways (uplink). If enabled, then uplink messages will not be shared with other tenants.",
				Optional:            true,
				Computed:            true,
			},
			"private_gateways_down": schema.BoolAttribute{
				MarkdownDescription: `Private gateways (downlink).
If enabled, then other tenants will not be able to schedule downlink
messages through the gateways of this tenant. For example, in case you
do want to share uplinks with other tenants (private_gateways_up=false),
but you want to prevent other tenants from using gateway airtime.`,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *TenantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func tenantFromData(data *TenantResourceModel) *api.Tenant {
	tenant := api.Tenant{
		Id:   data.Id.ValueString(),
		Name: data.Name.ValueString(),
	}
	if !data.Description.IsNull() {
		tenant.Description = data.Description.ValueString()
	}
	if !data.CanHaveGateways.IsNull() {
		tenant.CanHaveGateways = data.CanHaveGateways.ValueBool()
	}
	if !data.MaxGatewayCount.IsNull() {
		tenant.MaxGatewayCount = uint32(data.MaxGatewayCount.ValueInt64())
	}
	if !data.MaxDeviceCount.IsNull() {
		tenant.MaxDeviceCount = uint32(data.MaxDeviceCount.ValueInt64())
	}
	if !data.PrivateGatewaysUp.IsNull() {
		tenant.PrivateGatewaysUp = data.PrivateGatewaysUp.ValueBool()
	}
	if !data.PrivateGatewaysDown.IsNull() {
		tenant.PrivateGatewaysDown = data.PrivateGatewaysDown.ValueBool()
	}
	return &tenant
}
func tenantToData(tenant *api.Tenant, data *TenantResourceModel) {
	data.Name = types.StringValue(tenant.Name)
	if tenant.Description != "" {
		data.Description = types.StringValue(tenant.Description)
	}
	data.CanHaveGateways = types.BoolValue(tenant.CanHaveGateways)
	data.MaxGatewayCount = types.Int64Value(int64(tenant.MaxGatewayCount))
	data.MaxDeviceCount = types.Int64Value(int64(tenant.MaxDeviceCount))
	data.PrivateGatewaysUp = types.BoolValue(tenant.PrivateGatewaysUp)
	data.PrivateGatewaysDown = types.BoolValue(tenant.PrivateGatewaysDown)
}

func (r *TenantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TenantResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create tenant, got error: %s", err))
	//     return
	// }
	tenant := tenantFromData(&data)
	id, err := r.chirpstack.CreateTenant(ctx, tenant)
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to create tenant, got error: %s", err))
		return
	}

	// For the purposes of this tenant code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(id)
	tenantToData(tenant, &data)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TenantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TenantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read tenant, got error: %s", err))
	//     return
	// }
	tenant, err := r.chirpstack.GetTenant(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to read tenant, got error: %s", err))
		return
	}

	tenantToData(tenant, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TenantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TenantResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update tenant, got error: %s", err))
	//     return
	// }
	tenant := tenantFromData(&data)
	err := r.chirpstack.UpdateTenant(ctx, tenant)
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to update tenant, got error: %s", err))
		return
	}
	tenantToData(tenant, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TenantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TenantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete tenant, got error: %s", err))
	//     return
	// }
	err := r.chirpstack.DeleteTenant(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to delete tenant, got error: %s", err))
		return
	}
}

func (r *TenantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
