# Fix Issue #{issue_number}

Fix a GitHub issue by selecting appropriate skills based on task classification.

## Task Classification

### Step 1: Identify Language

Check issue labels or affected files to determine language:

| Label / Files | Skill to Use |
|---------------|--------------|
| `lang:go`, `internal/`, `pkg/`, `cmd/` | `go-development` |
| `lang:typescript`, `apps/ripple-lib-server/` | `typescript-development` |
| `lang:solidity`, `apps/erc20-token/contracts/` | `solidity-development` |

### Step 2: Identify Chain (if applicable)

| Label | Chain Context |
|-------|---------------|
| `chain:btc` | Bitcoin-specific considerations |
| `chain:bch` | Bitcoin Cash-specific considerations |
| `chain:eth` | Ethereum-specific considerations |
| `chain:erc20` | ERC-20 token-specific considerations |
| `chain:xrp` | XRP/Ripple-specific considerations |
| `chain:all` | Cross-chain considerations |

## Process

1. **Fetch issue**: `gh issue view {issue_number}`
2. **Classify**: Identify language and chain from labels/files
3. **Load Skill**: Use appropriate `{lang}-development` Skill
4. **Implement**: Follow the Skill's workflow
5. **Verify**: Run Skill-specific verification commands
6. **Commit & PR**: Follow standard workflow

## Skills Reference

| Skill | Path |
|-------|------|
| Go | `.claude/skills/go-development/SKILL.md` |
| TypeScript | `.claude/skills/typescript-development/SKILL.md` |
| Solidity | `.claude/skills/solidity-development/SKILL.md` |

## Example

```
/fix-issue #123
```

If issue #123 has label `lang:go` and `chain:btc`:
→ Use `go-development` Skill with Bitcoin context
