package client

import (
	"context"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
)

func (c *chirpstack) Enqueue(ctx context.Context, request *api.EnqueueDeviceQueueItemRequest) (*api.EnqueueDeviceQueueItemResponse, error) {
	return c.deviceServiceClient.Enqueue(ctx, request)
}

func (c *chirpstack) MulticastEnqueue(ctx context.Context, request *api.EnqueueMulticastGroupQueueItemRequest) (*api.EnqueueMulticastGroupQueueItemResponse, error) {
	return c.multicastGroupServiceClient.Enqueue(ctx, request)
}
