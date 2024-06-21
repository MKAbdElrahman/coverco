package reporter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/mkabdelrahman/coverco/finder"
)

var (
	ErrCoveragePercentageNotFound = fmt.Errorf("coverage percentage not found")
)

// Coverage represents the coverage information for a package
type Coverage struct {
	PackageName  string
	Percentage   float64
	CoverageFile string
}

// CoverageReporter represents a coverage reporter
type CoverageReporter struct {
	Packages                 []finder.Package
	DefaultCoverageThreshold float64
	ReportsDir               string
	OutputFormat             string
}

// NewCoverageReporter creates a new CoverageReporter instance
func NewCoverageReporter(packages []finder.Package, defaultThreshold float64, reportsDir, outputFormat string) (*CoverageReporter, error) {

	// Ensure the coverage reports directory exists
	err := ensureDir(reportsDir)
	if err != nil {
		return nil, fmt.Errorf("error ensuring coverage reports directory: %s", err.Error())
	}

	return &CoverageReporter{
		Packages:                 packages,
		DefaultCoverageThreshold: defaultThreshold,
		ReportsDir:               reportsDir,
		OutputFormat:             outputFormat,
	}, nil
}

// TestPackages tests all packages and returns their coverage information
func (cr *CoverageReporter) TestPackages() []Coverage {
	var coverages []Coverage
	for _, pkg := range cr.Packages {
		coverage := cr.testSinglePackage(pkg)
		coverages = append(coverages, coverage)
	}
	return coverages
}

// testSinglePackage tests a single package and returns its coverage information
func (cr *CoverageReporter) testSinglePackage(pkg finder.Package) Coverage {
	log.Infof("Testing package: %s", pkg.Name)

	coverProfileName := filepath.Join(cr.ReportsDir, fmt.Sprintf("coverage_%s.out", strings.ReplaceAll(pkg.Name, "/", "_")))
	cmd := exec.Command("go", "test", "-coverprofile="+coverProfileName, pkg.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("Error testing package %s: %s", pkg.Name, err.Error())
		return Coverage{PackageName: pkg.Name, Percentage: 0}
	}

	coverage, err := extractCoveragePercentage(output)
	if err != nil {
		log.Warnf("Failed to extract coverage percentage for package %s error: %s", pkg.Name, err)
		return Coverage{PackageName: pkg.Name, Percentage: 0}
	}

	if cr.OutputFormat == "lcov" {
		lcovFile := strings.Replace(coverProfileName, ".out", ".lcov", 1)
		err = convertToLcov(coverProfileName, lcovFile)
		if err != nil {
			log.Errorf("Error converting coverage profile to lcov for package %s: %s", pkg.Name, err.Error())
			return Coverage{PackageName: pkg.Name, Percentage: coverage, CoverageFile: coverProfileName}
		}
		return Coverage{PackageName: pkg.Name, Percentage: coverage, CoverageFile: lcovFile}
	}

	return Coverage{PackageName: pkg.Name, Percentage: coverage, CoverageFile: coverProfileName}
}

// convertToLcov converts a Go coverage profile to lcov format using gcov2lcov
func convertToLcov(inputFile, outputFile string) error {
	cmd := exec.Command("gcov2lcov", "-infile", inputFile, "-outfile", outputFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to convert to lcov: %w", err)
	}
	return nil
}

// extractCoveragePercentage extracts the coverage percentage from the command output
func extractCoveragePercentage(output []byte) (float64, error) {
	regex := regexp.MustCompile(`coverage: ([0-9.]+)% of statements`)
	match := regex.FindStringSubmatch(string(output))
	if len(match) < 2 {
		return 0, ErrCoveragePercentageNotFound
	}

	return strconv.ParseFloat(match[1], 64)
}

// ensureDir ensures the specified directory exists, creating it if necessary
func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Infof("Creating directory: %s", dir)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}
