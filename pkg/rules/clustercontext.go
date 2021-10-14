package rules

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/azure"
	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/kube"
	"k8s.io/client-go/kubernetes"
)

// CreateClusterContextOptions creates a cluster context.
type CreateClusterContextOptions struct {
	// ClusterKubeConfigPath is the path to the kubeconfig file for the cluster.
	ClusterKubeConfigPath string

	// AzureAuthorizer is the authorizer to use for Azure resources.
	AzureAuthorizer autorest.Authorizer

	// ManagedClusterInformation is the information about the cluster.
	ManagedClusterInformation azure.ManagedClusterInformation
}

// Create creates a cluster context.
func (opts CreateClusterContextOptions) Create() (ClusterContext, error) {
	return &clusterContextImpl{
		ClusterKubeConfigPath:     opts.ClusterKubeConfigPath,
		AzureAuthorizer:           opts.AzureAuthorizer,
		ManagedClusterInformation: opts.ManagedClusterInformation,
	}, nil
}

type clusterContextImpl struct {
	ClusterKubeConfigPath     string
	AzureAuthorizer           autorest.Authorizer
	ManagedClusterInformation azure.ManagedClusterInformation
}

var _ ClusterContext = &clusterContextImpl{}

func (clusterCtx *clusterContextImpl) GetClusterKubeClient() (kubernetes.Interface, error) {
	// TODO: use managed cluster's kubeconfig if possible
	return kube.NewKubeClient(clusterCtx.ClusterKubeConfigPath)
}

func (clusterCtx *clusterContextImpl) GetAzureAuthorizer() (autorest.Authorizer, error) {
	return clusterCtx.AzureAuthorizer, nil
}

func (clusterCtx *clusterContextImpl) GetManagedClusterInformation(
	ctx context.Context,
) (azure.ManagedClusterInformation, error) {
	return clusterCtx.ManagedClusterInformation, nil
}
