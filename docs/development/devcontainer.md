# DevContainer Development Environment

## Table of Contents

- [Overview](#overview)
- [Why DevContainer?](#why-devcontainer)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Using with AI Tools](#using-with-ai-tools)
- [Features](#features)
- [Development Workflow](#development-workflow)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Overview

This project provides an optional DevContainer configuration that creates a standardized, isolated development environment. DevContainer allows you to develop inside a Docker container with all necessary tools pre-installed, ensuring consistency across different machines and team members.

**Key Point:** DevContainer is **completely optional**. You can continue using your local development environment without any changes to the existing workflow.

---

## Why DevContainer?

### Benefits for All Developers

1. **Consistent Environment**: Everyone uses the same Go version, tools, and dependencies
2. **Isolated Development**: Project dependencies don't affect your host system
3. **Quick Setup**: New team members can start coding in minutes
4. **Clean Uninstall**: Remove the container to completely clean up the development environment

### Benefits for AI-Assisted Development

**DevContainer is particularly valuable when working with AI coding assistants like Claude Code, GitHub Copilot, or other AI tools:**

1. **Safety and Isolation**
   - AI tools may occasionally generate code that modifies or deletes files
   - DevContainer provides a sandboxed environment, protecting your host system
   - If something goes wrong, simply rebuild the container

2. **Consistent AI Context**
   - All team members and AI tools work with the same environment
   - AI-generated code works consistently across different machines
   - Reduces "works on my machine" issues with AI-suggested solutions

3. **Easy Rollback**
   - Quickly reset to a clean state if AI changes cause issues
   - Container snapshots allow experimentation without risk
   - Git combined with container isolation provides double safety

4. **Reproducible AI Sessions**
   - AI tools can rely on specific tool versions
   - Build tags and linter configurations are pre-configured
   - Reduces friction when AI suggests commands or configurations

---

## Prerequisites

### Required Software

1. **Docker Desktop** (or Docker Engine)
   - [Download Docker Desktop](https://www.docker.com/products/docker-desktop/)
   - Ensure Docker is running before starting

2. **VS Code** or **Cursor**
   - [VS Code Download](https://code.visualstudio.com/)
   - [Cursor Download](https://cursor.sh/)

3. **Dev Containers Extension** (for VS Code)
   - Install from: [Dev Containers Extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
   - Cursor has built-in DevContainer support

### System Requirements

- **macOS**: macOS 10.15 or later
- **Windows**: Windows 10/11 with WSL 2
- **Linux**: Any modern distribution with Docker support
- **RAM**: 8GB minimum, 16GB recommended
- **Disk**: 10GB free space for Docker images and containers

---

## Getting Started

### Option 1: Using VS Code

1. **Open the Project**

   ```bash
   cd go-crypto-wallet
   code .
   ```

2. **Reopen in Container**
   - VS Code will detect the `.devcontainer` configuration
   - Click "Reopen in Container" when prompted
   - **OR** press `F1` → "Dev Containers: Reopen in Container"

3. **Wait for Setup**
   - First time: 5-10 minutes (downloads image, installs tools)
   - Subsequent times: 30-60 seconds (uses cached image)

4. **Start Developing**
   - Terminal opens inside the container
   - All tools are pre-installed and ready to use

### Option 2: Using Cursor

1. **Open the Project**

   ```bash
   cd go-crypto-wallet
   cursor .
   ```

2. **Reopen in Container**
   - Cursor automatically detects DevContainer configuration
   - Click "Reopen in Container" when prompted
   - **OR** use Command Palette → "Reopen in Container"

3. **Start Developing**
   - Cursor AI works seamlessly inside the container
   - All project tools are available

### Option 3: Using Claude Code (CLI)

1. **Navigate to Project**

   ```bash
   cd go-crypto-wallet
   ```

2. **Start Claude Code**

   ```bash
   claude-code
   ```

3. **Claude Code will:**
   - Detect the DevContainer configuration
   - Ask if you want to use it (optional)
   - Work inside the container if you choose

4. **Benefits with Claude Code:**
   - Claude can safely execute commands inside the container
   - All tools (make, golangci-lint, Atlas) are available
   - Git operations work seamlessly
   - Docker Compose commands work through host socket

---

## Using with AI Tools

### Claude Code (Recommended)

**Why DevContainer with Claude Code?**

- Claude can execute terminal commands safely inside the container
- All verification commands (`make go-lint`, `make check-build`) work out of the box
- Database operations via Docker Compose work through the mounted socket
- Your host system stays clean and protected

**Example AI-Assisted Workflow:**

```bash
# 1. Open project in DevContainer
cursor .  # or code .

# 2. Reopen in Container when prompted

# 3. Start Claude Code inside the container
# Terminal is already inside the container
claude-code

# 4. Ask Claude to make changes
# Claude: "Let me fix the linting issues"
# Claude runs: make go-lint
# All tools work perfectly!

# 5. Claude can verify builds
# Claude runs: make check-build
# Go 1.25.6 matches exactly what's in go.mod

# 6. Claude can run tests
# Claude runs: make gotest
# Integration tests work with proper build tags
```

**Safety Features:**

- If Claude accidentally breaks something, just rebuild the container
- Host filesystem is protected
- Git repository is mounted, so commits are preserved
- Docker containers (Bitcoin, MySQL) run on host, accessible from container

### GitHub Copilot

**Using Copilot in DevContainer:**

1. Install GitHub Copilot extension in VS Code/Cursor
2. Copilot works normally inside the container
3. Code suggestions are based on the container's Go environment
4. All suggested commands will work with pre-installed tools

### Other AI Tools

**Compatible with:**

- Tabnine
- Codeium
- Amazon CodeWhisperer
- Any VS Code/Cursor extension that works with Remote Development

---

## Features

### Pre-Installed Tools

The DevContainer comes with all project tools pre-configured:

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.25.6 | Programming language |
| golangci-lint | v2.8.0 | Code linting |
| Atlas | v1.0.0 | Database migrations |
| GitHub CLI | Latest | GitHub operations |
| Docker | Host | Docker Compose support |
| Git | Latest | Version control |

### Pre-Configured Settings

All VS Code settings from `.vscode/settings.json` are applied:

- ✅ Go build tags: `integration`
- ✅ golangci-lint integration
- ✅ Format on save
- ✅ Organize imports on save
- ✅ Markdown linting
- ✅ Shell script checking
- ✅ YAML validation
- ✅ And more...

### Pre-Installed Extensions

All VS Code extensions from `.vscode/extensions.json` are automatically installed:

- Go extension
- Markdown linting
- Shell script formatting
- YAML/TOML/JSON support
- Docker support
- GitLens
- And more...

---

## Development Workflow

### Daily Development

```bash
# 1. Start Docker Desktop

# 2. Open project in editor
code .  # or cursor .

# 3. Reopen in Container (first time or after config changes)

# 4. Develop normally
make check-build
make gotest
make go-lint

# 5. Docker Compose works through mounted socket
docker compose up -d
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
docker compose up -d

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

## Troubleshooting

### Container Fails to Start

**Problem:** Container build fails or won't start

**Solutions:**

```bash
# 1. Check Docker is running
docker ps

# 2. Rebuild container without cache
# F1 → "Dev Containers: Rebuild Container Without Cache"

# 3. Check Docker Desktop has enough resources
# Docker Desktop → Preferences → Resources
# Recommended: 4GB RAM, 2 CPUs minimum

# 4. Remove all containers and start fresh
docker system prune -a
# Then rebuild container
```

### Tools Not Found

**Problem:** `golangci-lint` or `atlas` command not found

**Solutions:**

```bash
# 1. Check if postCreateCommand ran
cat /tmp/devcontainer-setup.log

# 2. Manually install tools
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0
curl -sSfL https://atlasgo.sh | sh -s -- --version v1.0.0

# 3. Add tools to PATH
export PATH=$PATH:$HOME/go/bin:$HOME/.local/bin

# 4. Rebuild container
# F1 → "Dev Containers: Rebuild Container"
```

### Docker Compose Not Working

**Problem:** `docker compose` commands fail

**Solutions:**

```bash
# 1. Verify Docker socket is mounted
ls -la /var/run/docker.sock

# 2. Check Docker host is accessible
docker ps

# 3. Ensure docker-outside-of-docker feature is enabled
# Check .devcontainer/devcontainer.json has:
# "features": {
#   "ghcr.io/devcontainers/features/docker-outside-of-docker:1": {}
# }

# 4. Restart Docker Desktop and rebuild container
```

### Performance Issues

**Problem:** Container is slow

**Solutions:**

```bash
# 1. Allocate more resources to Docker Desktop
# Docker Desktop → Preferences → Resources
# Increase CPUs to 4, RAM to 8GB

# 2. Use volume mounts instead of bind mounts
# (Already configured in .devcontainer/devcontainer.json)

# 3. Exclude large directories from sync
# Add to .dockerignore:
# node_modules/
# .git/
# data/

# 4. For macOS: Enable VirtioFS
# Docker Desktop → Settings → General
# Enable "Use the new Virtualization framework"
```

### Git Authentication Issues

**Problem:** Can't push to GitHub from container

**Solutions:**

```bash
# 1. SSH keys are automatically forwarded
# No action needed if you use SSH

# 2. For HTTPS, use GitHub CLI
gh auth login

# 3. Or configure Git credential helper
git config --global credential.helper store

# 4. Use VS Code's built-in Git authentication
# Git operations in VS Code UI handle auth automatically
```

### Extension Not Working

**Problem:** VS Code extension doesn't work in container

**Solutions:**

```bash
# 1. Check extension supports remote development
# Look for "Supports Remote Development" badge

# 2. Manually install extension in container
# Extensions panel → Install in Dev Container

# 3. Add extension to devcontainer.json
# Edit .devcontainer/devcontainer.json:
# "customizations": {
#   "vscode": {
#     "extensions": ["publisher.extension-name"]
#   }
# }

# 4. Rebuild container to apply changes
```

---

## FAQ

### Is DevContainer required?

**No.** DevContainer is completely optional. You can continue with local development without any changes.

### Does DevContainer affect my host system?

**No.** Everything runs in an isolated container. Your host system remains unchanged.

### Can I use DevContainer with multiple projects?

**Yes.** Each project has its own container with its own configuration.

### What happens to my files?

Your files are **bind-mounted** from your host system. Changes in the container are reflected on your host and vice versa.

### Can I customize the DevContainer?

**Yes.** Edit `.devcontainer/devcontainer.json`:

- Change the base image
- Add more tools in `postCreateCommand`
- Add/remove VS Code extensions
- Modify settings

### How much disk space does it use?

- Base Go image: ~1GB
- Tools and dependencies: ~500MB
- **Total: ~1.5GB** for the first setup
- Subsequent projects share the base image

### Can I use DevContainer offline?

**Partially.** After the first build, you can work offline. However:

- Initial build requires internet (downloads image)
- `go get` requires internet
- `docker compose pull` requires internet

### Does DevContainer slow down my editor?

**No.** VS Code/Cursor connects to the container via a lightweight client-server architecture. Performance is similar to local development.

### Can I access localhost from the container?

**Yes.** Ports are automatically forwarded. If your app runs on `localhost:8080` in the container, access it at `localhost:8080` on your host.

### What about M1/M2/M3 Mac compatibility?

**Fully supported.** The DevContainer images support ARM64 architecture.

### Can I use Docker Compose in DevContainer?

**Yes.** Docker Compose commands work through the mounted host socket. Containers run on the host, accessible from DevContainer.

### How do I exit DevContainer?

```bash
# Close the VS Code/Cursor window
# OR
# F1 → "Dev Containers: Reopen Folder Locally"
```

### How do I delete the DevContainer?

```bash
# Remove container and image
docker ps -a | grep go-crypto-wallet | awk '{print $1}' | xargs docker rm
docker images | grep go-crypto-wallet | awk '{print $3}' | xargs docker rmi

# Or use Docker Desktop UI
# Containers → Delete
# Images → Delete
```

---

## Comparison: DevContainer vs Local Development

| Aspect | DevContainer | Local Development |
|--------|--------------|-------------------|
| **Setup Time** | 5-10 min (first time), 30s (subsequent) | 30-60 min (install all tools) |
| **Consistency** | ✅ Guaranteed same environment | ⚠️ Depends on individual setup |
| **Isolation** | ✅ Fully isolated | ❌ Affects host system |
| **AI Safety** | ✅ Sandboxed, safe for AI tools | ⚠️ AI can affect host directly |
| **Tool Versions** | ✅ Always correct versions | ⚠️ May drift over time |
| **Disk Space** | ~1.5GB per project | ~500MB per project |
| **Performance** | ~95% of native | 100% native |
| **Offline Work** | ⚠️ Limited (after first build) | ✅ Full offline support |
| **Complexity** | Requires Docker knowledge | Familiar workflow |
| **Onboarding** | ✅ Minutes for new developers | ⚠️ Hours for new developers |

---

## Additional Resources

### Official Documentation

- [DevContainers Specification](https://containers.dev/)
- [VS Code Remote Development](https://code.visualstudio.com/docs/remote/remote-overview)
- [Docker Documentation](https://docs.docker.com/)

### Project-Specific Documentation

- [Installation Guide](../Installation.md)
- [Database Development](database.md)
- [AI Agent Skills](../ai-agents/agent-skills.md)
- [Architecture Overview](../architecture/)

### Getting Help

- **Project Issues**: [GitHub Issues](https://github.com/hiromaily/go-crypto-wallet/issues)
- **DevContainer Issues**: Include `[DevContainer]` in issue title
- **VS Code Remote**: [VS Code Remote GitHub](https://github.com/microsoft/vscode-remote-release)

---

## Contributing

If you improve the DevContainer configuration:

1. Test your changes thoroughly
2. Update this documentation
3. Submit a pull request
4. Include before/after comparisons

---

**Last Updated**: January 2026
**DevContainer Version**: 1.0.0
**Maintained By**: go-crypto-wallet team
