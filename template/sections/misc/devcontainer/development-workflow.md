## Development Workflow

### Daily Development

```bash
# 1. Start Docker Desktop

# 2. Open project in editor
code .  # or cursor .

# 3. Reopen in Container (first time or after config changes)

# 4. Develop normally
make check-build
make go-test
make go-lint

# 5. Docker Compose works through mounted socket
docker compose --profile mysql up -d
docker compose down

# 6. Git operations work normally
git status
git add .
git commit -m "feat: add new feature"
git push
```

### Working with Database

```bash
# Start database services
docker compose --profile mysql up -d

# Wait for database to be ready
docker compose exec db mysql -u root -proot -e "SELECT 1"

# Run migrations
make atlas-dev-reset

# Generate SQLC code
make sqlc

# Stop services when done
docker compose down
```

### Working with AI Assistant

```bash
# Inside DevContainer terminal

# Start Claude Code
claude-code

# Or use Cursor's built-in AI
# Press Cmd+K (Mac) or Ctrl+K (Windows/Linux)

# Ask AI to:
# - "Fix all linting issues"
# - "Run the tests"
# - "Create a new migration"
# - "Update the README"

# AI commands run safely inside the container
```

### Rebuilding Container

When you update `.devcontainer/devcontainer.json`:

```bash
# VS Code/Cursor:
# F1 → "Dev Containers: Rebuild Container"

# Or from terminal:
# Exit the container, then reopen
```

---
