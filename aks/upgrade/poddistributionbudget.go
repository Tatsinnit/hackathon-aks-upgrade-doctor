package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/kube"
	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/rules"
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
				// podlist, err := utils.GetPods(kubeClient, namespace.Name)
				// if err != nil {
				// 	return err
				// }
				for _, i := range podDistInterface.Items {
					fmt.Println("Pod Disruption Budget Name : ", i.Name)
					fmt.Println("Min Available : ", fmt.Sprint(i.Spec.MinAvailable))
					fmt.Println("Max Available : ", fmt.Sprint(i.Spec.MaxUnavailable))
					fmt.Println("DisruptionsAllowed : ", fmt.Sprint(i.Status.DisruptionsAllowed))

					// ok, reason := checkPDBForUpgrade(
					// 	podlist.Items,
					// 	i,
					// )
					// if !ok {
					// 	fmt.Println(reason)
					// }
				}
			}

			return nil
		},
	}

	return cmd
}

// func checkPDBForUpgrade(
// 	podsInTheSameNamespace []corev1.Pod,
// 	pdb policyv1beta1.PodDisruptionBudget,
// ) (ok bool, reason string) {
// 	// The non-zero value for ALLOWED DISRUPTIONS means that the disruption controller has seen the pods, counted the matching pods, and updated the status of the PDB.
// 	if fmt.Sprint(pdb.Status.DisruptionsAllowed) == "0" && fmt.Sprint(pdb.Spec.MaxUnavailable) == "0" {
// 		return false, "Upgrade operation will fail - you are requiring zero voluntary evictions, so cannot successfully drain a Node running one of the Pods"
// 	}

// 	if pdb.Spec.MinAvailable != nil && fmt.Sprint(pdb.Status.DisruptionsAllowed) == "0" && fmt.Sprint(pdb.Spec.MinAvailable) > "0" {
// 		count := 0
// 		for _, pod := range podsInTheSameNamespace {

// 			if utils.IsMapSubset(pod.Labels, pdb.Spec.Selector.MatchLabels) {
// 				count = count + 1
// 			}
// 		}

// 		if count != 0 {
// 			diff := count - pdb.Spec.MinAvailable.IntValue()

// 			if diff == 0 {
// 				return false, "Upgrade operation will fail - cannot successfully drain a Node running one of the Pods"
// 			}
// 		}
// 	}
// 	return true, ""
// }

func checkPDBForUpgrade(pdb policyv1beta1.PodDisruptionBudget) (ok bool, reason string) {

	// Adding this condition to reduce the number of if condition checks in this function.
	if fmt.Sprint(pdb.Status.DisruptionsAllowed) == "0" {

		// If there are no matching pods as per PDB spec, ExpectedPods (Total parameter in Kubectl output) count will be 0.
		if pdb.Status.ExpectedPods != 0 {

			// If current healthy pods are lesser than or equal to desired healthy pods, there will be 0 disruptions allowed.
			// And this will cause Node drain failure during Upgrade operation and the operation will fail.
			if pdb.Status.CurrentHealthy-pdb.Status.DesiredHealthy < 1 {
				return false, "Upgrade operation on this AKS cluster will fail. PDB " + pdb.Name + " in Namespace " + pdb.Namespace +
					" should allow at least 1 pod replica to be moved at a time. Please refer https://aka.ms/aksupgrade "

				// Below is used for debugging.
				//return false, "Failure, PDB - " + pdb.Name + " DisruptionsAllowed - " + fmt.Sprint(pdb.Status.DisruptionsAllowed) + " CurrentHealthy - " + fmt.Sprint(pdb.Status.CurrentHealthy) + " DesiredHealthy - " + fmt.Sprint(pdb.Status.DesiredHealthy) + " ExpectedPods - " + fmt.Sprint(pdb.Status.ExpectedPods) + " "
			}
		}
	}
	return true, ""
}

var upgradePDBRuleProvider = rules.NewRule(
	"upgrade/pdb",
	func(
		ctx context.Context,
		clusterCtx rules.ClusterContext,
	) ([]*rules.CheckResult, error) {
		kubeClient, err := clusterCtx.GetClusterKubeClient()
		if err != nil {
			return nil, err
		}

		namespacesList, err := kubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("unable to list namespaces in the cluster: %w", err)
		}

		var rv []*rules.CheckResult
		for _, namespace := range namespacesList.Items {
			podDistInterface, err := kubeClient.PolicyV1beta1().PodDisruptionBudgets(namespace.Name).List(context.Background(), metav1.ListOptions{}) //.Get(context.Background(), namespace.Name, metav1.GetOptions{})
			if err != nil {
				// FIXME: should check with best-effort
				return nil, fmt.Errorf("PDB error cluster: %w", err)
			}

			// podlist, err := utils.GetPods(kubeClient, namespace.Name)
			// if err != nil {
			// 	return nil, err
			// }

			// for _, i := range podDistInterface.Items {
			// 	ok, reason := checkPDBForUpgrade(
			// 		podlist.Items,
			// 		i,
			// 	)
			// 	if !ok {
			// 		rv = append(rv, &rules.CheckResult{
			// 			RuleID:      "upgrade/pdb",
			// 			Category:    rules.Warning,
			// 			Description: reason,
			// 		})
			// 	}
			// }

			for _, pdb := range podDistInterface.Items {
				ok, reason := checkPDBForUpgrade(pdb)
				if !ok {
					rv = append(rv, &rules.CheckResult{
						RuleID:      "Upgrade/PDB",
						Category:    rules.Warning,
						Description: reason,
					})
				}
			}
		}

		return rv, nil
	},
)
