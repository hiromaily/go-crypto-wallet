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
# Go 1.26.2 matches exactly what's in go.mod

# 6. Claude can run tests
# Claude runs: make go-test
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
