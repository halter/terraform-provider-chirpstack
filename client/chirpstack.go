package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/chirpstack/chirpstack/api/go/v4/api"
	"github.com/chirpstack/chirpstack/api/go/v4/common"
	"github.com/halter-corp/terraform-provider-chirpstack/client/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Chirpstack interface {
	// tenant
	CreateTenant(ctx context.Context, name, description string) (string, error)
	GetTenant(ctx context.Context, id string) (*api.Tenant, error)
	UpdateTenant(ctx context.Context, tenant *api.Tenant) error
	DeleteTenant(ctx context.Context, id string) error
	ListTenants(ctx context.Context, name string, limit uint32) ([]*api.TenantListItem, error)

	// application
	ListApplications(ctx context.Context, tenantID, name string, limit uint32) ([]*api.ApplicationListItem, error)
	CreateApplication(ctx context.Context, tenantID, name, description string) (string, error)
	GetApplication(ctx context.Context, id string) (*api.Application, error)
	UpdateApplication(ctx context.Context, application *api.Application) error
	DeleteApplication(ctx context.Context, id string) error

	// integrations
	CreateHttpIntegration(ctx context.Context, integration *api.HttpIntegration) error
	GetHttpIntegration(ctx context.Context, applicationId string) (*api.HttpIntegration, error)
	UpdateHttpIntegration(ctx context.Context, integration *api.HttpIntegration) error
	DeleteHttpIntegration(ctx context.Context, applicationId string) error

	// messaging
	Enqueue(ctx context.Context, request *api.EnqueueDeviceQueueItemRequest) (*api.EnqueueDeviceQueueItemResponse, error)
	MulticastEnqueue(ctx context.Context, request *api.EnqueueMulticastGroupQueueItemRequest) (*api.EnqueueMulticastGroupQueueItemResponse, error)

	// multicast group
	ListMulticastGroups(ctx context.Context, applicationID, name string, limit uint32) ([]*api.MulticastGroupListItem, error)
	GetMulticastGroup(ctx context.Context, id string) (*api.GetMulticastGroupResponse, error)
	CreateMulticastGroup(ctx context.Context, applicationID, name string, region common.Region, mcAddr, mcNwkSKey, mcAppSKey string, fCnt uint32, dr, frequency uint32) error
	DeleteMulticastGroup(ctx context.Context, id string) error
	AddGatewayToMulticastGroup(ctx context.Context, multicastGroupId, gatewayId string) error
	RemoveGatewayFromMulticastGroup(ctx context.Context, multicastGroupId, gatewayId string) error

	// gateway
	ListGateways(ctx context.Context, request *api.ListGatewaysRequest) ([]*api.GatewayListItem, error)
	CreateGateway(ctx context.Context, gatewayEui, tenantID string, latitude, longitude, altitude float64, accuracy float32, statsInterval uint32) error

	// device
	ListDevices(ctx context.Context, applicationID, name string, limit uint32) ([]*api.DeviceListItem, error)
	GetDevice(ctx context.Context, deviceEui string) (*model.GetDeviceResponse, error)
	CreateDevice(ctx context.Context, applicationID, deviceProfileID, deviceEui, name, joinEui, devAddr, appSKey, nwkSEncKey, appKey string) error
	DeleteDevice(ctx context.Context, deviceEui string) error

	// device profile
	ListDeviceProfiles(ctx context.Context, tenantID, name string, limit uint32) ([]*api.DeviceProfileListItem, error)
	GetDeviceProfile(ctx context.Context, id string) (*api.DeviceProfile, error)
	CreateDeviceProfile(ctx context.Context, deviceProfile *api.DeviceProfile) (string, error)
	UpdateDeviceProfile(ctx context.Context, deviceProfile *api.DeviceProfile) error
	DeleteDeviceProfile(ctx context.Context, id string) error
}

type apiToken string

func (a apiToken) GetRequestMetadata(
	ctx context.Context,
	url ...string,
) (map[string]string, error) {
	return map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", a),
	}, nil
}

func (a apiToken) RequireTransportSecurity() bool {
	return false
}

func getCerts(host string, port int) ([]*x509.Certificate, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), conf)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close()
	}()
	certs := conn.ConnectionState().PeerCertificates
	return certs, nil
}

func loadTLSCredentials(host string, port int) (credentials.TransportCredentials, error) {
	certs, err := getCerts(host, port)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	for _, cert := range certs {
		certPool.AddCert(cert)
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func GetChirpstackConn(ctx context.Context, host string, port int, apiKey string) (grpc.ClientConnInterface, error) {
	tlsCredentials, loadTLSCredErr := loadTLSCredentials(host, port)
	if loadTLSCredErr != nil {
		return nil, fmt.Errorf("cannot load TLS credentials: %v", loadTLSCredErr)
	}

	// debug issues with: export GRPC_GO_LOG_SEVERITY_LEVEL=info
	optionList := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(apiToken(apiKey)),
		grpc.WithTransportCredentials(tlsCredentials),
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	conn, dialErr := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), optionList...)
	return conn, dialErr
}

type chirpstack struct {
	tenantServiceClient         api.TenantServiceClient
	applicationServiceClient    api.ApplicationServiceClient
	deviceServiceClient         api.DeviceServiceClient
	deviceProfileServiceClient  api.DeviceProfileServiceClient
	multicastGroupServiceClient api.MulticastGroupServiceClient
	gatewayServiceClient        api.GatewayServiceClient
}

func NewChirpstack(conn grpc.ClientConnInterface) Chirpstack {
	return &chirpstack{
		tenantServiceClient:         api.NewTenantServiceClient(conn),
		applicationServiceClient:    api.NewApplicationServiceClient(conn),
		deviceServiceClient:         api.NewDeviceServiceClient(conn),
		deviceProfileServiceClient:  api.NewDeviceProfileServiceClient(conn),
		multicastGroupServiceClient: api.NewMulticastGroupServiceClient(conn),
		gatewayServiceClient:        api.NewGatewayServiceClient(conn),
	}
}
