package kube

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// NewKubeClient creates a kubernetes clientset with environment detection.
// It loads kube client config with following orders:
// 1. If kubeConfigFilePath is set and pointed to a valid path, use it
// 2. If recommended kubeconfig env variable is set, use it
// 3. If it's in-cluster environment, use in-cluster config
// 4. Try looks up the kube config from default recommended locations ($HOME/.kube/config)
func NewKubeClient(kubeConfigFilePath string) (kubernetes.Interface, error) {
	restConfig, err := loadRestConfig(kubeConfigFilePath)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(restConfig)
}

// based on: https://github.com/kubernetes-sigs/controller-runtime/blob/3c54acbad091d40ef80d467943c5d7126c4b8291/pkg/client/config/config.go#L99
func loadRestConfig(kubeconfigPath string) (*rest.Config, error) {
	if len(kubeconfigPath) > 0 {
		return loadConfigWithDefaultContext("", &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath})
	}

	// If the recommended kubeconfig env variable is not specified,
	// try the in-cluster config.
	kubeconfigPath = os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(kubeconfigPath) == 0 {
		if c, err := rest.InClusterConfig(); err == nil {
			return c, nil
		}
	}

	// If the recommended kubeconfig env variable is set, or there
	// is no in-cluster config, try the default recommended locations.
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if _, ok := os.LookupEnv("HOME"); !ok {
		u, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("could not get current user: %w", err)
		}
		loadingRules.Precedence = append(loadingRules.Precedence, path.Join(u.HomeDir, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName))
	}

	return loadConfigWithDefaultContext("", loadingRules)
}

func loadConfigWithDefaultContext(apiServerURL string, loader clientcmd.ClientConfigLoader) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		loader,
		&clientcmd.ConfigOverrides{
			ClusterInfo: clientcmdapi.Cluster{
				Server: apiServerURL,
			},
			CurrentContext: "",
		}).ClientConfig()
}
