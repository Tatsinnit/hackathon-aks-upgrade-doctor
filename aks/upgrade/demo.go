package main

import (
	"context"
	"fmt"
	"time"

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
	)

	cmd := &cobra.Command{
		Use:   "engine-demo",
		Short: "demo for rules engine",
		RunE: func(cmd *cobra.Command, args []string) error {
			createClusterCtx := rules.CreateClusterContextOptions{
				ClusterKubeConfigPath: flagClusterKubeConfigFilePath,
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
					rules.NewRule(
						"upgrade/armtest",
						func(
							ctx context.Context,
							clusterCtx rules.ClusterContext,
						) ([]*rules.CheckResult, error) {
							// details := clusterCtx.GetAKSClusterResourceDetails()

							return []*rules.CheckResult{
								{
									RuleID:      "upgrade/armtest",
									Category:    rules.Advisory,
									Description: fmt.Sprintf("Got details from cluster: %s", ""),
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

	return cmd
}
