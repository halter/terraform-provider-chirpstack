package client

import (
	"context"
	"fmt"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
)

func (c *chirpstack) ListApplications(ctx context.Context, tenantID, name string, limit uint32) ([]*api.ApplicationListItem, error) {
	listApplicationsRequest := api.ListApplicationsRequest{
		TenantId: tenantID,
		Search:   name,
		Limit:    limit,
	}
	listApplicationsResponse, listErr := c.applicationServiceClient.List(ctx, &listApplicationsRequest)
	if listErr != nil {
		return nil, fmt.Errorf("failed to list applications; err: %s;", listErr)
	}
	return listApplicationsResponse.Result, nil
}

func (c *chirpstack) GetApplication(ctx context.Context, id string) (*api.Application, error) {
	req := api.GetApplicationRequest{
		Id: id,
	}
	resp, err := c.applicationServiceClient.Get(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to get application %s; err: %s;", id, err)
	}
	return resp.Application, nil
}

func (c *chirpstack) CreateApplication(ctx context.Context, tenantId, name, description string) (string, error) {
	createApplicationsRequest := api.CreateApplicationRequest{
		Application: &api.Application{
			TenantId:    tenantId,
			Name:        name,
			Description: description,
		},
	}
	listApplicationsResponse, err := c.applicationServiceClient.Create(ctx, &createApplicationsRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create application %s; err: %s;", name, err)
	}
	return listApplicationsResponse.Id, nil
}

func (c *chirpstack) UpdateApplication(ctx context.Context, application *api.Application) error {
	updateApplicationsRequest := api.UpdateApplicationRequest{
		Application: application,
	}
	_, err := c.applicationServiceClient.Update(ctx, &updateApplicationsRequest)
	if err != nil {
		return fmt.Errorf("failed to update application %s; err: %s;", application.Id, err)
	}
	return nil
}

func (c *chirpstack) DeleteApplication(ctx context.Context, id string) error {
	deleteApplicationsRequest := api.DeleteApplicationRequest{
		Id: id,
	}
	_, err := c.applicationServiceClient.Delete(ctx, &deleteApplicationsRequest)
	if err != nil {
		return fmt.Errorf("failed to delete application id %s; err: %s;", id, err)
	}
	return nil
}
