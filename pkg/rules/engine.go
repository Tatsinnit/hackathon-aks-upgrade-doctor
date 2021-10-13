package rules

import (
	"context"
	"fmt"
	"io"

	"github.com/gosuri/uiprogress"
)

type engineImpl struct {
	Stdout io.Writer
}

// NewEngine creates a new engine.
func NewEngine(stdout io.Writer) Engine {
	return &engineImpl{
		Stdout: stdout,
	}
}

var _ Engine = &engineImpl{}

func (e *engineImpl) CheckRulesSet(
	ctx context.Context,
	clusterCtx ClusterContext,
	rs RulesSet,
) ([]*CheckResult, error) {
	if len(rs) < 1 {
		return nil, nil
	}

	progress := uiprogress.New()

	bar := progress.AddBar(len(rs))
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("checking %d/%d", b.Current(), b.Total)
	})

	progress.Start()
	defer progress.Stop()

	var checkResults []*CheckResult
	for _, rule := range rs {
		results, err := rule.GetCheckResults(ctx, clusterCtx)
		if err != nil {
			results = []*CheckResult{{
				RuleID:      rule.RuleID(),
				Category:    Failed,
				Description: fmt.Sprintf("check result failed: %s", err.Error()),
			}}
		}
		checkResults = append(checkResults, results...)

		bar.Incr()
	}

	return checkResults, nil
}
