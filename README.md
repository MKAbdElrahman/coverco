## Coverco

### Purpose

Coverco is designed to fine-tune the coverage percentages for core business logic. It is meant to be part of a CI/CD pipeline or to provide developers with a useful tool to see how much of their code is tested. Additionally, it helps ignore some noisy coverage that downgrades the overall coverage and has no comparable value to the core logic.

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
   coverco [flags...] [dir]
   ```

   - Replace `config.yaml` with your actual configuration file path. If no `-config` flag is provided, Coverco will use internal defaults.
   - `[dir]` is the path to the folder to list Go packages, with a default value of `.`.

4. **Command-Line Flags**:
   - `-config`: Path to the configuration file (optional; default: `config.yaml`).
   - `-default-threshold`: Default coverage threshold (default: `80.0`).
   - `-coverage-dir`: Directory for coverage reports (default: `./coverage_reports`).
   - `-log-level`: Log level (`debug`, `info`, `warn`, `error`; default: `info`).
   - `-log-file`: Log file path (default: log to stdout).
   - `-keep-reports`: Keep coverage reports after printing (default: `false`).

5. **Configuration Priority**:
   - Command-line flags have the highest priority.
   - YAML configuration file values have higher priority than defaults.
   - Internal defaults are used if neither flags nor configuration file values are provided.



### Contributions

Contributions are welcome! Please fork the repository and submit pull requests for any improvements you'd like to make.
