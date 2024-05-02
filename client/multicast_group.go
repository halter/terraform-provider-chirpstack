package client

import (
	"context"
	"fmt"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/chirpstack/chirpstack/api/go/v4/common"
)

func (c *chirpstack) ListMulticastGroups(ctx context.Context, applicationID, name string, limit uint32) ([]*api.MulticastGroupListItem, error) {
	resp, err := c.multicastGroupServiceClient.List(ctx, &api.ListMulticastGroupsRequest{
		ApplicationId: applicationID,
		Search:        name,
		Limit:         limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list multicast groups from chirpstack: %v", err)
	}
	return resp.Result, nil
}

func (c *chirpstack) GetMulticastGroup(ctx context.Context, id string) (*api.GetMulticastGroupResponse, error) {
	resp, err := c.multicastGroupServiceClient.Get(ctx, &api.GetMulticastGroupRequest{
		Id: id,
	})
	return resp, err
}

func (c *chirpstack) CreateMulticastGroup(ctx context.Context, applicationID, name string, region common.Region, mcAddr, mcNwkSKey, mcAppSKey string, fCnt uint32, dr, frequency uint32) error {
	device := api.MulticastGroup{
		Name:                 name,
		ApplicationId:        applicationID,
		Region:               region,
		McAddr:               mcAddr,
		McNwkSKey:            mcNwkSKey,
		McAppSKey:            mcAppSKey,
		FCnt:                 fCnt,
		GroupType:            api.MulticastGroupType_CLASS_C,
		Dr:                   dr,
		Frequency:            frequency,
		ClassCSchedulingType: api.MulticastGroupSchedulingType_GPS_TIME,
	}
	_, err := c.multicastGroupServiceClient.Create(ctx, &api.CreateMulticastGroupRequest{
		MulticastGroup: &device,
	})
	if err != nil {
		return fmt.Errorf("failed to create multicast group in chirpstack; device: %s; err: %+v;", device.String(), err)
	}
	return nil
}

func (c *chirpstack) DeleteMulticastGroup(ctx context.Context, id string) error {
	_, err := c.multicastGroupServiceClient.Delete(ctx, &api.DeleteMulticastGroupRequest{
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("failed to delete multicast group in chirpstack; id: %s; err: %+v;", id, err)
	}
	return nil
}

func (c *chirpstack) AddGatewayToMulticastGroup(ctx context.Context, multicastGroupId, gatewayId string) error {
	_, err := c.multicastGroupServiceClient.AddGateway(ctx, &api.AddGatewayToMulticastGroupRequest{
		MulticastGroupId: multicastGroupId,
		GatewayId:        gatewayId,
	})
	if err != nil {
		return fmt.Errorf("failed to add gateway to multicast group in chirpstack; multicast group ID: %s; gatewayId: %s err: %+v;", multicastGroupId, gatewayId, err)
	}
	return nil
}

func (c *chirpstack) RemoveGatewayFromMulticastGroup(ctx context.Context, multicastGroupId, gatewayId string) error {
	_, err := c.multicastGroupServiceClient.RemoveGateway(ctx, &api.RemoveGatewayFromMulticastGroupRequest{
		MulticastGroupId: multicastGroupId,
		GatewayId:        gatewayId,
	})
	if err != nil {
		return fmt.Errorf("failed to remove gateway from multicast group in chirpstack; multicast group ID: %s; gatewayId: %s err: %+v;", multicastGroupId, gatewayId, err)
	}
	return nil
}
