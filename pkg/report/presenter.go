package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/Tatsinnit/hackathon-aks-upgrade-doctor/pkg/rules"
	"github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"
)

type Presentable interface {
	Present(io.Writer) error
}

type FancyCheckResultPresenter struct {
	ReportName   string
	CheckResults []*rules.CheckResult
}

func (p FancyCheckResultPresenter) formatCategory(c rules.ResultCategory) interface{} {
	content := strings.ToUpper(string(c))
	switch c {
	case rules.Warning:
		return aurora.Blink(aurora.BgBlack(aurora.Yellow(aurora.Bold("!" + content))))
	case rules.Advisory:
		return aurora.BgGray(5, aurora.Bold(content))
	case rules.Failed:
		return aurora.BgRed(aurora.Bold(content))
	case rules.Healthy:
		return aurora.Green(content)
	}

	return content
}

func (p FancyCheckResultPresenter) Present(out io.Writer) error {
	var (
		warningCheckResults  []*rules.CheckResult
		advisoryCheckResults []*rules.CheckResult
		healthyCheckResults  []*rules.CheckResult
		failedCheckResults   []*rules.CheckResult
		otherCheckResults    []*rules.CheckResult
	)
	for _, r := range p.CheckResults {
		switch r.Category {
		case rules.Healthy:
			healthyCheckResults = append(healthyCheckResults, r)
		case rules.Advisory:
			advisoryCheckResults = append(advisoryCheckResults, r)
		case rules.Warning:
			warningCheckResults = append(warningCheckResults, r)
		case rules.Failed:
			failedCheckResults = append(failedCheckResults, r)
		default:
			otherCheckResults = append(otherCheckResults, r)
		}
	}

	// table header
	fmt.Fprintln(
		out,
		aurora.Underline(aurora.Bold(p.ReportName)),
	)
	fmt.Fprintln(out)

	table := tablewriter.NewWriter(out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetColWidth(100)

	// table rows
	for _, result := range warningCheckResults {
		table.Append([]string{
			fmt.Sprint(p.formatCategory(result.Category)),
			aurora.Bold(result.RuleID).String(),
			result.Description,
		})
	}
	for _, result := range failedCheckResults {
		table.Append([]string{
			fmt.Sprint(p.formatCategory(result.Category)),
			aurora.Bold(result.RuleID).String(),
			result.Description,
		})
	}
	for _, result := range advisoryCheckResults {
		table.Append([]string{
			fmt.Sprint(p.formatCategory(result.Category)),
			aurora.Bold(result.RuleID).String(),
			result.Description,
		})
	}
	for _, result := range healthyCheckResults {
		table.Append([]string{
			fmt.Sprint(p.formatCategory(result.Category)),
			aurora.Bold(result.RuleID).String(),
			result.Description,
		})
	}

	table.Render()
	fmt.Fprintln(out)

	return nil
}
