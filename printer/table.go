package printer

import (
	"github.com/mkabdelrahman/coverco/reporter"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

// PrintCoverageTable prints the coverage data as a table
func (cp *CoveragePrinter) PrintCoverageTable(coverages []reporter.Coverage) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Package Name", "Coverage Percentage", "Threshold"})

	for _, cov := range coverages {
		packageThreshold := cp.Reporter.DefaultCoverageThreshold
		for _, pkg := range cp.Reporter.Packages {
			if pkg.Name == cov.PackageName {
				packageThreshold = pkg.Threshold
				break
			}
		}

		row := []string{
			cov.PackageName,
			fmt.Sprintf("%.2f%%", cov.Percentage),
			fmt.Sprintf("%.2f%%", packageThreshold),
		}

		if cov.Percentage < float64(packageThreshold) {
			// Set text color to red for packages that do not meet the threshold
			table.Rich(row, []tablewriter.Colors{
				{tablewriter.FgRedColor},
				{tablewriter.FgRedColor},
				{tablewriter.FgRedColor},
			})
		} else {
			// Set text color to green for packages that meet or exceed the threshold
			table.Rich(row, []tablewriter.Colors{
				{tablewriter.FgGreenColor},
				{tablewriter.FgGreenColor},
				{tablewriter.FgGreenColor},
			})
		}
	}

	table.Render()
}
