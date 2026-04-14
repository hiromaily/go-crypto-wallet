---
description: Go Path Configuration Rules
globs: ["**/*.go"]
alwaysApply: false
---


# Go Path Configuration Rules

## Overview

This project uses **Homebrew-installed Go only**. Version managers (e.g., `mise`, `asdf`, `goenv`) must NOT affect this project.

## Required Go Version

The Go version is specified in `go.mod`. Always use the Homebrew-installed Go that matches this version.

## Common Symptoms of Version Manager Interference

- `compile: version "go1.X.X" does not match go tool version "go1.Y.Y"`
- `go version` shows different version than expected
- Build failures after Go upgrade

## Diagnosis

```bash
# Check which Go is being used
which go
go version
go env GOROOT

# Check for version managers in PATH (should NOT be present)
echo $PATH | tr ':' '\n' | grep -E 'mise|asdf|goenv'
```

## Solution: Disable Version Managers

If a version manager is interfering, remove it from PATH for this session:

```bash
# Remove version managers from PATH and use Homebrew Go
export PATH=$(echo $PATH | tr ':' '\n' | grep -v 'mise' | tr '\n' ':' | sed 's/:$//')
export GOROOT=/opt/homebrew/Cellar/go/1.26.2/libexec

# Verify
go version
```

## Expected Environment

| Variable     | Expected Value                              |
| ------------ | ------------------------------------------- |
| `which go`   | `/opt/homebrew/bin/go`                      |
| `go version` | Matches `go.mod` version                    |
| `GOROOT`     | `/opt/homebrew/Cellar/go/<version>/libexec` |
