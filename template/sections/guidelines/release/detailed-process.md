### Detailed Process

#### Step 1: Prepare for Release

Ensure all changes are merged to `main`:

```bash
git checkout main
git pull origin main
```

Verify there are no uncommitted changes:

```bash
git status
```

#### Step 2: Determine Version

This project follows [Semantic Versioning](https://semver.org/):

| Change Type | Version Bump | Example |
|-------------|--------------|---------|
| Breaking changes | MAJOR | v5.0.0 → v6.0.0 |
| New features (backward compatible) | MINOR | v6.0.0 → v6.1.0 |
| Bug fixes | PATCH | v6.1.0 → v6.1.1 |

Check current version and commits since last release:

```bash
make release-check
```

#### Step 3: Create Release Tag

Create an annotated tag and push:

```bash
# Using Makefile (recommended)
make release-tag VERSION=v6.1.0

# Or manually
git tag -a v6.1.0 -m "Release v6.1.0"
git push origin v6.1.0
```

#### Step 4: Monitor Release Workflow

The GitHub Actions workflow will automatically run:

1. **Check workflow status**:

   ```bash
   gh run list --workflow=release.yml --limit=3
   ```

2. **Watch live progress**:

   ```bash
   gh run watch <run-id>
   ```

3. **View on GitHub**:
   <https://github.com/hiromaily/go-crypto-wallet/actions/workflows/release.yml>

#### Step 5: Verify Release

Once the workflow completes:

1. **Check release page**:
   <https://github.com/hiromaily/go-crypto-wallet/releases>

2. **Verify artifacts**:

   ```bash
   gh release view v6.1.0
   ```
