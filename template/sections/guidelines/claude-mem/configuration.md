### Configuration

Settings are stored in `~/.claude-mem/settings.json` (auto-created on first run).

Key settings to consider:

```json
{
  "worker": {
    "port": 37777
  },
  "context": {
    "maxTokens": 2000
  }
}
```

Adjust `maxTokens` based on your needs:
- **Lower (1000-1500)**: Faster sessions, less past context
- **Higher (2500-4000)**: More past context, slightly slower start
