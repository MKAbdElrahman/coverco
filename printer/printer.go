package printer

import (
	"github.com/mkabdelrahman/coverco/reporter"
	"io"
)

// CoveragePrinter handles printing of coverage results
type CoveragePrinter struct {
	Reporter *reporter.CoverageReporter
	Output   io.Writer
}

// NewCoveragePrinter creates a new CoveragePrinter instance
func NewCoveragePrinter(reporter *reporter.CoverageReporter, output io.Writer) *CoveragePrinter {
	return &CoveragePrinter{
		Reporter: reporter,
		Output:   output,
	}
}
