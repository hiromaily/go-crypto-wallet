### Changelog Generation

Changelog is automatically generated from commit messages using Conventional Commits:

| Commit Type | Changelog Section |
|-------------|-------------------|
| `feat` | 🚀 Features |
| `fix` | 🐛 Bug Fixes |
| `refactor` | 🔄 Refactoring |
| `deps`, `build` | 📦 Dependencies |
| Other | Others |

Excluded from changelog:

- `docs:` commits
- `test:` commits
- `chore:` commits
- `ci:` commits
- Merge commits
