# brightlocal-cli

Command-line interface for the BrightLocal API. Manage local SEO rankings, citation audits, location searches, and reports from your terminal.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap builtbyrobben/tap
brew install brightlocal-cli
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/builtbyrobben/brightlocal-cli/releases).

### Build from Source

```bash
git clone https://github.com/builtbyrobben/brightlocal-cli.git
cd brightlocal-cli
make build
```

## Configuration

brightlocal-cli authenticates via a BrightLocal API key. You can provide it in two ways:

**Environment variable (recommended for CI/scripts):**

```bash
export BRIGHTLOCAL_API_KEY="your-api-key"
```

**Keyring storage (recommended for interactive use):**

```bash
# Interactive prompt (secure)
brightlocal-cli auth set-key --stdin

# Pipe from environment
echo "$BRIGHTLOCAL_API_KEY" | brightlocal-cli auth set-key --stdin
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `BRIGHTLOCAL_API_KEY` | API key (overrides keyring) |
| `BRIGHTLOCAL_CLI_COLOR` | Color output: `auto`, `always`, `never` |
| `BRIGHTLOCAL_CLI_OUTPUT` | Default output mode: `json`, `plain` |

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output JSON to stdout (best for scripting) |
| `--plain` | Output stable, parseable text (TSV; no colors) |
| `--color` | Color output: `auto`, `always`, `never` |
| `--verbose` | Enable verbose logging |
| `--force` | Skip confirmations for destructive commands |
| `--no-input` | Never prompt; fail instead (useful for CI) |

## Commands

### auth

Manage authentication credentials.

```bash
# Store API key in system keyring
brightlocal-cli auth set-key --stdin

# Check authentication status
brightlocal-cli auth status

# Remove stored credentials
brightlocal-cli auth remove
```

### locations

Search for locations.

```bash
# Search for a location
brightlocal-cli locations search --query "Columbia, MO"

# Specify country (default: USA)
brightlocal-cli locations search --query "London" --country GBR

# Limit results
brightlocal-cli locations search --query "New York" --limit 5

# Output as JSON
brightlocal-cli locations search --query "Columbia, MO" --json
```

### rankings

Check local search rankings for a business.

```bash
# Check rankings for specific search terms
brightlocal-cli rankings check --business "Joe's Pizza" --location "Columbia, MO" --terms "pizza,best pizza,pizza delivery"

# Get a rankings report by ID
brightlocal-cli rankings get 12345

# Output as JSON
brightlocal-cli rankings check --business "Joe's Pizza" --location "Columbia, MO" --terms "pizza" --json
```

### citations

Run citation audits for a business.

```bash
# Run a citation audit
brightlocal-cli citations audit --business "Joe's Pizza" --location "Columbia, MO"

# Output as JSON
brightlocal-cli citations audit --business "Joe's Pizza" --location "Columbia, MO" --json
```

### reports

Manage BrightLocal reports.

```bash
# List all reports
brightlocal-cli reports list

# List with pagination
brightlocal-cli reports list --page 2 --page-size 20

# Create a new report
brightlocal-cli reports create --name "Q1 Rankings" --type rankings
brightlocal-cli reports create --name "Citation Audit" --type citations

# Output as JSON
brightlocal-cli reports list --json
```

### version

Print version information.

```bash
brightlocal-cli version
```

## License

MIT
