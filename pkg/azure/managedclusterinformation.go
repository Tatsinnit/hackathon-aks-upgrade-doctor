package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-08-01/containerservice"
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
	GetModel(ctx context.Context) (containerservice.ManagedCluster, error)

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
	GetLatestModel(ctx context.Context) (containerservice.ManagedClusterAgentPoolProfile, error)

	// TODO: GetVMSSName()
}
