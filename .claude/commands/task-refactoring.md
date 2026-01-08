# Refactoring Task - {description} [Chain: {chain}]

## Task Type

Refactoring (リファクタリング)

## Required Context

このタスクを開始する前に、以下のドキュメントを読み込んでください：

### 必須ドキュメント

1. **Refactoring Context**: `docs/ai-agents/task-contexts/refactoring.md`
2. **Architecture Guidelines**: `agents/architecture.md`
3. **Refactoring Plan**: `docs/issues/REFACTORING_PLAN.md`
4. **Testing Guidelines**: `agents/testing.md`

### チェーン固有ドキュメント（{chain} が指定された場合）

{chain} パラメータに応じて追加で読み込むドキュメント：

| Chain | Documents to Load |
|-------|-------------------|
| BTC | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/btc.md` |
| BCH | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/bch.md` |
| ETH | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/eth.md` |
| XRP | `docs/ai-agents/task-contexts/chain-specific.md`, `docs/ai-agents/task-contexts/chains/xrp.md` |

## Parameters

- `{description}`: リファクタリングの概要（例: "Use Case層への移動", "インターフェース抽出"）
- `{chain}`: 対象の暗号通貨（BTC, BCH, ETH, XRP）。チェーン固有でない場合は省略可能

## Important Rules

⚠️ **リファクタリングの原則**:

1. **機能変更なし**: 外部から見た動作は変更しない
2. **段階的な変更**: 小さな変更を積み重ねる
3. **テストの維持**: 既存のテストが全てパスすること

## Process

### Step 1: コンテキストのロード

1. 上記の必須ドキュメントを読み込む
2. {chain} が指定されている場合、チェーン固有ドキュメントも読み込む

### Step 2: 現状分析

1. リファクタリング対象のコードを特定
2. 依存関係を調査（どのパッケージがこのコードを使用しているか）
3. 既存のテストカバレッジを確認
4. REFACTORING_PLAN.md で関連タスクを確認

### Step 3: 計画

1. 変更を小さな単位に分割
2. 各ステップの検証方法を決定
3. リスクの高い変更を特定

### Step 4: 実装

`docs/ai-agents/task-contexts/refactoring.md` の Refactoring Workflow に従って実装

段階的に実装:
1. 準備（必要なテストを追加）
2. リファクタリング実行
3. クリーンアップ

### Step 5: 各ステップ後の検証

```bash
make go-lint && make gotest
```

### Step 6: 最終検証

```bash
make go-lint && make tidy && make check-build && make gotest
```

### Step 7: コミット & PR

```bash
git add .
git commit -m "refactor: {description}

No functional changes."

gh pr create --title "Refactor: {description}"
```

## Examples

```
/task-refactoring BTCトランザクション作成ロジックをUse Case層に移動 Chain: BTC
/task-refactoring Repository実装をインターフェースで抽象化
/task-refactoring BTC/BCH共通ロジックの抽出
```

## Related Documents

- [Refactoring Context](../../docs/ai-agents/task-contexts/refactoring.md)
- [Architecture Guidelines](../../agents/architecture.md)
- [Refactoring Plan](../../docs/issues/REFACTORING_PLAN.md)

