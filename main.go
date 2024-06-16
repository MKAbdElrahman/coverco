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

func main() {
	log.SetLevel(log.DebugLevel)

	// Define a command line flag for the configuration file path
	configFilePath := flag.String("config", "config.yaml", "Path to the configuration file")
	dirPath := flag.String("dir", ".", "Path to the folder to list Go packages")

	flag.Parse()

	// Load configuration from YAML file
	config, err := conf.LoadConfig(*configFilePath)
	if err != nil {
		log.Errorf("Error loading config: %s", err.Error())
		return
	}

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

	// Remove the coverage reports directory
	err = os.RemoveAll(config.CoverageReportsDir)
	if err != nil {
		log.Errorf("Error removing coverage reports directory: %s", err.Error())
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
