### Installation

#### Prerequisites

- Node.js >= 18.0.0
- Claude Code with plugin support
- Bun and uv (auto-installed if missing)

#### Steps

```bash
# 1. In a Claude Code session, run:
/plugin marketplace add thedotmack/claude-mem

# 2. Install the plugin:
/plugin install claude-mem

# 3. Restart Claude Code
```

#### Verify Installation

After restarting, claude-mem should automatically:
- Start the worker service on port 37777
- Begin capturing session events
- Provide the `mem-search` skill

You can verify by visiting `http://localhost:37777` in your browser to see the web viewer UI.
