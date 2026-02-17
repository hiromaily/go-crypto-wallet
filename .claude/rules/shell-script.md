---
paths: ["**/*.sh", "**/Makefile", "**/make/*.mk"]
---

# Shell Script Best Practices

## Overview

Rules for writing robust, maintainable shell scripts in go-crypto-wallet.
These rules are based on code review feedback and industry best practices.

## Critical Shell Options

### Always Use Strict Mode

```bash
set -euo pipefail
```

**Explanation:**

- `-e`: Exit immediately if any command fails
- `-u`: Treat unset variables as errors
- `-o pipefail`: Fail if any command in a pipeline fails (not just the last one)

**Why pipefail matters:**

```bash
# WITHOUT pipefail (BAD)
set -eu
cat nonexistent.txt | grep "pattern"  # Only grep's exit code matters

# WITH pipefail (GOOD)
set -euo pipefail
cat nonexistent.txt | grep "pattern"  # Fails immediately if cat fails
```

## Variable Configuration

### Make Hardcoded Values Configurable

**Bad:**

```bash
# Hardcoded volume name - brittle if project name changes
docker volume rm "go-crypto-wallet_wallet-mysql"
```

**Good:**

```bash
# Configurable with default
DOCKER_VOLUME_NAME="${DOCKER_VOLUME_NAME:-go-crypto-wallet_wallet-mysql}"
docker volume rm "$DOCKER_VOLUME_NAME"
```

**Benefits:**

- Flexibility for different environments
- No breakage if project structure changes
- Easy to override in CI/CD

### Environment Variable Naming

Use descriptive, uppercase names with underscores:

```bash
RPC_USER="${RPC_USER:-xyz}"
RPC_PASSWORD="${RPC_PASSWORD:-xyz}"
WALLET_PASSPHRASE="${WALLET_PASSPHRASE:-test}"
DOCKER_VOLUME_NAME="${DOCKER_VOLUME_NAME:-go-crypto-wallet_wallet-mysql}"
```

### Use Environment Variables Instead of Modifying Config Files

**Bad (creates backups, modifies files):**

```bash
# Backup and modify config file
sed -i.bak 's|host: "127.0.0.1:18332"|host: "127.0.0.1:18332/wallet/watch"|' config.yaml

# ... operations ...

# Restore backup
if [ -f "config.yaml.bak" ]; then
    mv "config.yaml.bak" "config.yaml"
fi
```

**Good (use environment variables):**

```bash
# Create wrapper function to set environment variable per-command
watch_with_wallet() {
    WALLET_BITCOIN_HOST="127.0.0.1:18332/wallet/watch" watch "$@"
}

# Use wrapper function
watch_with_wallet -c config.yaml create payment
```

**Benefits:**

- No risk of leaving modified configs if script fails
- No backup files to manage
- Cleaner and more robust
- Easy to override in CI/CD
- No file permission issues

**When to use this pattern:**

- Applications that support environment variable overrides (check config documentation)
- Temporary config changes for scripts or tests
- Multi-environment setups (dev, staging, prod)

## Error Handling

### Refactor Duplicate Error Handling

**Bad (Duplicated):**

```bash
if echo "$output" | grep -q "error"; then
    log_error "Operation failed"
    log_error "This could indicate:"
    log_error "  - Reason 1"
    log_error "  - Reason 2"
    return 1
fi

# ... later in code ...

if echo "$output2" | grep -q "error"; then
    log_error "Operation failed"
    log_error "This could indicate:"
    log_error "  - Reason 1"
    log_error "  - Reason 2"
    return 1
fi
```

**Good (DRY with helper function):**

```bash
# Create helper function
log_operation_error() {
    log_error "Operation failed"
    log_error "This could indicate:"
    log_error "  - Reason 1"
    log_error "  - Reason 2"
    return 1
}

# Use it
if echo "$output" | grep -q "error"; then
    log_operation_error
fi

if echo "$output2" | grep -q "error"; then
    log_operation_error
fi
```

### Error Handling in Command Substitution

Always handle errors in command substitution:

```bash
# Get output with error handling
balance_json=$(bitcoin-cli getbalances 2>&1 || true)

# Parse with fallback
balance=$(echo "$balance_json" | jq -r '.amount // 0' 2>/dev/null || echo "0")
```

## Robust Comparisons

### Variable Quoting in Test Conditions (SC2086)

**Always quote variables in `[ ]` test conditions** to prevent globbing and word splitting.

**Bad (Unquoted - can fail if variable is empty or contains spaces):**

```bash
# These can fail unexpectedly
while [ $counter -lt $max ]; do
if [ $removal_attempts -lt $max_attempts ]; then
```

**Good (Properly quoted):**

```bash
# Always quote variables in test conditions
while [ "$counter" -lt "$max" ]; do
if [ "$removal_attempts" -lt "$max_attempts" ]; then
```

**Why this matters:**

- If `$counter` is empty, `[ -lt 5 ]` fails with "unary operator expected"
- If `$counter` contains spaces, it splits into multiple arguments
- Quoted variables ensure reliable behavior in all cases

**This applies to all test operators:**

```bash
# Numeric comparisons - quote both sides
[ "$a" -lt "$b" ]
[ "$a" -gt "$b" ]
[ "$a" -eq "$b" ]

# String comparisons - quote both sides
[ "$str" = "value" ]
[ "$str" != "value" ]
[ -z "$str" ]
[ -n "$str" ]
```

### Floating Point Comparisons with bc

**Bad (Fragile):**

```bash
if (($(echo "$balance > 0" | bc -l))); then
    # Fails if bc is not installed or input is malformed
fi
```

**Good (Robust):**

```bash
if [ -n "$balance" ] && [ "$(echo "$balance > 0" | bc -l 2>/dev/null || echo 0)" -eq 1 ]; then
    # Handles bc failures gracefully
fi
```

**Why this is better:**

- Checks variable is not empty first
- Redirects bc errors to /dev/null
- Defaults to 0 on failure
- Uses `-eq 1` for reliable comparison

### String Comparisons

Always quote variables:

```bash
# BAD
if [ $var = "value" ]; then

# GOOD
if [ "$var" = "value" ]; then

# EVEN BETTER (handles empty/unset)
if [ "${var:-}" = "value" ]; then
```

## Script Structure

### Standard Script Template

```bash
#!/usr/bin/env bash

# Script Description
# Usage: ./script.sh [OPTIONS]
# Options:
#   --option  Description

set -euo pipefail

# Script directory for relative paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# Source common utilities
# shellcheck source=path/to/common.sh
source "${SCRIPT_DIR}/common.sh"

# Configuration with defaults
VERBOSE="${VERBOSE:-false}"
CONFIG_FILE="${CONFIG_FILE:-config.yaml}"

###############################################################################
# Functions
###############################################################################

show_help() {
    cat <<EOF
Usage: $0 [OPTIONS]

Options:
    --help      Show this help message
    --verbose   Enable verbose output
EOF
}

main() {
    # Parse arguments
    while [ $# -gt 0 ]; do
        case "$1" in
            --help|-h)
                show_help
                exit 0
                ;;
            --verbose)
                VERBOSE=true
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done

    # Main logic here
}

# Run main
main "$@"
```

## Common Patterns

### Directory Operations

```bash
# Create directory if needed
mkdir -p "${TARGET_DIR}"

# Clean directory except .gitkeep
find "$dir" -type f ! -name '.gitkeep' -delete 2>/dev/null || true

# Remove directory contents except specific file
find "$wallet_dir" -mindepth 1 ! -name 'bitcoin.conf' -exec rm -rf {} + 2>/dev/null || true
```

### File Path Handling

```bash
# Always quote paths with spaces
cd "/path/with spaces"  # GOOD
cd /path/with spaces    # BAD - will fail

# Use double quotes for variables
cp "${SOURCE_FILE}" "${DEST_FILE}"

# Extract filename from path
filename="${path##*/}"

# Extract directory from path
dirname="${path%/*}"
```

### Loop Patterns

```bash
# Loop over arguments
for arg in "$@"; do
    echo "$arg"
done

# Loop with counter
for i in $(seq 1 10); do
    echo "$i"
done

# Loop over array
wallets=(watch keygen sign1 sign2)
for wallet in "${wallets[@]}"; do
    echo "$wallet"
done
```

## Docker Operations

### Safe Docker Volume Removal

```bash
# Try multiple times with backoff
removal_attempts=0
max_removal_attempts=5

# Note: Always quote variables in test conditions
while [ "$removal_attempts" -lt "$max_removal_attempts" ]; do
    if docker volume rm "$volume_name" 2>/dev/null; then
        log_info "Volume removed successfully"
        break
    fi
    removal_attempts=$((removal_attempts + 1))
    if [ "$removal_attempts" -lt "$max_removal_attempts" ]; then
        log_warn "Retrying in 2 seconds... (attempt $removal_attempts/$max_removal_attempts)"
        sleep 2
    fi
done
```

### Wait for Container Health

```bash
wait_for_healthy() {
    local container_name=$1
    local max_wait=${2:-60}
    local counter=0

    # Note: Always quote variables in test conditions
    while [ "$counter" -lt "$max_wait" ]; do
        status=$(docker inspect --format='{{.State.Health.Status}}' "$container_name" 2>/dev/null || echo "not_found")

        if [ "$status" = "healthy" ]; then
            return 0
        fi

        counter=$((counter + 1))
        sleep 1
    done

    return 1
}
```

## Logging Best Practices

### Use Descriptive Log Levels

```bash
log_info "Starting operation..."
log_warn "Retrying due to temporary failure"
log_error "Critical error occurred"
log_debug "Variable value: $var"  # Only shown if VERBOSE=true
```

### Log Context for Debugging

```bash
# BAD
log_error "Failed"

# GOOD
log_error "Failed to create wallet"
log_error "Container: $container_name"
log_error "Output: $output"
log_error "Exit code: $?"
```

## Security Considerations

### Never Log Sensitive Information

```bash
# BAD
log_info "Password: $PASSWORD"
log_info "Private key: $PRIVATE_KEY"

# GOOD
log_info "Credentials loaded"
log_info "Private key file: ${KEY_FILE} (not displaying contents)"
```

### Validate User Input

```bash
# Validate before use
if [ -z "$user_input" ]; then
    log_error "Input cannot be empty"
    exit 1
fi

# Sanitize paths
if [[ "$path" =~ \.\. ]]; then
    log_error "Path traversal detected"
    exit 1
fi
```

## Testing and Validation

### Check Prerequisites

```bash
check_prerequisites() {
    # Check required commands
    for cmd in docker jq bc; do
        if ! command -v "$cmd" >/dev/null 2>&1; then
            log_error "$cmd is not installed"
            exit 1
        fi
    done

    # Check required files
    for file in "$CONFIG_FILE" "$KEY_FILE"; do
        if [ ! -f "$file" ]; then
            log_error "Required file not found: $file"
            exit 1
        fi
    done
}
```

### Dry Run Mode

```bash
DRY_RUN="${DRY_RUN:-false}"

execute_command() {
    local cmd="$1"

    if [ "$DRY_RUN" = "true" ]; then
        log_info "[DRY RUN] Would execute: $cmd"
    else
        eval "$cmd"
    fi
}
```

## Common Utilities from common.sh

This project provides common utilities in `scripts/operation/common.sh`:

```bash
# Source it
source "${SCRIPT_DIR}/../../common.sh"

# Available functions
log_info "message"
log_warn "message"
log_error "message"
log_step "Major Section"
log_substep "Minor Section"

check_docker
wait_for_healthy "container-name" 60
btc_cli "btc-watch" "getblockcount"
clean_dir_except_gitkeep "data/address/btc"
```

## ShellCheck Integration

**Always run Makefile targets after modifying shell scripts:**

```bash
# Format all shell scripts
make shfmt

# Lint all shell scripts
make shellcheck
```

These targets automatically process all `.sh` files in the `scripts/` directory.

For individual file testing (optional):

```bash
# Format single file
shfmt -l -w script.sh

# Check single file
shellcheck script.sh
```

### Common ShellCheck Directives

```bash
# Disable specific warning
# shellcheck disable=SC2086
echo $var

# Mark sourced file
# shellcheck source=path/to/file.sh
source "${SCRIPT_DIR}/common.sh"

# Disable for entire file (use sparingly)
# shellcheck disable=SC1090,SC2034
```

## File Naming Conventions

- Use lowercase with hyphens: `e2e-p2pkh-2of3.sh`
- Descriptive names: `setup-bitcoin-nodes.sh` not `setup.sh`
- Prefix with purpose: `test-`, `e2e-`, `setup-`, `cleanup-`

## Documentation

### Comprehensive Header Comments

```bash
#!/usr/bin/env bash

# Bitcoin E2E Workflow Script - Pattern 2: P2PKH 2-of-3 Multisig
# This script automates the complete Bitcoin workflow
#
# Usage: ./script.sh [OPTIONS]
#
# Options:
#   --reset    Full reset and run from scratch
#   --cleanup  Stop containers and cleanup state
#   --verbose  Enable verbose output
#   -h, --help Display help message
#
# Environment Variables:
#   RPC_USER          Bitcoin RPC username (default: xyz)
#   RPC_PASSWORD      Bitcoin RPC password (default: xyz)
#
# Reference Documentation:
#   docs/chains/btc/operations/e2e-transaction-patterns.md
```

### Function Documentation

```bash
# Wait for Docker container to be healthy
# Usage: wait_for_healthy "container-name" [max_wait_seconds]
# Arguments:
#   $1 - Container name
#   $2 - Maximum wait time in seconds (default: 60)
# Returns:
#   0 on success, 1 on timeout
wait_for_healthy() {
    local container_name=$1
    local max_wait=${2:-60}
    # ... implementation ...
}
```

## Language Requirements

### Write All Comments and Messages in English

All shell scripts must use English for:

- **Header comments** - Script description, usage, options
- **Inline comments** - Code explanations
- **Function documentation** - Usage, arguments, returns
- **Log messages** - Info, warning, error messages
- **Help text** - Command-line help output

**Bad (Japanese):**

```bash
# 変換完了
log_success "変換完了!"
log_info "DRY-RUN モードでした"
```

**Good (English):**

```bash
# Conversion complete
log_success "Conversion complete!"
log_info "DRY-RUN mode was enabled"
```

**Rationale:**

- Ensures accessibility for international contributors
- Maintains consistency across the codebase
- Facilitates collaboration with global developer community
- Aligns with project documentation standards (see `documentation-language` rule)

## Makefile Integration

Shell scripts are often called from Makefiles:

```makefile
.PHONY: btc-e2e-reset
btc-e2e-reset:
 ./scripts/operation/btc/e2e/e2e-workflow.sh --reset

.PHONY: btc-e2e-verbose
btc-e2e-verbose:
 ./scripts/operation/btc/e2e/e2e-workflow.sh --verbose
```

## References

- [ShellCheck](https://www.shellcheck.net/) - Shell script analysis tool
- [Google Shell Style Guide](https://google.github.io/styleguide/shellguide.html)
- [Bash Hackers Wiki](https://wiki.bash-hackers.org/)
- Project: @scripts/operation/common.sh - Common utilities
- Project: @scripts/operation/btc/e2e/ - E2E test examples

## Quick Checklist

Before committing shell scripts:

- [ ] `set -euo pipefail` at the top
- [ ] All hardcoded values are configurable via environment variables
- [ ] Floating point comparisons use robust bc patterns
- [ ] No duplicate error handling code (use helper functions)
- [ ] All variables are quoted: `"$var"`
- [ ] Error messages include context
- [ ] No sensitive information in logs
- [ ] Prerequisites are checked
- [ ] Help message is comprehensive
- [ ] **All comments and messages are in English**
- [ ] Run `make shfmt` (format all shell scripts)
- [ ] Run `make shellcheck` (lint all shell scripts)
- [ ] File is executable: `chmod +x script.sh`
