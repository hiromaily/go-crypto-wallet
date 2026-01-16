# Claude Code Commands

## Commands

| Command | Description |
|---------|-------------|
| `/fix-issue #123` | Work on a GitHub issue |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review #123` | Address PR review comments |
| `/recreate-pr #123` | Copy PR with review comments as fix checklist |

### BTC E2E Commands

| Command | Pattern | Description |
|---------|---------|-------------|
| `/e2e/fix-btc-e2e-p1` | 1 | P2PKH Single-sig |
| `/e2e/fix-btc-e2e-p2` | 2 | P2PKH 2-of-3 Multisig |
| `/e2e/fix-btc-e2e-p3` | 3 | P2SH-P2WPKH Single-sig |
| `/e2e/fix-btc-e2e-p4` | 4 | P2SH-P2WSH 2-of-3 Multisig |
| `/e2e/fix-btc-e2e-p5` | 5 | P2WPKH Native SegWit Single-sig |
| `/e2e/fix-btc-e2e-p6` | 6 | P2WSH 2-of-3 Multisig |
| `/e2e/fix-btc-e2e-p7` | 7 | P2WSH 3-of-3 Multisig |
| `/e2e/fix-btc-e2e-p8` | 8 | P2SH-P2WSH 3-of-3 Multisig |
| `/e2e/fix-btc-e2e-p9` | 9 | **P2TR Taproot Single-sig** (in progress) |

## Workflow

```
1. Create Issue                     2. Work on Issue
   github-issue-creation               git-workflow + task skill
         ↓                                    ↓
   Classify: Type + Lang/Scope        Use skill based on label
   Assign labels
```

## Related

- Label → Skill mapping: See [github-issue-creation](../skills/github-issue-creation/SKILL.md)
- Git workflow: See [git-workflow](../skills/git-workflow/SKILL.md)
