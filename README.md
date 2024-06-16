## Coverco

### Purpose

Coverco is a tool designed to analyze the test coverage percentage of each package within a Go module. It allows setting specific coverage thresholds per package, excludes certain packages from analysis, provides configurable logging options, and supports command-line overrides for enhanced flexibility. Coverco is designed to fine-tune the coverage percentages for core business logic. It is meant to be part of a CI/CD pipeline or to provide developers with a useful tool to see how much of their code is tested and gain confidence in their tests. Additionally, it helps ignore some noisy coverage that downgrades the overall coverage and has no comparable value to the core logic.

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
   - `-exclude`: Comma-separated list of package patterns to exclude (e.g., `-exclude=demo/exclude/*,demo/skip/*`).

5. **Configuration Priority**:
   - Command-line flags have the highest priority.
   - YAML configuration file values have higher priority than defaults.
   - Internal defaults are used if neither flags nor configuration file values are provided.


### Quick Example

```bash
❯ git clone https://github.com/fyne-io/fyne.git
❯ git cd fyne

❯ coverco  -exclude  "*/cmd/*,*/driver/*,*/app" -default-threshold 70
2024/06/17 01:06:28 INFO Running go mod tidy...
2024/06/17 01:06:28 INFO Excluding package: fyne.io/fyne/v2/app
2024/06/17 01:06:28 INFO Excluding package: fyne.io/fyne/
.
.
.

+-------------------------------------------+---------------------+-----------+
|               PACKAGE NAME                | COVERAGE PERCENTAGE | THRESHOLD |
+-------------------------------------------+---------------------+-----------+
| fyne.io/fyne/v2                           | 65.80%              | 70.00%    |
| fyne.io/fyne/v2/canvas                    | 59.80%              | 70.00%    |
| fyne.io/fyne/v2/container                 | 86.20%              | 70.00%    |
| fyne.io/fyne/v2/data/binding              | 43.40%              | 70.00%    |
| fyne.io/fyne/v2/data/validation           | 88.20%              | 70.00%    |
| fyne.io/fyne/v2/dialog                    | 79.70%              | 70.00%    |
| fyne.io/fyne/v2/driver                    | 100.00%             | 70.00%    |
| fyne.io/fyne/v2/internal                  | 69.40%              | 70.00%    |
| fyne.io/fyne/v2/internal/animation        | 77.90%              | 70.00%    |
| fyne.io/fyne/v2/internal/async            | 54.50%              | 70.00%    |
| fyne.io/fyne/v2/internal/cache            | 61.60%              | 70.00%    |
| fyne.io/fyne/v2/internal/color            | 16.40%              | 70.00%    |
| fyne.io/fyne/v2/internal/driver           | 90.80%              | 70.00%    |
| fyne.io/fyne/v2/internal/painter          | 52.10%              | 70.00%    |
| fyne.io/fyne/v2/internal/painter/gl       | 0.00%               | 70.00%    |
| fyne.io/fyne/v2/internal/painter/software | 92.10%              | 70.00%    |
| fyne.io/fyne/v2/internal/repository       | 86.10%              | 70.00%    |
| fyne.io/fyne/v2/internal/repository/mime  | 100.00%             | 70.00%    |
| fyne.io/fyne/v2/internal/scale            | 0.00%               | 70.00%    |
| fyne.io/fyne/v2/internal/svg              | 60.60%              | 70.00%    |
| fyne.io/fyne/v2/internal/test             | 68.20%              | 70.00%    |
| fyne.io/fyne/v2/internal/widget           | 88.80%              | 70.00%    |
| fyne.io/fyne/v2/layout                    | 94.50%              | 70.00%    |
| fyne.io/fyne/v2/storage                   | 62.00%              | 70.00%    |
| fyne.io/fyne/v2/storage/repository        | 47.40%              | 70.00%    |
| fyne.io/fyne/v2/test                      | 63.80%              | 70.00%    |
| fyne.io/fyne/v2/theme                     | 57.10%              | 70.00%    |
| fyne.io/fyne/v2/tools/playground          | 0.00%               | 70.00%    |
| fyne.io/fyne/v2/widget                    | 92.20%              | 70.00%    |
+-------------------------------------------+---------------------+-----------+
```


### Contributions

Contributions are welcome! Please fork the repository and submit pull requests for any improvements you'd like to make.
