### Usage

#### Automatic Operation

Claude-mem works automatically once installed. No manual intervention needed for:
- Capturing observations during sessions
- Generating session summaries
- Injecting relevant context at session start

#### Manual Search (mem-search Skill)

Use natural language queries to search past sessions:

```
/mem-search "How did we implement the PostgreSQL repository layer?"
/mem-search "What approach was used for BCH address validation?"
/mem-search "What errors occurred during the XRP integration?"
```

#### 3-Layer Progressive Disclosure

Claude-mem uses a token-efficient 3-layer approach:

1. **`search`** - Returns compact index with observation IDs (~50-100 tokens per result)
2. **`timeline`** - Shows chronological context around specific results
3. **`get_observations`** - Fetches full details for filtered IDs (~500-1,000 tokens per result)

This yields ~10x token savings compared to fetching all details upfront.

#### Privacy Controls

Wrap sensitive content in `<private>` tags to exclude it from memory:

```
<private>
This content will NOT be stored in claude-mem.
Use for private keys, credentials, or sensitive discussion.
</private>
```

**Important for this project**: Since we handle cryptocurrency private keys, always use `<private>` tags when discussing or working with:
- Private keys and seeds
- Wallet passwords
- API credentials
- Any content that would violate our [security guidelines](../../../../docs/guidelines/security.md)
