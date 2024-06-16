package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the YAML configuration file structure
type Config struct {
	DefaultCoverageThreshold float64 `yaml:"default_coverage_threshold"`
	CoverageReportsDir       string  `yaml:"coverage_reports_dir"`
	CoverPackages            []struct {
		Name      string   `yaml:"name"`
		Threshold *float64 `yaml:"threshold,omitempty"`
	} `yaml:"cover_packages"`
	ExcludePackages []string `yaml:"exclude_packages"`
	Logging         struct {
		Level string `yaml:"level"`
		File  string `yaml:"file,omitempty"`
	} `yaml:"logging"`
}

// loadConfig loads the configuration from the specified file
func LoadConfig(configFilePath string) (Config, error) {
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}
