package main

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest/to"
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

		clusterInfo, err := clusterCtx.GetManagedClusterInformation(ctx)
		if err != nil {
			return nil, err
		}

		auth, err := clusterCtx.GetAzureAuthorizer()
		if err != nil {
			return nil, err
		}

		clusterModel, err := clusterInfo.GetLatestModel(ctx)
		if err != nil {
			return nil, err
		}

		subnetsClient := network.NewSubnetsClient(clusterInfo.GetSubscriptionID())
		subnetsClient.Authorizer = auth

		vnetsClient := network.NewVirtualNetworksClient(clusterInfo.GetSubscriptionID())
		vnetsClient.Authorizer = auth

		var results []*rules.CheckResult
		for _, apProfile := range *clusterModel.AgentPoolProfiles {
			ap, err := clusterInfo.GetAgentPoolInformation(ctx, to.String(apProfile.Name))
			if err != nil {
				return nil, err
			}

			apModel, err := ap.GetLatestModel(ctx)
			if err != nil {
				return nil, err
			}

			// TODO: Get the subnet object from *apModel.VnetSubnetID
			// Problem: that value is nil. Can't see a way to tell the go client to retrieve it.
			// From the rest API docs (https://docs.microsoft.com/en-us/rest/api/aks/agent-pools/get#agentpool)
			// it looks like this property will be populated if the customer BYOs their own subnet....
			// So may have to fall back to picking a vnet/subnet from the cluster node pool based on
			// agent pool name, if the value is not specified.

			// Try looping through the vnets in the node resource group
			vnetsResult, err := vnetsClient.List(ctx, clusterInfo.GetNodeResourceGroup())
			if err != nil {
				return nil, err
			}

			for _, vnet := range vnetsResult.Values() {
				// Does subnet name match agent pool name?
				subnetsClient.Get(ctx, clusterInfo.GetNodeResourceGroup(), *vnet.Name, *apModel.Name, "")
				// for _, subnet := vnet.Subnets {
				// 	subnetsClient.Get(ctx, clusterInfo.GetNodeResourceGroup(), *vnet.Name, ag, "")
				// }
			}

			// results = append(results, &rules.CheckResult{
			// 	RuleID:   "upgrade/armtest-agent-pool",
			// 	Category: category,
			// 	Description: fmt.Sprintf(
			// 		"Got subnet from agent pool: %s - subnetid: %s",
			// 		*apProfile.Name,
			// 		*apModel.VnetSubnetID,
			// 	),
			// })
		}

		return results, nil
	},
)
