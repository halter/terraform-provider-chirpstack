package client

import (
	"context"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/chirpstack/chirpstack/api/go/v4/common"
)

func (c *chirpstack) CreateGateway(ctx context.Context, gatewayEui, tenantID string, latitude, longitude, altitude float64, accuracy float32, statsInterval uint32) error {
	gw := api.Gateway{
		GatewayId:   gatewayEui,
		Name:        "eui-" + gatewayEui,
		Description: "eui-" + gatewayEui,
		TenantId:    tenantID,
		Location: &common.Location{
			Latitude:  latitude,
			Longitude: longitude,
			Altitude:  altitude,
			Accuracy:  accuracy,
		},
		StatsInterval: statsInterval,
	}

	_, err := c.gatewayServiceClient.Create(ctx, &api.CreateGatewayRequest{
		Gateway: &gw,
	})
	if err != nil {
		_, err = c.gatewayServiceClient.Update(ctx, &api.UpdateGatewayRequest{
			Gateway: &gw,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *chirpstack) ListGateways(ctx context.Context, request *api.ListGatewaysRequest) ([]*api.GatewayListItem, error) {
	resp, err := c.gatewayServiceClient.List(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.Result, nil
}
