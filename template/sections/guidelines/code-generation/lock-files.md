### Dependency Lock Files

**Tool**: Go modules, npm/yarn

**Generated Files**:

- `go.sum` - Go module checksums
- `web/*/yarn.lock` - Yarn package lock files
- `web/*/package-lock.json` - npm package lock files

**Note**: These files track exact dependency versions and should be committed to version control.
