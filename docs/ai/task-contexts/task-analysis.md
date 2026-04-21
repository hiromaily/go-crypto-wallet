<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/ai/task-contexts/task-analysis.tpl.md · Run `make docs` to regenerate.
-->

# Task Pattern Analysis

このドキュメントは、GitHub Issues と Commit 履歴から分析したタスクパターンをまとめています。
定期的に更新することで、AI Agent のコンテキストローディングを最適化できます。

## Issue Label Distribution

最新の分析結果（直近100件のIssue）:

| Label | Count | Priority |
|-------|-------|----------|
| technical-debt | 17 | P1 |
| refactoring | 15 | P0 (既存) |
| enhancement | 15 | P0 (既存) |
| bug | 3 | P0 (既存) |

### 分析結果

1. **technical-debt** が最も多い → `refactoring` コンテキストでカバー可能
2. **refactoring** と **enhancement** が同数 → 両タスクタイプが重要
3. **bug** は比較的少ない → バグ修正より新機能・改善が多い

## Commit Type Distribution

最新の分析結果（直近200件のコミット）:

| Commit Type | Count | Mapped Task Type |
|-------------|-------|------------------|
| feat: | 24 | `feature-add` |
| fix: | 9 | `bug-fix` |
| refactor: | 3 | `refactoring` |
| docs: | 3 | `documentation` |

### 分析結果

1. **feat:** が最も多い → 新機能追加が主要タスク
2. **fix:** が2番目 → バグ修正も継続的に発生
3. **refactor:** は少なめ → 大きなリファクタリングは Issue ベースで管理
4. **docs:** は限定的 → ドキュメントタスクは比較的少ない

## Recent Activity Patterns

直近のコミットから見えるパターン：

### 1. Feature Development

- YAML設定サポート追加 (#284, #285)
- E2Eワークフロースクリプト追加 (#276, #277)
- DevContainer導入 (#264, #265)

### 2. Refactoring

- Clean Architecture への移行（Phase 1-3）
- Persistence Port インターフェース導入 (#271-275)

### 3. Bug Fixes

- Bitcoin Core バージョン対応 (#266-269)
- BCH コンテナ起動問題 (#280, #281)
- Multisig descriptor インポート (#259, #260)

### 4. Chain-Specific Work

- **BTC**: Descriptor, E2E workflow, Bitcoin Core 対応
- **BCH**: E2E workflow, コンテナ修正
- **ETH**: Anvil ローカルノード追加

## Recommended Context Loading Priority

分析に基づく、タスクタイプの読み込み優先度：

| Priority | Task Type | Reason |
|----------|-----------|--------|
| P0 | `feature-add` | 最頻出コミットタイプ |
| P0 | `bug-fix` | 2番目に多いコミットタイプ |
| P1 | `refactoring` | Issue で多い、Phase分けで進行中 |
| P2 | `db-change` | 発生時は重要だが頻度低 |
| P2 | `documentation` | 頻度低だが独立したタスク |

## Chain-Specific Activity

最近のチェーン別アクティビティ：

| Chain | Recent Issues/PRs | Focus Area |
|-------|-------------------|------------|
| BTC | #266-269, #276-279 | Bitcoin Core 対応, E2E テスト |
| BCH | #280-281 | コンテナ環境, E2E テスト |
| ETH | #262-263 | ローカル開発環境 |
| XRP | - | 最近のアクティビティ少 |

## File Change Patterns

頻繁に変更されるファイルパターン：

1. **internal/infrastructure/** - インフラ層の実装
2. **internal/application/usecase/** - ユースケース実装
3. **config/wallet/** - 設定ファイル
4. **scripts/operation/** - 運用スクリプト
5. **docs/** - ドキュメント

## Recommendations for Context Optimization

### 1. タスクタイプ別コンテキスト

現在のコンテキストファイルは適切。追加検討：

- [ ] `test-add.md` - テスト追加タスク用（E2E テストの増加）
- [ ] `e2e-workflow.md` - E2E ワークフロー関連

### 2. チェーン別コンテキスト

- BTC: 最もアクティブ、現在のコンテキストで十分
- BCH: BTCとの関係性が重要（オーバーライドパターン）
- ETH: ローカル環境構築の情報追加を検討
- XRP: 現状維持

### 3. Verification 最適化

Issue/Commit パターンから：

- Go ファイル変更が多い → Go lint + build は必須
- ドキュメント変更は独立 → Go 検証スキップで効率化
- 設定ファイル変更は頻繁 → 構文チェックのみ

## Analysis Commands

このドキュメントの更新に使用するコマンド：

```bash
# Issue ラベル分析
gh issue list --state all --json labels --limit 100 | jq -r '.[].labels[].name' | sort | uniq -c | sort -rn

# Commit タイプ分析
git log --oneline -200 | grep -oE '^[a-f0-9]+ (feat|fix|refactor|docs|chore|test):' | sed 's/^[a-f0-9]* //' | sort | uniq -c | sort -rn

# 最近のコミット確認
git log --oneline -30

# 頻繁に変更されるファイル
git log --pretty=format: --name-only -100 | grep -E '\.(go|md|yaml|toml|sql|hcl)$' | sort | uniq -c | sort -rn | head -15
```

## Update Schedule

このドキュメントは以下のタイミングで更新を推奨：

- 月次: Issue/Commit パターンの再分析
- 四半期: コンテキスト優先度の見直し
- 必要時: 新しいタスクパターンの発見時

## Related Documents

- [Task Contexts README](./README.md)
- [Task-Oriented Context Management](./task-oriented-context.md)
- `AGENTS.md`
