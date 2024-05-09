// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/chirpstack/chirpstack/api/go/v4/common"
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
var _ resource.Resource = &DeviceProfileResource{}
var _ resource.ResourceWithImportState = &DeviceProfileResource{}

func NewDeviceProfileResource() resource.Resource {
	return &DeviceProfileResource{}
}

// DeviceProfileResource defines the resource implementation.
type DeviceProfileResource struct {
	chirpstack client.Chirpstack
}

// DeviceProfileResourceModel describes the resource data model.
type DeviceProfileResourceModel struct {
	Id                           types.String `tfsdk:"id"`
	TenantId                     types.String `tfsdk:"tenant_id"`
	Name                         types.String `tfsdk:"name"`
	Description                  types.String `tfsdk:"description"`
	Region                       types.String `tfsdk:"region"`
	RegionConfigId               types.String `tfsdk:"region_config_id"`
	RegionParametersRevision     types.String `tfsdk:"region_parameters_revision"`
	MacVersion                   types.String `tfsdk:"mac_version"`
	AdrAlgorithm                 types.String `tfsdk:"adr_algorithm"`
	FlushQueueOnActivate         types.Bool   `tfsdk:"flush_queue_on_activate"`
	AllowRoaming                 types.Bool   `tfsdk:"allow_roaming"`
	ExpectedUplinkInterval       types.Int64  `tfsdk:"expected_uplink_interval"`
	DeviceStatusRequestFrequency types.Int64  `tfsdk:"device_status_request_frequency"`
	DeviceSupportsOTAA           types.Bool   `tfsdk:"device_supports_otaa"`
	DeviceSupportsClassB         types.Bool   `tfsdk:"device_supports_class_b"`
	DeviceSupportsClassC         types.Bool   `tfsdk:"device_supports_class_c"`
	ClassCTimeout                types.Int64  `tfsdk:"class_c_timeout"`
}

func (r *DeviceProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_profile"
}

func (r *DeviceProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DeviceProfile resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "DeviceProfile identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "Tenant ID",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Device profile name",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Device profile description",
				Optional:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Device profile region",
				Required:            true,
			},
			"region_config_id": schema.StringAttribute{
				MarkdownDescription: "Region configuration ID",
				Optional:            true,
			},
			"mac_version": schema.StringAttribute{
				MarkdownDescription: "The LoRaWAN MAC version supported by the device.",
				Required:            true,
			},
			"region_parameters_revision": schema.StringAttribute{
				MarkdownDescription: "Revision of the Regional Parameters specification supported by the device.",
				Required:            true,
			},
			"adr_algorithm": schema.StringAttribute{
				MarkdownDescription: "The ADR algorithm that will be used for controlling the device data-rate.",
				Optional:            true,
			},
			"flush_queue_on_activate": schema.BoolAttribute{
				MarkdownDescription: "The ADR algorithm that will be used for controlling the device data-rate.",
				Optional:            true,
			},
			"allow_roaming": schema.BoolAttribute{
				MarkdownDescription: "If enabled (and if roaming is configured on the server), this allows the device to use roaming.",
				Optional:            true,
			},
			"expected_uplink_interval": schema.Int64Attribute{
				MarkdownDescription: "The expected interval in seconds in which the device sends uplink messages. This is used to determine if a device is active or inactive.",
				Optional:            true,
			},
			"device_status_request_frequency": schema.Int64Attribute{
				MarkdownDescription: "Frequency to initiate an End-Device status request (request/day). Set to 0 to disable.",
				Optional:            true,
			},
			"device_supports_otaa": schema.BoolAttribute{
				MarkdownDescription: "Device supports OTAA",
				Optional:            true,
			},
			"device_supports_class_b": schema.BoolAttribute{
				MarkdownDescription: "Device supports Class-B",
				Optional:            true,
			},
			"device_supports_class_c": schema.BoolAttribute{
				MarkdownDescription: "Device supports Class-C",
				Optional:            true,
			},
			"class_c_timeout": schema.Int64Attribute{
				MarkdownDescription: "Class-C timeout (seconds). This is the maximum time ChirpStack will wait to receive an acknowledgement from the device (if requested).",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *DeviceProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func deviceProfileFromData(data *DeviceProfileResourceModel) *api.DeviceProfile {
	deviceProfile := &api.DeviceProfile{
		Id:       data.Id.ValueString(),
		TenantId: data.TenantId.ValueString(),
		Name:     data.Name.ValueString(),
	}

	if !data.Description.IsNull() {
		deviceProfile.Description = data.Description.ValueString()
	}
	if !data.Region.IsNull() {
		deviceProfile.Region = common.Region(common.Region_value[data.Region.ValueString()])
	}
	if !data.RegionConfigId.IsNull() {
		deviceProfile.RegionConfigId = data.RegionConfigId.ValueString()
	}
	if !data.RegionParametersRevision.IsNull() {
		deviceProfile.RegParamsRevision = common.RegParamsRevision(common.RegParamsRevision_value[data.RegionParametersRevision.ValueString()])
	}
	if !data.MacVersion.IsNull() {
		deviceProfile.MacVersion = common.MacVersion(common.MacVersion_value[data.MacVersion.ValueString()])
	}
	if !data.AdrAlgorithm.IsNull() {
		deviceProfile.AdrAlgorithmId = data.AdrAlgorithm.ValueString()
	}
	if !data.FlushQueueOnActivate.IsNull() {
		deviceProfile.FlushQueueOnActivate = data.FlushQueueOnActivate.ValueBool()
	}
	if !data.AllowRoaming.IsNull() {
		deviceProfile.AllowRoaming = data.AllowRoaming.ValueBool()
	}
	if !data.ExpectedUplinkInterval.IsNull() {
		deviceProfile.UplinkInterval = uint32(data.ExpectedUplinkInterval.ValueInt64())
	}
	if !data.DeviceStatusRequestFrequency.IsNull() {
		deviceProfile.DeviceStatusReqInterval = uint32(data.DeviceStatusRequestFrequency.ValueInt64())
	}
	if !data.DeviceSupportsOTAA.IsNull() {
		deviceProfile.SupportsOtaa = data.DeviceSupportsOTAA.ValueBool()
	}
	if !data.DeviceSupportsClassB.IsNull() {
		deviceProfile.SupportsClassB = data.DeviceSupportsClassB.ValueBool()
	}
	if !data.DeviceSupportsClassC.IsNull() {
		deviceProfile.SupportsClassC = data.DeviceSupportsClassC.ValueBool()
	}
	deviceProfile.ClassCTimeout = uint32(data.ClassCTimeout.ValueInt64())

	return deviceProfile
}

func deviceProfileToData(deviceProfile *api.DeviceProfile, data *DeviceProfileResourceModel) {
	data.TenantId = types.StringValue(deviceProfile.TenantId)
	data.Name = types.StringValue(deviceProfile.Name)
	if deviceProfile.Description != "" {
		data.Description = types.StringValue(deviceProfile.Description)
	}
	data.Region = types.StringValue(deviceProfile.Region.String())
	if deviceProfile.RegionConfigId != "" {
		data.RegionConfigId = types.StringValue(deviceProfile.RegionConfigId)
	}
	if deviceProfile.RegParamsRevision.String() != "" {
		data.RegionParametersRevision = types.StringValue(deviceProfile.RegParamsRevision.String())
	}
	data.MacVersion = types.StringValue(deviceProfile.MacVersion.String())
	if deviceProfile.AdrAlgorithmId != "" {
		data.AdrAlgorithm = types.StringValue(deviceProfile.AdrAlgorithmId)
	}
	data.FlushQueueOnActivate = types.BoolValue(deviceProfile.FlushQueueOnActivate)
	data.AllowRoaming = types.BoolValue(deviceProfile.AllowRoaming)
	data.ExpectedUplinkInterval = types.Int64Value(int64(deviceProfile.UplinkInterval))
	data.DeviceStatusRequestFrequency = types.Int64Value(int64(deviceProfile.DeviceStatusReqInterval))
	data.DeviceSupportsOTAA = types.BoolValue(deviceProfile.SupportsOtaa)
	data.DeviceSupportsClassB = types.BoolValue(deviceProfile.SupportsClassB)
	data.DeviceSupportsClassC = types.BoolValue(deviceProfile.SupportsClassC)
	data.ClassCTimeout = types.Int64Value(int64(deviceProfile.ClassCTimeout))
}

func (r *DeviceProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeviceProfileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create device profile, got error: %s", err))
	//     return
	// }
	deviceProfile := deviceProfileFromData(&data)

	id, err := r.chirpstack.CreateDeviceProfile(ctx, deviceProfile)
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to create device profile, got error: %s", err))
		return
	}

	// For the purposes of this device profile code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.StringValue(id)
	deviceProfileToData(deviceProfile, &data)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeviceProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read device profile, got error: %s", err))
	//     return
	// }
	deviceProfile, err := r.chirpstack.GetDeviceProfile(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to read device profile, got error: %s", err))
		return
	}

	deviceProfileToData(deviceProfile, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeviceProfileResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update device profile, got error: %s", err))
	//     return
	// }
	deviceProfile := deviceProfileFromData(&data)
	err := r.chirpstack.UpdateDeviceProfile(ctx, deviceProfile)
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to update device profile, got error: %s", err))
		return
	}
	deviceProfileToData(deviceProfile, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeviceProfileResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete device profile, got error: %s", err))
	//     return
	// }
	err := r.chirpstack.DeleteDeviceProfile(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Chirpstack Error", fmt.Sprintf("Unable to delete device profile, got error: %s", err))
		return
	}
}

func (r *DeviceProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
