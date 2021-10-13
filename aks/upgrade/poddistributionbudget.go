package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/aks/upgrade/utils"
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
				podDistInterface, err := kubeClient.PolicyV1beta1().PodDisruptionBudgets(namespace.Name).List(context.Background(), metav1.ListOptions{}) //.Get(context.Background(), namespace.Name, metav1.GetOptions{})
				if err != nil {
					return fmt.Errorf("PDB error cluster: %w", err)
				}
				for _, i := range podDistInterface.Items {
					fmt.Println("Pod Disruption Budget Name : ", i.Name)
					fmt.Println("Min Available : ", fmt.Sprint(i.Spec.MinAvailable))
					fmt.Println("Max Available : ", fmt.Sprint(i.Spec.MaxUnavailable))
					fmt.Println("DisruptionsAllowed : ", fmt.Sprint(i.Status.DisruptionsAllowed))

					// The non-zero value for ALLOWED DISRUPTIONS means that the disruption controller has seen the pods, counted the matching pods, and updated the status of the PDB.
					if fmt.Sprint(i.Status.DisruptionsAllowed) == "0" && fmt.Sprint(i.Spec.MaxUnavailable) == "0" {
						fmt.Println("Upgrade operation will fail - you are requiring zero voluntary evictions, so cannot successfully drain a Node running one of the Pods")
					}

					if i.Spec.MinAvailable != nil && fmt.Sprint(i.Status.DisruptionsAllowed) == "0" && fmt.Sprint(i.Spec.MinAvailable) > "0" {
						podlist, err := utils.GetPods(kubeClient, namespace.Name)
						if err != nil {
							return err
						}
						count := 0
						for _, pod := range podlist.Items {

							if utils.IsMapSubset(pod.Labels, i.Spec.Selector.MatchLabels) {
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
