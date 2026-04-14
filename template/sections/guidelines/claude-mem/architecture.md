### Architecture

```
Claude Code Session
    │
    ├─ SessionStart hook ──────► Worker (port 37777)
    ├─ UserPromptSubmit hook ──►     │
    ├─ PostToolUse hook ───────►     ├─ SQLite DB (sessions, observations)
    ├─ Stop hook ──────────────►     ├─ Chroma Vector DB (semantic search)
    └─ SessionEnd hook ────────►     └─ AI Summarizer
```

#### 5 Lifecycle Hooks

| Hook | When | What It Captures |
|------|------|-----------------|
| `SessionStart` | Session begins | Injects relevant past context |
| `UserPromptSubmit` | User sends prompt | Records user intent |
| `PostToolUse` | After each tool call | Records tool usage and results |
| `Stop` | Claude stops responding | Captures completion state |
| `SessionEnd` | Session ends | Generates session summary |
