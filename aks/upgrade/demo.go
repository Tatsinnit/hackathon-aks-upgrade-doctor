package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/azure"
	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/report"
	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/rules"
	"github.com/spf13/cobra"
)

type demoRule struct {
	ruleID      string
	category    rules.ResultCategory
	description string
}

func (d demoRule) RuleID() string {
	return d.ruleID
}

func (d demoRule) GetCheckResults(
	ctx context.Context,
	clusterCtx rules.ClusterContext,
) ([]*rules.CheckResult, error) {
	time.Sleep(500 * time.Millisecond)

	return []*rules.CheckResult{{
		RuleID:      d.ruleID,
		Category:    d.category,
		Description: d.description,
	}}, nil
}

func createEngineDemoCommand() *cobra.Command {
	var (
		flagClusterKubeConfigFilePath string
		flagClusterResourceID         string
	)

	cmd := &cobra.Command{
		Use:   "engine-demo",
		Short: "demo for rules engine",
		RunE: func(cmd *cobra.Command, args []string) error {
			authorizer, err := auth.NewAuthorizerFromCLI()
			if err != nil {
				return fmt.Errorf("needs AZ CLI authentcation support: %w", err)
			}

			createClusterCtx := rules.CreateClusterContextOptions{
				ClusterKubeConfigPath:     flagClusterKubeConfigFilePath,
				AzureAuthorizer:           authorizer,
				ManagedClusterInformation: azure.NilManagedClsuterInformation(),
			}
			if flagClusterResourceID != "" {
				// user has specified an cluster resource id, try load cluster from it
				cluster, err := azure.LoadManagedClusterInformationFromResourceID(authorizer, flagClusterResourceID)
				if err != nil {
					// user specified a wrong input...
					return err
				}
				// successfully loaded the cluster
				createClusterCtx.ManagedClusterInformation = cluster
			}

			clusterCtx, err := createClusterCtx.Create()
			if err != nil {
				return err
			}

			engine := rules.NewEngine(cmd.OutOrStdout())

			results, err := engine.CheckRulesSet(
				context.Background(),
				clusterCtx,
				rules.RulesSet{
					upgradePDBRuleProvider,
					subnetCapacityRuleProvider,
					rules.NewRule(
						"upgrade/armtest-managed-cluster",
						func(
							ctx context.Context,
							clusterCtx rules.ClusterContext,
						) ([]*rules.CheckResult, error) {
							// load the cluster from cluster context
							cluster, err := clusterCtx.GetManagedClusterInformation(ctx)
							if err != nil {
								return nil, err
							}

							// load the ARM representation of the cluster
							// the model details: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-08-01/containerservice#ManagedCluster
							latestModel, err := cluster.GetLatestModel(ctx)
							if err != nil {
								return nil, err
							}

							// check provisioning state
							category := rules.Healthy
							provisionState := to.String(latestModel.ProvisioningState)
							if provisionState != "Succeeded" {
								category = rules.Warning
							}

							return []*rules.CheckResult{
								{
									RuleID:   "upgrade/armtest-managed-cluster",
									Category: category,
									Description: fmt.Sprintf(
										"Got details from cluster: %s - state: %s (%s)",
										cluster.GetResourceName(),
										provisionState,
										cluster.GetNodeResourceGroup(),
									),
								},
							}, nil
						},
					),
					rules.NewRule(
						"upgrade/armtest-agent-pool",
						func(
							ctx context.Context,
							clusterCtx rules.ClusterContext,
						) ([]*rules.CheckResult, error) {
							cluster, err := clusterCtx.GetManagedClusterInformation(ctx)
							if err != nil {
								return nil, err
							}

							clusterModel, err := cluster.GetLatestModel(ctx)
							if err != nil {
								return nil, err
							}

							ap, err := cluster.GetAgentPoolInformation(ctx, to.String((*clusterModel.AgentPoolProfiles)[0].Name))
							if err != nil {
								return nil, err
							}

							// load the ARM representation of the cluster
							// the model details: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-08-01/containerservice#AgentPool
							latestModel, err := ap.GetLatestModel(ctx)
							if err != nil {
								return nil, err
							}

							category := rules.Healthy
							provisionState := to.String(latestModel.ProvisioningState)
							if provisionState != "Succeeded" {
								category = rules.Warning
							}

							return []*rules.CheckResult{
								{
									RuleID:   "upgrade/armtest-agent-pool",
									Category: category,
									Description: fmt.Sprintf(
										"Got details from agent pool: %s - state: %s (%s)",
										ap.GetResourceName(),
										provisionState,
										ap.GetManagedClusterName(),
									),
								},
							}, nil
						},
					),
					demoRule{
						ruleID:      "demo/upgrade/subnet",
						category:    rules.Warning,
						description: "cluster subnet is almost full",
					},
					demoRule{
						ruleID:      "demo/version/out-of-date-version",
						category:    rules.Advisory,
						description: "cluster version 1.17.11 is out-of-date",
					},
					demoRule{
						ruleID:      "demo/control-plane/coredns",
						category:    rules.Healthy,
						description: "CoreDNS pods are running normally",
					},
				},
			)
			if err != nil {
				return err
			}

			p := report.FancyCheckResultPresenter{
				ReportName:   "demo cluster result",
				CheckResults: results,
			}
			if err := p.Present(cmd.OutOrStdout()); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(
		&flagClusterKubeConfigFilePath,
		"kube-config",
		"",
		"cluster kubeconfig to use",
	)
	cmd.Flags().StringVar(
		&flagClusterResourceID,
		"aks-resource-id",
		"",
		"resource id for the AKS cluster",
	)

	return cmd
}
