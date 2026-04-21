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
