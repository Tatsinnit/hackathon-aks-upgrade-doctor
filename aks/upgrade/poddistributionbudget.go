package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	fmt.Println("Hello, world.")

	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Cannot get user home dir: %v", err)
	}

	master := ""
	kubeconfig := path.Join(dirname, ".kube/config")
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)

	// config, err := restclient.InClusterConfig()
	if err != nil {
		log.Fatalf("Cannot load kubeconfig: %v", err)
	}

	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(fmt.Errorf("Getting access to K8S failed: %w", err))
	}

	namespacesList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(fmt.Errorf("Unable to list namespaces in the cluster: %w", err))
	}
	for _, namespace := range namespacesList.Items {

		podDistInterface, err := clientset.PolicyV1().PodDisruptionBudgets(namespace.Name).List(context.Background(), metav1.ListOptions{}) //.Get(context.Background(), namespace.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(fmt.Errorf("PDB error cluster: %w", err))
		}
		for _, i := range podDistInterface.Items {
			fmt.Println("Pod Disruption Budget Name : ", i.Name)
			fmt.Println("Min Available : ", i.Spec.MinAvailable)
			fmt.Println("Max Available : ", i.Spec.MaxUnavailable)
			fmt.Println("DisruptionsAllowed : ", i.Status.DisruptionsAllowed)
		}
	}
}
