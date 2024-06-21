package conf

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	DefaultThreshold             = 80.0
	DefaultCoverageReportsDir    = "./coverage_reports"
	DefaultCoverageReportsFormat = "lcov"
	DefaultLoggingLevel          = "info"
	DefaultLoggingFile           = ""
	DefaultKeepReports           = true
	DefaultCoverPackageName      = "*"
)

var (
	DefaultExcludePackages = []string{}
	DefaultCoverPackages   = []struct {
		Name      string   `yaml:"name"`
		Threshold *float64 `yaml:"threshold,omitempty"`
	}{
		{Name: DefaultCoverPackageName, Threshold: nil}, // Default cover all packages
	}
)

// Config represents the configuration file structure
type Config struct {
	DefaultCoverageThreshold float64 `yaml:"default_coverage_threshold"`
	CoverageReportsDir       string  `yaml:"coverage_reports_dir"`
	CoverageReportsFormat    string  `yaml:"coverage_reports_format"`

	CoverPackages []struct {
		Name      string   `yaml:"name"`
		Threshold *float64 `yaml:"threshold,omitempty"`
	} `yaml:"cover_packages"`
	ExcludePackages []string `yaml:"exclude_packages"`
	Logging         struct {
		Level string `yaml:"level"`
		File  string `yaml:"file,omitempty"`
	} `yaml:"logging"`
	KeepReports bool `yaml:"keep_reports"`
}

// LoadConfigFromFile loads the configuration from the specified file and overlays it onto the default configuration
func LoadConfigFromFile(config *Config, configFilePath string) error {
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	var fileConfig Config
	err = yaml.Unmarshal(configData, &fileConfig)
	if err != nil {
		return fmt.Errorf("error parsing config file: %w", err)
	}

	// Overlay fileConfig onto the default config
	if fileConfig.DefaultCoverageThreshold != 0 {
		config.DefaultCoverageThreshold = fileConfig.DefaultCoverageThreshold
	}
	if fileConfig.CoverageReportsDir != "" {
		config.CoverageReportsDir = fileConfig.CoverageReportsDir
	}
	if fileConfig.CoverageReportsFormat != "" {
		config.CoverageReportsFormat = fileConfig.CoverageReportsFormat
	}
	if len(fileConfig.CoverPackages) > 0 {
		config.CoverPackages = fileConfig.CoverPackages
	}
	if len(fileConfig.ExcludePackages) > 0 {
		config.ExcludePackages = fileConfig.ExcludePackages
	}
	if fileConfig.Logging.Level != "" {
		config.Logging.Level = fileConfig.Logging.Level
	}
	if fileConfig.Logging.File != "" {
		config.Logging.File = fileConfig.Logging.File
	}
	config.KeepReports = fileConfig.KeepReports

	return nil
}

// GetDefaultConfig returns a Config struct with default values
func GetDefaultConfig() Config {
	return Config{
		DefaultCoverageThreshold: DefaultThreshold,
		CoverageReportsDir:       DefaultCoverageReportsDir,
		CoverageReportsFormat:    DefaultCoverageReportsFormat,
		CoverPackages:            DefaultCoverPackages,
		ExcludePackages:          DefaultExcludePackages,
		Logging: struct {
			Level string `yaml:"level"`
			File  string `yaml:"file,omitempty"`
		}{
			Level: DefaultLoggingLevel,
			File:  DefaultLoggingFile,
		},
		KeepReports: DefaultKeepReports,
	}
}

// OverrideWithFlags overrides configuration values with command line flags if they are set
func OverrideWithFlags(config *Config, defaultCoverageThreshold *float64, excludePatterns, coverageReportsDir, coverageReportsFormat, logLevel, logFile *string, keepReports *bool) {
	if defaultCoverageThreshold != nil && *defaultCoverageThreshold != 0 {
		config.DefaultCoverageThreshold = *defaultCoverageThreshold
	}
	if coverageReportsDir != nil && *coverageReportsDir != "" {
		config.CoverageReportsDir = *coverageReportsDir
	}
	if coverageReportsFormat != nil && *coverageReportsFormat != "" {
		config.CoverageReportsFormat = *coverageReportsFormat
	}
	if logLevel != nil && *logLevel != "" {
		config.Logging.Level = *logLevel
	}
	if logFile != nil && *logFile != "" {
		config.Logging.File = *logFile
	}
	if excludePatterns != nil && *excludePatterns != "" {
		config.ExcludePackages = append(config.ExcludePackages, strings.Split(*excludePatterns, ",")...)
	}
	if keepReports != nil {
		config.KeepReports = *keepReports
	}
}

// ExtractFinalConfig extracts the final configuration based on defaults, file, and CLI flags
func ExtractFinalConfig() (Config, error) {
	// Define command line flags
	configFilePath := flag.String("config", "", "Path to the configuration file")
	defaultCoverageThreshold := flag.Float64("default-threshold", 0, "Default coverage threshold")
	coverageReportsDir := flag.String("coverage-dir", "", "Directory for coverage reports")
	coverageReportsFormat := flag.String("coverage-reports-format", "", "Output format for coverage reports (out or lcov)")
	logLevel := flag.String("log-level", "", "Log level (debug, info, warn, error)")
	logFile := flag.String("log-file", "", "Log file (default: log to stdout)")
	keepReports := flag.Bool("keep-reports", true, "Keep coverage reports after printing (default: true)")
	excludePatterns := flag.String("exclude", "", "Comma-separated list of package patterns to exclude")

	flag.Parse()

	// Load default configuration
	config := GetDefaultConfig()

	// Load configuration from file if specified
	if *configFilePath != "" {
		err := LoadConfigFromFile(&config, *configFilePath)
		if err != nil {
			return Config{}, fmt.Errorf("error loading config from file: %w", err)
		}
	}

	// Override config with flags if they are set
	OverrideWithFlags(&config, defaultCoverageThreshold, excludePatterns, coverageReportsDir, coverageReportsFormat, logLevel, logFile, keepReports)

	return config, nil
}
