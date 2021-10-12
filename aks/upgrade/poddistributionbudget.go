package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

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
		fmt.Println(fmt.Errorf("getting access to K8S failed: %w", err))
	}

	namespacesList, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(fmt.Errorf("unable to list namespaces in the cluster: %w", err))
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

			// The non-zero value for ALLOWED DISRUPTIONS means that the disruption controller has seen the pods, counted the matching pods, and updated the status of the PDB.
			if i.Status.DisruptionsAllowed == 0 && i.Spec.MaxUnavailable.String() == "0" {
				fmt.Println("Upgrade operation will fail - you are requiring zero voluntary evictions, so cannot successfully drain a Node running one of the Pods")
			}

			// Below code tries to implement following pseudo code
			// 	Count = Count of (kubectl get pods --selector=<Labels provided in PDB>)
			// 	if (Count != 0)
			// 		{
			// 			If (Count - "MIN AVAILABLE" == 0)
			// 			{
			// 				Printf("Upgrade operation will fail - cannot successfully drain a Node running one of the Pods")
			//
			// 			}
			// 		}
			// Note: Above pseudo code implemenation below

			if i.Status.DisruptionsAllowed == 0 && i.Spec.MinAvailable.String() > "0" {
				podlist, err := GetPods(clientset, namespace.Name)
				if err != nil {
					fmt.Println(err.Error())
				}
				count := 0
				for _, pod := range podlist.Items {
					if reflect.DeepEqual(pod.Labels, i.Labels) {
						count = count + 1
					}
				}
				if count != 0 {
					diff := count - i.Spec.MinAvailable.IntValue()
					if diff == 0 {
						fmt.Println("Upgrade operation will fail - cannot successfully drain a Node running one of the Pods")
					}
				}
			}

		}
	}
}

func GetPods(clientset *kubernetes.Clientset, namespace string) (*v1.PodList, error) {
	// Create a pod interface for the given namespace
	podInterface := clientset.CoreV1().Pods(namespace)

	// List the pods in the given namespace
	podList, err := podInterface.List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("getting pods failed: %w", err)
	}

	return podList, nil
}
