package azure

import (
	"context"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-08-01/containerservice"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

var ErrNotAvilable = fmt.Errorf("managed cluster resource is not available")

type nilManagedClusterInformation struct{}

var _ ManagedClusterInformation = &nilManagedClusterInformation{}

func (m *nilManagedClusterInformation) IsAvailable() bool {
	return false
}

func (m *nilManagedClusterInformation) GetResourceID() string {
	return ""
}

func (m *nilManagedClusterInformation) GetSubscriptionID() string {
	return ""
}

func (m *nilManagedClusterInformation) GetRegion() string {
	return ""
}

func (m *nilManagedClusterInformation) GetResourceGroup() string {
	return ""
}

func (m *nilManagedClusterInformation) GetResourceName() string {
	return ""
}

func (m *nilManagedClusterInformation) GetNodeResourceGroup() string {
	return ""
}

func (m *nilManagedClusterInformation) GetLatestModel(ctx context.Context) (containerservice.ManagedCluster, error) {
	return containerservice.ManagedCluster{}, ErrNotAvilable
}

func (m *nilManagedClusterInformation) GetKubeConfig(ctx context.Context) (string, error) {
	return "", ErrNotAvilable
}

func (m *nilManagedClusterInformation) GetAgentPoolInformation(
	ctx context.Context, agentPoolName string,
) (ManagedClusterAgentPoolInformation, error) {
	return nil, ErrNotAvilable
}

type nilManagedClusterAgentPoolInformation struct{}

var _ ManagedClusterAgentPoolInformation = &nilManagedClusterAgentPoolInformation{}

func (m *nilManagedClusterAgentPoolInformation) IsAvailable() bool {
	return true
}

func (m *nilManagedClusterAgentPoolInformation) GetResourceID() string {
	return ""
}

func (m *nilManagedClusterAgentPoolInformation) GetSubscriptionID() string {
	return ""
}

func (m *nilManagedClusterAgentPoolInformation) GetResourceGroup() string {
	return ""
}

func (m *nilManagedClusterAgentPoolInformation) GetManagedClusterName() string {
	return ""
}

func (m *nilManagedClusterAgentPoolInformation) GetResourceName() string {
	return ""
}

func (m *nilManagedClusterAgentPoolInformation) GetLatestModel(
	ctx context.Context,
) (containerservice.ManagedClusterAgentPoolProfile, error) {
	return containerservice.ManagedClusterAgentPoolProfile{}, ErrNotAvilable
}

type managedClusterInfomration struct {
	azureAuthorizer autorest.Authorizer

	resourceID *ARMResourceID

	mutex *sync.RWMutex // protects the following fields
	model *containerservice.ManagedCluster
}

var _ ManagedClusterInformation = &managedClusterInfomration{}

func (m *managedClusterInfomration) IsAvailable() bool {
	return true
}

func (m *managedClusterInfomration) GetResourceID() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return to.String(m.model.ID)
}

func (m *managedClusterInfomration) GetSubscriptionID() string {
	return m.resourceID.Subscription
}

func (m *managedClusterInfomration) GetRegion() string {
	m.mutex.RLock()
	defer m.mutex.Unlock()

	return to.String(m.model.Location)
}

func (m *managedClusterInfomration) GetResourceGroup() string {
	return m.resourceID.ResourceGroup
}

func (m *managedClusterInfomration) GetResourceName() string {
	return m.resourceID.ResourceName
}

func (m *managedClusterInfomration) GetNodeResourceGroup() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return *m.model.NodeResourceGroup
}

func (m *managedClusterInfomration) GetLatestModel(ctx context.Context) (containerservice.ManagedCluster, error) {
	client := aksManagedClusterClient(m.azureAuthorizer, m.resourceID.Subscription)
	managedCluster, err := client.Get(ctx, m.resourceID.ResourceGroup, m.resourceID.ResourceName)
	if err != nil {
		return containerservice.ManagedCluster{}, err
	}

	// update record
	m.mutex.Lock()
	m.model = &managedCluster
	m.mutex.Unlock()

	return managedCluster, nil
}

func (m *managedClusterInfomration) GetKubeConfig(ctx context.Context) (string, error) {
	// TODO
	return "", ErrNotAvilable
}

func (m *managedClusterInfomration) GetAgentPoolInformation(
	ctx context.Context, agentPoolName string,
) (ManagedClusterAgentPoolInformation, error) {
	return nil, ErrNotAvilable
}
