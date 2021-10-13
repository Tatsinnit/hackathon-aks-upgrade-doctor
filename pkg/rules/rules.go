package rules

import (
	"context"
)

type GetCheckResultsFunc func(ctx context.Context, clusterCtx ClusterContext) ([]*CheckResult, error)

type simpleRulesProvider struct {
	ruleID    string
	checkFunc GetCheckResultsFunc
}

var _ RuleProvider = &simpleRulesProvider{}

func (r *simpleRulesProvider) RuleID() string {
	return r.ruleID
}

func (r *simpleRulesProvider) GetCheckResults(
	ctx context.Context,
	clusterCtx ClusterContext,
) ([]*CheckResult, error) {
	return r.checkFunc(ctx, clusterCtx)
}

// NewRule creates rule provider from function.
func NewRule(
	ruleID string,
	checkFunc GetCheckResultsFunc,
) RuleProvider {
	return &simpleRulesProvider{
		ruleID:    ruleID,
		checkFunc: checkFunc,
	}
}
