package client

import (
	"context"
	"fmt"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
)

func (c *chirpstack) ListDeviceProfiles(ctx context.Context, tenantID, name string, limit uint32) ([]*api.DeviceProfileListItem, error) {
	resp, err := c.deviceProfileServiceClient.List(ctx, &api.ListDeviceProfilesRequest{
		TenantId: tenantID,
		Search:   name,
		Limit:    limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list device profiles from chirpstack: %v", err)
	}
	return resp.Result, nil
}

func (c *chirpstack) GetDeviceProfile(ctx context.Context, id string) (*api.DeviceProfile, error) {
	req := api.GetDeviceProfileRequest{
		Id: id,
	}
	resp, err := c.deviceProfileServiceClient.Get(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to get device profile %s; err: %s;", id, err)
	}
	return resp.DeviceProfile, nil
}

func (c *chirpstack) CreateDeviceProfile(ctx context.Context, deviceProfile *api.DeviceProfile) (string, error) {
	createDeviceProfilesRequest := api.CreateDeviceProfileRequest{
		DeviceProfile: deviceProfile,
	}
	listDeviceProfilesResponse, err := c.deviceProfileServiceClient.Create(ctx, &createDeviceProfilesRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create device profile: %+v; err: %s;", deviceProfile, err)
	}
	return listDeviceProfilesResponse.Id, nil
}

func (c *chirpstack) UpdateDeviceProfile(ctx context.Context, deviceProfile *api.DeviceProfile) error {
	updateDeviceProfilesRequest := api.UpdateDeviceProfileRequest{
		DeviceProfile: deviceProfile,
	}
	_, err := c.deviceProfileServiceClient.Update(ctx, &updateDeviceProfilesRequest)
	if err != nil {
		return fmt.Errorf("failed to update device profile %s; err: %s;", deviceProfile.Id, err)
	}
	return nil
}

func (c *chirpstack) DeleteDeviceProfile(ctx context.Context, id string) error {
	deleteDeviceProfilesRequest := api.DeleteDeviceProfileRequest{
		Id: id,
	}
	_, err := c.deviceProfileServiceClient.Delete(ctx, &deleteDeviceProfilesRequest)
	if err != nil {
		return fmt.Errorf("failed to delete device profile id %s; err: %s;", id, err)
	}
	return nil
}
