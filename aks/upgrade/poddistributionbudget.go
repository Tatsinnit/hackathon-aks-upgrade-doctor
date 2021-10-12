package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/kube"
)

func createDemoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "demo for pdb",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Creates the clientset
			kubeClient, err := kube.NewKubeClient("")
			if err != nil {
				return fmt.Errorf("construct kube client failed: %w", err)
			}

			namespacesList, err := kubeClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
			if err != nil {
				return fmt.Errorf("unable to list namespaces in the cluster: %w", err)
			}
			for _, namespace := range namespacesList.Items {
				podDistInterface, err := kubeClient.PolicyV1().PodDisruptionBudgets(namespace.Name).List(context.Background(), metav1.ListOptions{}) //.Get(context.Background(), namespace.Name, metav1.GetOptions{})
				if err != nil {
					return fmt.Errorf("PDB error cluster: %w", err)
				}
				for _, i := range podDistInterface.Items {
					fmt.Println("Pod Disruption Budget Name : ", i.Name)
					fmt.Println("Min Available : ", i.Spec.MinAvailable)
					fmt.Println("Max Available : ", i.Spec.MaxUnavailable)
					fmt.Println("DisruptionsAllowed : ", i.Status.DisruptionsAllowed)

					// The non-zero value for ALLOWED DISRUPTIONS means that the disruption controller has seen the pods, counted the matching pods, and updated the status of the PDB.
					if i.Status.DisruptionsAllowed == 0 && fmt.Sprint(i.Spec.MaxUnavailable) == "0" {
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

					if i.Status.DisruptionsAllowed == 0 && fmt.Sprint(i.Spec.MinAvailable) > "0" {
						podlist, err := GetPods(kubeClient, namespace.Name)
						if err != nil {
							return err
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

			return nil
		},
	}

	return cmd
}

func GetPods(clientset kubernetes.Interface, namespace string) (*corev1.PodList, error) {
	// Create a pod interface for the given namespace
	podInterface := clientset.CoreV1().Pods(namespace)

	// List the pods in the given namespace
	podList, err := podInterface.List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("getting pods failed: %w", err)
	}

	return podList, nil
}
