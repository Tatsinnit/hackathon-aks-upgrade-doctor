package azure

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-08-01/containerservice"
	"github.com/Azure/go-autorest/autorest"
)

// ManagedClusterInformation provides information about a managed cluster.
type ManagedClusterInformation interface {
	// IsAvailable tells if the managed cluster is available.
	IsAvailable() bool

	// GetResourceID returns the ARM resource ID of the managed cluster.
	GetResourceID() string

	// GetSubscription returns the subscription id of the managed cluster.
	GetSubscriptionID() string

	// GetResourceID returns the resource group of the managed cluster.
	GetResourceGroup() string

	// GetResourceName returns the resource name of the managed cluster.
	GetResourceName() string

	// GetRegion returns the region of the managed cluster.
	GetRegion() string

	// GetNodeResourceGroup returns the node resource group of the cluster.
	GetNodeResourceGroup() string

	// GetLatestModel returns the latest properties of the managed cluster.
	GetLatestModel(ctx context.Context) (containerservice.ManagedCluster, error)

	// GetKubeConfig retrieves the cluster kubeconfig.
	GetKubeConfig(ctx context.Context) (string, error)

	// GetAgentPoolInformation retrieves the agent pool information.
	GetAgentPoolInformation(ctx context.Context, agentPoolName string) (ManagedClusterAgentPoolInformation, error)
}

// ManagedClusterAgentPoolInformation provides information about a managed cluster agent pool.
type ManagedClusterAgentPoolInformation interface {
	// IsAvailable tells if the managed cluster agent pool is available.
	IsAvailable() bool

	// GetResourceID returns the ARM resource ID of the managed cluster agent pool.
	GetResourceID() string

	// GetSubscription returns the subscription id of the managed cluster.
	GetSubscriptionID() string

	// GetResourceID returns the resource group of the managed cluster.
	GetResourceGroup() string

	// GetResourceName returns the resource name of the managed cluster.
	GetManagedClusterName() string

	// GetResourceName returns the resource name of the managed cluster.
	GetResourceName() string

	// GetLatestModel returns the latest properties of the managed cluster agent pool.
	GetLatestModel(ctx context.Context) (containerservice.AgentPool, error)

	// TODO: GetVMSSName()
}

// NilManagedClsuterInformation creates a nil ManagedClusterInformation instance.
func NilManagedClsuterInformation() ManagedClusterInformation {
	return &nilManagedClusterInformation{}
}

// LoadManagedClusterInformationFromResourceID loads a ManagedClusterInformation instance from cluster's resource ID.
func LoadManagedClusterInformationFromResourceID(
	authorizer autorest.Authorizer,
	resourceID string,
) (ManagedClusterInformation, error) {
	armResourceID, err := ParseARMResourceID(resourceID)
	if err != nil {
		return &nilManagedClusterInformation{}, err
	}

	rv := &managedClusterInfomration{
		azureAuthorizer: authorizer,
		resourceID:      armResourceID,
		mutex:           &sync.RWMutex{},
	}

	// load the cluster for first time
	loadCtx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if _, err := rv.GetLatestModel(loadCtx); err != nil {
		return &nilManagedClusterInformation{}, fmt.Errorf("unable to get cluster form %s: %w", resourceID, err)
	}

	return rv, nil
}
