package client

import (
	"context"
	"fmt"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
)

func (c *chirpstack) ListTenants(ctx context.Context, name string, limit uint32) ([]*api.TenantListItem, error) {
	listTenantsRequest := api.ListTenantsRequest{
		Limit:  limit,
		Search: name,
	}
	listTenantsResponse, listErr := c.tenantServiceClient.List(ctx, &listTenantsRequest)
	if listErr != nil {
		return nil, fmt.Errorf("failed to list tenants; err: %s;", listErr)
	}
	return listTenantsResponse.Result, nil
}

func (c *chirpstack) GetTenant(ctx context.Context, id string) (*api.Tenant, error) {
	req := api.GetTenantRequest{
		Id: id,
	}
	resp, err := c.tenantServiceClient.Get(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant %s; err: %s;", id, err)
	}
	return resp.Tenant, nil
}

func (c *chirpstack) CreateTenant(ctx context.Context, name, description string) (string, error) {
	createTenantsRequest := api.CreateTenantRequest{
		Tenant: &api.Tenant{
			Name:        name,
			Description: description,
		},
	}
	listTenantsResponse, err := c.tenantServiceClient.Create(ctx, &createTenantsRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create tenant %s; err: %s;", name, err)
	}
	return listTenantsResponse.Id, nil
}

func (c *chirpstack) UpdateTenant(ctx context.Context, tenant *api.Tenant) error {
	updateTenantsRequest := api.UpdateTenantRequest{
		Tenant: tenant,
	}
	_, err := c.tenantServiceClient.Update(ctx, &updateTenantsRequest)
	if err != nil {
		return fmt.Errorf("failed to update tenant %s; err: %s;", tenant.Id, err)
	}
	return nil
}

func (c *chirpstack) DeleteTenant(ctx context.Context, id string) error {
	deleteTenantsRequest := api.DeleteTenantRequest{
		Id: id,
	}
	_, err := c.tenantServiceClient.Delete(ctx, &deleteTenantsRequest)
	if err != nil {
		return fmt.Errorf("failed to delete tenant id %s; err: %s;", id, err)
	}
	return nil
}
