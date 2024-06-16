## Coverco

Coverco is a tool designed to analyze the test coverage percentage of each package within a Go module. It allows setting specific coverage thresholds per package, excludes certain packages from analysis, provides configurable logging options, and supports command-line overrides for enhanced flexibility.

### Features

- **Thresholds**: Define coverage thresholds for individual packages or patterns.
- **Exclusions**: Exclude specific packages from coverage analysis.
- **Logging**: Configure logging levels and optionally log to a file.
- **Command-Line Flags**: Override default configurations using command-line flags.

### Configuration

Configuration is managed via a YAML file (`config.yaml` by default) with the following structure:

```yaml
# Default coverage threshold applied to all packages not explicitly listed
default_coverage_threshold: 80.0

# Directory to save coverage reports
coverage_reports_dir: "coverage_reports"

# List of package coverage configurations
cover_packages:
  # Pattern to match package and its specific threshold
  - name: "demo/arrays"
    threshold: 95.0
  - name: "demo/services/*"
    threshold: 90.0
  - name: "demo/utility/*"
    threshold: 85.0
  - name: "*"  # Default pattern covering all packages
    threshold: 80.0

# Patterns to exclude specific packages
exclude_packages:
  - "demo/exclude/*"
  - "demo/skip/*"

# Logging configuration
logging:
  level: "info"       # Can be "debug", "info", "warn", "error"
  file: "coverage.log" # Log file path (optional)
```

### Usage

1. **Install Coverco**: Install the `Coverco` tool.

   ```sh
   go install github.com/mkabdelrahman/coverco
   ```

2. **Create Configuration File**: Create a `config.yaml` file in your project root with desired configurations. Example configurations are provided above.

3. **Run Coverco**:

   ```sh
   coverco -config=config.yaml
   ```

   - Replace `config.yaml` with your actual configuration file path. If no `-config` flag is provided, Coverco will use internal defaults.

4. **Command-Line Flags**:
   - `-config`: Path to the configuration file (default: `config.yaml`).
   - `-dir`: Path to the folder to list Go packages.
   - `-default-threshold`: Default coverage threshold (default: `80.0`).
   - `-coverage-dir`: Directory for coverage reports (default: `./coverage_reports`).
   - `-log-level`: Log level (`debug`, `info`, `warn`, `error`; default: `info`).
   - `-log-file`: Log file path (default: log to stdout).
   - `-keep-reports`: Keep coverage reports after printing (default: `false`).

### Contributions

Contributions are welcome! Please fork the repository and submit pull requests for any improvements you'd like to make.

