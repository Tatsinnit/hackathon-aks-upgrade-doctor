package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/rules"
	"github.com/spf13/cobra"
)

type demoRule struct {
}

func (d demoRule) RuleID() string {
	return "demo"
}

func (d demoRule) GetCheckResult(
	ctx context.Context,
	clusterCtx rules.ClusterContext,
) (*rules.CheckResult, error) {
	time.Sleep(500 * time.Millisecond)

	return &rules.CheckResult{
		RuleID:      d.RuleID(),
		Category:    rules.Healthy,
		Description: "hello, world",
	}, nil
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
					demoRule{},
					demoRule{},
					demoRule{},
					demoRule{},
				},
			)
			if err != nil {
				return err
			}

			fmt.Println(results)

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
