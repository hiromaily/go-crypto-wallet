### Local Testing

Test GoReleaser locally before pushing a tag:

```bash
# Install GoReleaser (if not installed)
go install github.com/goreleaser/goreleaser/v2@latest

# Dry run (builds but doesn't publish)
make release-dry-run

# Or directly
goreleaser release --snapshot --clean --skip=publish
```
