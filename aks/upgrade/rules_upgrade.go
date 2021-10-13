package main

import "github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/rules"

// rulesSetForUpgrade - list of rules to check for AKS cluster upgrade operation
var rulesSetForUpgrade = []rules.RuleProvider{
	upgradePDBRuleProvider,
}
