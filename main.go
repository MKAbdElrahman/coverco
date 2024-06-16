package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mkabdelrahman/coverco/finder"
	"github.com/mkabdelrahman/coverco/printer"
	"github.com/mkabdelrahman/coverco/reporter"

	"github.com/mkabdelrahman/coverco/conf"

	"github.com/charmbracelet/log"
)

// check if config yaml is supplied and use it if exitst
// if not load the defaults
// ovveride with any command line flags supplied

func main() {
	// Define command line flags
	configFilePath := flag.String("config", "", "Path to the configuration file")
	dirPath := flag.String("dir", ".", "Path to the folder to list Go packages")
	defaultCoverageThreshold := flag.Float64("default-threshold", 80.0, "Default coverage threshold")
	coverageReportsDir := flag.String("coverage-dir", "./coverage_reports", "Directory for coverage reports")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	logFile := flag.String("log-file", "", "Log file (default: log to stdout)")
	keepReports := flag.Bool("keep-reports", false, "Keep coverage reports after printing (default: false)")

	flag.Parse()
	// Load configuration from YAML file or use defaults if config flag is not set

	// Load configuration
	config, err := loadConfiguration(*configFilePath)
	if err != nil {
		log.Errorf("Error loading config: %s", err.Error())
		return
	}

	// Override config with flags if they are set
	overrideConfigWithFlags(&config, defaultCoverageThreshold, coverageReportsDir, logLevel, logFile)

	// Setup logging
	err = setupLogging(config)
	if err != nil {
		log.Errorf("Error setting up logging: %s", err.Error())
		return
	}

	packages, err := finder.FilterCoveredPackages(config, *dirPath)
	if err != nil {
		fmt.Printf("Failed to create packages list: %v\n", err)
		return
	}

	reporter, err := reporter.NewCoverageReporter(packages, config.DefaultCoverageThreshold, config.CoverageReportsDir)
	if err != nil {
		log.Errorf("%s", err.Error())
	}

	coverages := reporter.TestPackages()

	// Print coverage results
	printer := printer.NewCoveragePrinter(reporter, os.Stdout)
	// printer.PrintCoverageResults(coverages)
	printer.PrintCoverageTable(coverages)

	if !*keepReports {
		err = os.RemoveAll(config.CoverageReportsDir)
		if err != nil {
			log.Errorf("Error removing coverage reports directory: %s", err.Error())
		}
	}
}

// setupLogging sets up logging based on the configuration
func setupLogging(cfg conf.Config) error {
	logLevel := cfg.Logging.Level
	logFile := cfg.Logging.File

	var logOutput *os.File
	if logFile != "" {
		var err error
		logOutput, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return fmt.Errorf("error setting up log file: %w", err)
		}
		log.SetOutput(logOutput)
	} else {
		logOutput = os.Stdout
	}

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		return fmt.Errorf("invalid log level: %s", logLevel)
	}

	return nil
}

// loadConfiguration loads the configuration from YAML file if provided, or returns defaults
func loadConfiguration(configFilePath string) (conf.Config, error) {
	if configFilePath != "" {
		return conf.LoadConfig(configFilePath)
	}
	return conf.GetDefaultConfig(), nil
}

// overrideConfigWithFlags overrides configuration values with command line flags
func overrideConfigWithFlags(config *conf.Config, defaultCoverageThreshold *float64, coverageReportsDir, logLevel, logFile *string) {
	config.DefaultCoverageThreshold = *defaultCoverageThreshold
	config.CoverageReportsDir = *coverageReportsDir
	config.Logging.Level = *logLevel
	config.Logging.File = *logFile

}
