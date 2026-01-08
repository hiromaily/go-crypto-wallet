# Feature Add Task - {description} [Chain: {chain}]

## Task Type

Feature Add (新機能追加)

## Required Context

このタスクを開始する前に、以下のドキュメントを読み込んでください：

### 必須ドキュメント

1. **Feature Add Context**: `docs/ai-agents/task-contexts/feature-add.md`
2. **Architecture Guidelines**: `agents/architecture.md`
3. **Internal Guidelines**: `internal/AGENTS.md`
4. **Workflow Guidelines**: `agents/workflow.md`

### チェーン固有ドキュメント（{chain} が指定された場合）

{chain} パラメータに応じて追加で読み込むドキュメント：

| Chain | Documents to Load |
|-------|-------------------|
| BTC | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/btc.md`, `docs/crypto/btc/README.md` |
| BCH | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/bch.md`, `docs/crypto/bch/README.md` |
| ETH | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/eth.md`, `docs/crypto/eth/README.md` |
| XRP | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/xrp.md`, `docs/crypto/xrp/README.md` |

## Parameters

- `{description}`: 追加する機能の概要（例: "MuSig2署名集約", "ERC-20残高確認"）
- `{chain}`: 対象の暗号通貨（BTC, BCH, ETH, XRP）。チェーン固有でない場合は省略可能

## Process

### Step 1: コンテキストのロード

1. 上記の必須ドキュメントを読み込む
2. {chain} が指定されている場合、チェーン固有ドキュメントも読み込む

### Step 2: 要件分析

1. 機能要件を明確にする
2. 影響するレイヤーを特定（Domain/Application/Infrastructure/Interface Adapters）
3. 既存の類似実装を確認

### Step 3: 設計

1. レイヤー別の実装計画を立てる
2. 必要なインターフェース（Ports）を特定
3. ファイル配置を決定

### Step 4: 実装

`docs/ai-agents/task-contexts/feature-add.md` の Feature Add Workflow に従って実装

実装順序:
1. Domain層（エンティティ、値オブジェクト）
2. Application層（Ports、DTO、Use Case）
3. Infrastructure層（Repository、API Client）
4. Interface Adapters層（CLI、HTTP Handler）
5. DI（internal/di/container.go）

### Step 5: テスト

```bash
make gotest
```

### Step 6: 検証

```bash
make go-lint && make tidy && make check-build && make gotest
```

### Step 7: コミット & PR

```bash
git add .
git commit -m "feat: add {description}

- Implementation details..."

gh pr create --title "Feature: {description}"
```

## Examples

```
/task-feature-add MuSig2署名集約機能 Chain: BTC
/task-feature-add ERC-20トークン残高確認 Chain: ETH
/task-feature-add HD Wallet生成の改善
```

## Related Documents

- [Feature Add Context](../../docs/ai-agents/task-contexts/feature-add.md)
- [Architecture Guidelines](../../agents/architecture.md)
- [Chain-Specific Context](../../docs/ai-agents/task-contexts/chain-specific.md)

