# Task Contexts

このディレクトリには、タスクタイプ別のコンテキスト定義ファイルが含まれています。

## Quick Reference

| Task Type | File | When to Use |
|-----------|------|-------------|
| Bug Fix | [bug-fix.md](./bug-fix.md) | Issue修正、バグ対応、エラー修正 |
| Feature Add | [feature-add.md](./feature-add.md) | 新機能追加、Use Case追加、API追加 |
| Refactoring | [refactoring.md](./refactoring.md) | コード整理、アーキテクチャ改善 |
| DB Change | [db-change.md](./db-change.md) | スキーマ変更、マイグレーション |
| Documentation | [documentation.md](./documentation.md) | ドキュメント更新、README整備、コメント追加 |
| Test | [test.md](./test.md) | テスト追加・修正、カバレッジ向上 |
| Chain-Specific | [chain-specific.md](./chain-specific.md) | 暗号通貨固有のタスク（BTC/BCH/ETH/XRP） |
| Verification | [verification.md](./verification.md) | ファイルタイプ別の検証コマンドマトリックス |

## Chain-Specific References

暗号通貨固有のタスクを処理する場合は、[chain-specific.md](./chain-specific.md) を参照し、チェーンを特定した後に `docs/chains/` の該当ドキュメントを読み込んでください：

| Chain | Documentation | Key Features |
|-------|---------------|--------------|
| BTC | [docs/chains/btc/README.md](../chains/btc/README.md) | UTXO, Descriptor, Taproot, MuSig2 |
| BCH | [docs/chains/bch/README.md](../chains/bch/README.md) | UTXO, CashAddr |
| ETH | [docs/chains/eth/README.md](../chains/eth/README.md) | Account, Gas, ERC-20 |
| XRP | [docs/chains/xrp/README.md](../chains/xrp/README.md) | Account, gRPC, Destination Tag |

## Usage

### For AI Agents

タスクを受け取った際、以下の手順でコンテキストをロードしてください：

1. **タスクタイプの判定**: ユーザーの依頼内容からタスクタイプを特定
2. **コンテキストファイルの読み込み**: 該当する `task-contexts/*.md` を読み込む
3. **必須ドキュメントのロード**: コンテキストファイルで指定された必須ドキュメントを読み込む
4. **タスク実行**: コンテキストファイルのルールに従ってタスクを実行

### For Users

タスク依頼時にタスクタイプを明示することで、Agentが適切なコンテキストをロードしやすくなります：

```
# 明示的な指定
"Task Type: bug-fix. Issue #123 の問題を修正して"

# 暗黙的（Agentが判定）
"このエラーを修正して"  → bug-fix と判定
"新しい機能を追加して"  → feature-add と判定
```

## Task Type Details

### bug-fix

**用途**: バグ修正、Issue対応、エラー修正

**主な読み込みドキュメント**:

- `docs/guidelines/core.md` - エラーハンドリング
- `docs/guidelines/coding-standards.md` - コーディング規約
- `docs/guidelines/workflow.md` - 検証ステップ

### feature-add

**用途**: 新機能追加、Use Case実装、API追加

**主な読み込みドキュメント**:

- `docs/guidelines/architecture.md` - Clean Architecture
- `docs/guidelines/coding-standards.md` - コーディング規約
- `internal/AGENTS.md` - レイヤー構造

### refactoring

**用途**: コードリファクタリング、アーキテクチャ改善

**主な読み込みドキュメント**:

- `docs/guidelines/architecture.md` - Clean Architecture
- `docs/issues/REFACTORING_PLAN.md` - リファクタリング計画
- `docs/guidelines/testing.md` - 既存テストの維持

### db-change

**用途**: データベーススキーマ変更、マイグレーション

**主な読み込みドキュメント**:

- `docs/database/db-management.md` - Atlas/SQLC手順
- `docs/guidelines/code-generation.md` - コード生成
- `tools/atlas/` - スキーマファイル

## Adding New Task Types

新しいタスクタイプを追加する場合：

1. このディレクトリに `{task-type}.md` ファイルを作成
2. [Task Context File Format](./task-oriented-context.md#task-context-file-format) に従って記述
3. このREADMEのQuick Referenceテーブルを更新
4. 必要に応じて `AGENTS.md` のナビゲーションを更新

## Related Documents

- [Task-Oriented Context Management](./task-oriented-context.md) - 概念と戦略
- [Task Analysis](./task-analysis.md) - Issue/Commit パターン分析
- `AGENTS.md` - プロジェクトガイドライン
- [Workflow Guidelines](../guidelines/workflow.md) - 共通ワークフロー
