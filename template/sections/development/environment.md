## Development Environment

### DevContainer (Recommended for AI-Assisted Development)

This project provides an **optional** DevContainer configuration for a standardized, isolated development environment. DevContainer is particularly useful when working with AI coding assistants like Claude Code, GitHub Copilot, or Cursor.

**Key Benefits:**

- ✅ **Safe AI Development**: Isolated environment protects your host system from accidental AI-generated changes
- ✅ **Consistent Setup**: Pre-configured with Go 1.25.6, golangci-lint v2.8.0, Atlas v1.0.0, and GitHub CLI
- ✅ **Quick Start**: New developers can start coding in minutes
- ✅ **Zero Impact**: Local development workflow remains completely unchanged

**Quick Start:**

```bash
# 1. Open project in VS Code or Cursor
code .

# 2. Click "Reopen in Container" when prompted
# (or press F1 → "Dev Containers: Reopen in Container")

# 3. Start developing!
# All tools are pre-installed and ready
```

**Documentation:**

- 📖 [Complete DevContainer Guide](../../../docs/getting-started/devcontainer.md) - Setup, usage, and troubleshooting
- 🤖 [AI-Assisted Development with DevContainer](../../../docs/getting-started/devcontainer.md#using-with-ai-tools) - Claude Code, Copilot integration

**Note:** DevContainer is completely optional. Continue with local development if you prefer.

### Local Development

For traditional local development setup, follow the installation guide below.
