package rules

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/containerservice/armcontainerservice"
	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/kube"
	"k8s.io/client-go/kubernetes"
)

// CreateClusterContextOptions creates a cluster context.
type CreateClusterContextOptions struct {
	// ClusterKubeConfigPath is the path to the kubeconfig file for the cluster.
	ClusterKubeConfigPath string
}

// Create creates a cluster context.
func (opts CreateClusterContextOptions) Create(
	subscriptionId string,
	resourceGroup string,
	resourceName string,
) (ClusterContext, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO
	}

	con := arm.NewDefaultConnection(cred, nil)

	return &clusterContextImpl{
		ClusterKubeConfigPath: opts.ClusterKubeConfigPath,
		ArmConnection:         con,
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		ResourceName:          resourceName,
	}, nil
}

type clusterContextImpl struct {
	ClusterKubeConfigPath string
	ArmConnection         *arm.Connection
	SubscriptionId        string
	ResourceGroup         string
	ResourceName          string
}

var _ ClusterContext = &clusterContextImpl{}

func (ctx *clusterContextImpl) GetClusterKubeClient() (kubernetes.Interface, error) {
	return kube.NewKubeClient(ctx.ClusterKubeConfigPath)
}

func (ctx *clusterContextImpl) GetAKSClusterResourceDetails() string {
	client := armcontainerservice.NewManagedClustersClient(ctx.ArmConnection, ctx.SubscriptionId)
	resp, err := client.Get(context.TODO(), ctx.ResourceGroup, ctx.ResourceName, &armcontainerservice.ManagedClustersGetOptions{})
	if err != nil {
		// TODO:
	}

	return *resp.Properties.Fqdn
}
