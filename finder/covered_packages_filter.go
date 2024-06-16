package finder

import (
	"github.com/mkabdelrahman/coverco/conf"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
)

type Package struct {
	Name      string
	Threshold float64
}

// PatternMatchError provides detailed information about pattern matching errors.
type PatternMatchError struct {
	Pattern string
	Err     error
}

func (e *PatternMatchError) Error() string {
	return fmt.Sprintf("pattern '%s': %v", e.Pattern, e.Err)
}

func (e *PatternMatchError) Unwrap() error {
	return e.Err
}

// PackageFilter manages the filtering of packages based on configuration.
type packageFilter struct {
	config       conf.Config
	allPackages  []string
	matchedPkgs  []Package
	excludedPkgs []Package
}

// NewPackageFilter creates a new PackageFilter instance.
func newPackageFilter(config conf.Config, allPackages []string) *packageFilter {
	return &packageFilter{
		config:      config,
		allPackages: allPackages,
	}
}

// MatchPackages matches packages based on the cover patterns specified in the configuration.
func (pf *packageFilter) matchPackages() error {
	for _, coverPattern := range pf.config.CoverPackages {
		found := false
		for _, pkg := range pf.allPackages {
			matched, err := matchPattern(pkg, []string{coverPattern.Name})
			if err != nil {
				return err
			}
			if matched {
				threshold := pf.config.DefaultCoverageThreshold
				if coverPattern.Threshold != nil {
					threshold = *coverPattern.Threshold
				}
				pf.matchedPkgs = append(pf.matchedPkgs, Package{Name: pkg, Threshold: threshold})
				found = true
			}
		}
		if !found {
			log.Warnf("No packages found matching cover pattern: %s", coverPattern.Name)
		}
	}
	return nil
}

// ExcludePackages excludes packages based on the patterns specified in the configuration.
func (pf *packageFilter) excludePackages() error {
	for _, pkg := range pf.matchedPkgs {
		excluded, err := matchPattern(pkg.Name, pf.config.ExcludePackages)
		if err != nil {
			return err
		}
		if excluded {
			pf.excludedPkgs = append(pf.excludedPkgs, pkg)
			log.Infof("Excluding package: %s", pkg.Name)
		}
	}
	return nil
}

// GetFilteredPackages returns the final list of packages after applying cover and exclude filters.
func (pf *packageFilter) getFilteredPackages() []Package {
	excludedMap := make(map[string]bool)
	for _, pkg := range pf.excludedPkgs {
		excludedMap[pkg.Name] = true
	}

	var filteredPackages []Package
	for _, pkg := range pf.matchedPkgs {
		if !excludedMap[pkg.Name] {
			filteredPackages = append(filteredPackages, pkg)
		}
	}
	return filteredPackages
}

// FilterPackages creates and applies filters to return the final list of packages based on the configuration.
func FilterCoveredPackages(cfg conf.Config, dirPath string) ([]Package, error) {

	// List Go packages in the specified folder
	allPackages, err := ListGoPackages(dirPath)
	if err != nil {
		return nil, fmt.Errorf("error listing packages: %w", err)
	}

	pf := newPackageFilter(cfg, allPackages)

	if err := pf.matchPackages(); err != nil {
		return nil, fmt.Errorf("error matching packages: %w", err)
	}

	if err := pf.excludePackages(); err != nil {
		return nil, fmt.Errorf("error excluding packages: %w", err)
	}

	return pf.getFilteredPackages(), nil
}

// Supported Patterns:
//   - "github.com/example/*": Matches any package path starting with "github.com/example/"
//   - "github.com/*/project": Matches any package path starting with "github.com/", followed by any single segment, and ending with "/project"
//   - "github.com/example/project/pkg": Matches exactly "github.com/example/project/pkg"
func matchPattern(packageName string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		// Convert pattern to regex
		regexPattern := "^" + strings.ReplaceAll(pattern, "*", ".*") + "$"
		matched, err := regexp.MatchString(regexPattern, packageName)
		if err != nil {
			return false, &PatternMatchError{Pattern: pattern, Err: err}
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

// listGoPackages lists all Go packages in the specified folder and its subdirectories
func ListGoPackages(folder string) ([]string, error) {
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = folder
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error listing Go packages in folder %s: %w", folder, err)
	}
	packages := strings.Fields(string(output))
	return packages, nil
}
