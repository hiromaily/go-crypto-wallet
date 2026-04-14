### Data Management

#### Storage Location

All data is stored locally at `~/.claude-mem/`:
- `sessions.db` - SQLite database
- `chroma/` - Vector database for semantic search

#### Viewing Past Sessions

Visit `http://localhost:37777` for the web viewer UI showing:
- Session timeline
- Observations per session
- AI-generated summaries

#### Clearing Data

```bash
# Remove all claude-mem data (irreversible)
rm -rf ~/.claude-mem/
```
