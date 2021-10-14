package rules

import (
	"context"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/azure"
	"k8s.io/client-go/kubernetes"
)

// ResultCategory specifies the result category of a health check.
type ResultCategory string

const (
	// Healthy - The health check returned a healthy result.
	Healthy ResultCategory = "healthy"
	// Advisory - The health check returned a advisory result.
	Advisory ResultCategory = "advisory"
	// Warning - The health check returned a warning result.
	Warning ResultCategory = "warning"
	// Failed - The health check failed. This indicates checker error, (usually) not the issue in cluster.
	Failed ResultCategory = "failed"
)

// CheckResult defines the health check result of an AKS cluster.
type CheckResult struct {
	// RuleID - the id of the health check rule.
	RuleID string `json:"ruleID"`
	// Category - the result category of the health check.
	Category ResultCategory `json:"category"`
	// Description - health check reuslt description.
	Description string `json:"description"`
}

// ClusterContext provides the information for a cluster.
type ClusterContext interface {
	// GetClusterKubeClient returns a kubernetes client instance.
	GetClusterKubeClient() (kubernetes.Interface, error)

	// GetAzureAuthorizer returns the Azure SDK authorizer.
	GetAzureAuthorizer() (autorest.Authorizer, error)

	// GetAKSClusterResourceDetails returns the AKS cluster resource details.
	GetManagedClusterInformation(ctx context.Context) (azure.ManagedClusterInformation, error)
}

// RuleProvider provides and checks a health check rule for an AKS cluster.
type RuleProvider interface {
	// RuleID returns the rule id.
	// Rule id should be globally unique.
	RuleID() string

	// GetCheckResults executes the health check rule and returns the health check results.
	GetCheckResults(ctx context.Context, clusterCtx ClusterContext) ([]*CheckResult, error)
}

// RulesSet defines a set of rules.
type RulesSet []RuleProvider

// Engine executes a collection of health check rules for an AKS cluster.
type Engine interface {
	// CheckRulesSet checks a set of rules for an AKS cluster.
	CheckRulesSet(ctx context.Context, clusterCtx ClusterContext, rs RulesSet) ([]*CheckResult, error)
}
