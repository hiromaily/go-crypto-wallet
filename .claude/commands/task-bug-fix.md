# Bug Fix Task - {description} [Chain: {chain}]

## Task Type

Bug Fix (バグ修正)

## Required Context

このタスクを開始する前に、以下のドキュメントを読み込んでください：

### 必須ドキュメント

1. **Bug Fix Context**: `docs/ai-agents/task-contexts/bug-fix.md`
2. **Workflow Guidelines**: `agents/workflow.md`
3. **Core Principles**: `agents/core.md`

### チェーン固有ドキュメント（{chain} が指定された場合）

{chain} パラメータに応じて追加で読み込むドキュメント：

| Chain | Documents to Load |
|-------|-------------------|
| BTC | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/btc.md` |
| BCH | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/bch.md` |
| ETH | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/eth.md` |
| XRP | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/xrp.md` |

## Parameters

- `{description}`: バグの概要または Issue 番号（例: "Issue #123", "GetAddressInfo error"）
- `{chain}`: 対象の暗号通貨（BTC, BCH, ETH, XRP）。チェーン固有でない場合は省略可能

## Process

### Step 1: コンテキストのロード

1. 上記の必須ドキュメントを読み込む
2. {chain} が指定されている場合、チェーン固有ドキュメントも読み込む

### Step 2: 問題の分析

1. {description} が Issue 番号の場合: `gh issue view {issue_number}` で詳細を取得
2. 問題の再現手順を確認
3. 影響範囲を特定
4. 根本原因を調査

### Step 3: 実装

`docs/ai-agents/task-contexts/bug-fix.md` の Bug Fix Workflow に従って実装

### Step 4: 検証

```bash
make go-lint && make tidy && make check-build && make gotest
```

### Step 5: コミット & PR

```bash
git add <files>
git commit -m "fix: {description}

Closes #{issue_number}"

gh pr create --title "Fix: {description}"
```

## Examples

```
/task-bug-fix Issue #123 Chain: BTC
/task-bug-fix GetAddressInfo returns wrong result Chain: BCH
/task-bug-fix Database connection timeout
```

## Related Documents

- [Bug Fix Context](../../docs/ai-agents/task-contexts/bug-fix.md)
- [Chain-Specific Context](../../docs/ai-agents/task-contexts/chain-specific.md)
- [Workflow Guidelines](../../agents/workflow.md)

