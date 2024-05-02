package client

import (
	"context"
	"fmt"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/halter-corp/terraform-provider-chirpstack/client/model"
)

func (c *chirpstack) GetDevice(ctx context.Context, deviceEui string) (*model.GetDeviceResponse, error) {
	result := model.GetDeviceResponse{}
	getResp, err := c.deviceServiceClient.Get(ctx, &api.GetDeviceRequest{
		DevEui: deviceEui,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get device from chirpstack; device: %+v; err: %+v;", getResp, err)
	}
	result.Device = getResp.Device
	result.DeviceStatus = getResp.DeviceStatus
	result.ClassEnabled = getResp.ClassEnabled
	getKeysResp, err := c.deviceServiceClient.GetKeys(ctx, &api.GetDeviceKeysRequest{
		DevEui: deviceEui,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get keys for device chirpstack; keys %+v; err: %+v;", getKeysResp, err)
	}
	result.DeviceKeys = getKeysResp.DeviceKeys
	getActivitionResp, err := c.deviceServiceClient.GetActivation(ctx, &api.GetDeviceActivationRequest{
		DevEui: deviceEui,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get activation chirpstack; activation %+v; err: %+v;", getActivitionResp, err)
	}
	result.DeviceActivation = getActivitionResp.DeviceActivation
	return &result, nil
}

func (c *chirpstack) CreateDevice(ctx context.Context, applicationID, deviceProfileID, deviceEui, name, joinEui, devAddr, appSKey, nwkSEncKey, appKey string) error {
	device := api.Device{
		DevEui:          deviceEui,
		Name:            name,
		JoinEui:         joinEui,
		ApplicationId:   applicationID,
		DeviceProfileId: deviceProfileID,
	}
	_, err := c.deviceServiceClient.Create(ctx, &api.CreateDeviceRequest{
		Device: &device,
	})
	if err != nil {
		return fmt.Errorf("failed to create device in chirpstack; device: %s; err: %+v;", device.String(), err)
	}
	keys := api.DeviceKeys{
		DevEui: deviceEui,
		NwkKey: appKey,
	}
	_, err = c.deviceServiceClient.CreateKeys(ctx, &api.CreateDeviceKeysRequest{
		DeviceKeys: &keys,
	})

	if err != nil {
		return fmt.Errorf("failed to create keys for device chirpstack; keys %s; err: %+v;", keys.String(), err)
	}
	activation := api.DeviceActivation{
		DevEui:      deviceEui,
		DevAddr:     devAddr,
		AppSKey:     appSKey,
		NwkSEncKey:  nwkSEncKey,
		SNwkSIntKey: nwkSEncKey,
		FNwkSIntKey: nwkSEncKey,
	}
	_, err = c.deviceServiceClient.Activate(ctx, &api.ActivateDeviceRequest{
		DeviceActivation: &activation,
	})

	if err != nil {
		return fmt.Errorf("failed to activate device chirpstack; activation %s; err: %v;", activation.String(), err)
	}
	return nil
}

func (c *chirpstack) ListDevices(ctx context.Context, applicationID, name string, limit uint32) ([]*api.DeviceListItem, error) {
	resp, err := c.deviceServiceClient.List(ctx, &api.ListDevicesRequest{
		ApplicationId: applicationID,
		Search:        name,
		Limit:         limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list devices from chirpstack: %v", err)
	}
	return resp.Result, nil
}

func (c *chirpstack) DeleteDevice(ctx context.Context, deviceEui string) error {
	_, err := c.deviceServiceClient.Delete(ctx, &api.DeleteDeviceRequest{
		DevEui: deviceEui,
	})
	return err
}
