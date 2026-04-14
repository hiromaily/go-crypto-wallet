## Expected Behavior

### Always Do

- **Check current branch before starting any task** (see `git-workflow` skill)
- Read relevant documentation before making changes
- Run verification commands after code changes
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Consider impact on offline wallet operations

### Never Do

- ❌ Log private keys or sensitive information
- ❌ Edit files marked `DO NOT EDIT` (auto-generated)
- ❌ Push directly to `main` branch
- ❌ Run `git merge` or `gh pr merge`
- ❌ Run `protoc` or `buf` commands directly (always use Makefile targets like `make proto`, `make proto-ts`)

### Ask Before

- Making security-related changes
- Breaking changes to public APIs
- Changes affecting multiple layers
