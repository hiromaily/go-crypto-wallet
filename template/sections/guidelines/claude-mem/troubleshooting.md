### Troubleshooting

#### Worker Not Starting

```bash
# Check if port 37777 is in use
lsof -i :37777

# Manually start worker (if needed)
claude-mem worker start
```

#### No Context Being Injected

1. Verify the plugin is installed: check for claude-mem in `/plugin list`
2. Ensure the worker is running: visit `http://localhost:37777`
3. Check that previous sessions exist in the database

#### High Token Usage

If sessions start slowly due to too much injected context:
- Reduce `maxTokens` in settings
- Use `<private>` tags more liberally for non-essential discussions
