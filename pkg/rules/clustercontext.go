package rules

import (
	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/kube"
	"k8s.io/client-go/kubernetes"
)

// CreateClusterContextOptions creates a cluster context.
type CreateClusterContextOptions struct {
	// ClusterKubeConfigPath is the path to the kubeconfig file for the cluster.
	ClusterKubeConfigPath string
}

// Create creates a cluster context.
func (opts CreateClusterContextOptions) Create() (ClusterContext, error) {
	// TODO: acquire Azure client
	return &clusterContextImpl{
		ClusterKubeConfigPath: opts.ClusterKubeConfigPath,
	}, nil
}

type clusterContextImpl struct {
	ClusterKubeConfigPath string
}

var _ ClusterContext = &clusterContextImpl{}

func (ctx *clusterContextImpl) GetClusterKubeClient() (kubernetes.Interface, error) {
	return kube.NewKubeClient(ctx.ClusterKubeConfigPath)
}
