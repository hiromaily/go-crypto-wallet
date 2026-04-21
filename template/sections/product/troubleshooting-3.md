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
