package main

import (
	"context"

	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/rules"
)

var subnetCapacityRuleProvider = rules.NewRule(
	"network/subnet",
	func(
		ctx context.Context,
		clusterCtx rules.ClusterContext,
	) ([]*rules.CheckResult, error) {
		// TODO process:
		// 1. get the cluster information from clusterCtx
		// 2. get the node resource group
		// 3. list the vnet / subnet resources under the node resource group
		//     * we can follow this demo to construct the client: https://github.com/Azure-Samples/azure-sdk-for-go-samples/blob/b49c4162aa1d96bc2b1b42afecbf4a21b420e568/network/subnets.go#L18-L24
		// 4. calculate the capacity of the subnet and calculate result

		return nil, nil
	},
)
