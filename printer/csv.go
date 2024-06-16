package printer

import (
	"github.com/mkabdelrahman/coverco/reporter"
	"encoding/csv"
	"fmt"
)

func (cp *CoveragePrinter) PrintCoverageCSV(coverages []reporter.Coverage) error {
	writer := csv.NewWriter(cp.Output)

	// Write CSV header
	if err := writer.Write([]string{"Package Name", "Coverage Percentage", "Threshold"}); err != nil {
		return err
	}

	// Write CSV rows
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

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}
