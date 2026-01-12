---
task_type: bug-fix
description: バグ修正・Issue対応タスク用のコンテキスト
version: 1.0.0
---

# Task: Bug Fix (バグ修正)

## When to Use

- GitHub Issueに報告されたバグを修正する場合
- 既存機能のエラーや不具合を修正する場合
- テスト失敗の原因を特定・修正する場合
- パフォーマンス問題を解決する場合

## Required Context Documents

### Must Read (必須)

以下のドキュメントは必ず読み込んでください：

| Document | Path | Purpose |
|----------|------|---------|
| Core Principles | `docs/guidelines/core.md` | エラーハンドリング、セキュリティ |
| Coding Standards | `docs/guidelines/coding-standards.md` | コーディング規約 |
| Workflow Guidelines | `docs/guidelines/workflow.md` | 検証ステップ、PRフロー |

### Conditional Read (条件付き)

条件に応じて追加で読み込むドキュメント：

| Condition | Document | Path |
|-----------|----------|------|
| アーキテクチャ関連のバグ | Architecture | `docs/guidelines/architecture.md` |
| データベース関連のバグ | Database | `docs/guidelines/database.md` |
| 特定チェーン固有の問題 | Multi-Chain | `docs/guidelines/multi-chain.md` |
| セキュリティ関連のバグ | Core (Security) | `docs/guidelines/core.md` |
| テスト関連の問題 | Testing | `docs/guidelines/testing.md` |

## Task-Specific Rules

### 1. 根本原因の特定

```
❌ 症状だけを修正する
✅ 根本原因を特定してから修正する
```

- 問題を再現する手順を明確にする
- スタックトレースやログを分析する
- 関連するコードパスを追跡する

### 2. 修正範囲の最小化

```
❌ "ついでに" 他の箇所も修正する
✅ バグに直接関連する箇所のみ修正する
```

- 修正は最小限に留める
- 無関係なリファクタリングは別PRで行う

### 3. エラーハンドリング

```go
// ❌ Bad: エラーを握りつぶす
result, _ := someFunc()

// ✅ Good: エラーを適切に処理
result, err := someFunc()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

### 4. セキュリティ考慮

- **NEVER**: 秘密鍵やセンシティブ情報をログに出力しない
- **ALWAYS**: セキュリティ関連の修正は慎重にレビュー
- **CHECK**: オフラインウォレット操作への影響を確認

## Pre-Task Checklist

- [ ] 問題を再現できるか確認したか
- [ ] 影響範囲を特定したか
- [ ] 既存のテストが何をカバーしているか確認したか
- [ ] セキュリティへの影響を評価したか
- [ ] 関連するIssueやPRを確認したか

## Verification Commands

```bash
# 必須の検証コマンド（修正後に必ず実行）
make go-lint      # リンターチェック
make tidy         # 依存関係の整理
make check-build  # ビルド確認
make gotest       # テスト実行

# 追加の検証（該当する場合）
make gotest-integration  # 統合テスト
make check-vuln          # セキュリティ脆弱性チェック（セキュリティ関連の修正時）
```

## Bug Fix Workflow

### Step 1: 問題の分析

```bash
# Issueの詳細を取得
gh issue view {issue_number}

# 関連するコードを特定
# - エラーメッセージで検索
# - スタックトレースを追跡
```

### Step 2: 再現と調査

1. 問題を再現する
2. 関連するテストがあれば実行して失敗を確認
3. デバッグして根本原因を特定

### Step 3: 修正の実装

1. 最小限の修正を実装
2. 修正が問題を解決することを確認
3. 回帰を防ぐテストを追加

### Step 4: 検証

```bash
# 全ての検証コマンドを実行
make go-lint && make tidy && make check-build && make gotest
```

### Step 5: コミットとPR

```bash
# コミット
git add <modified-files>
git commit -m "fix: resolve issue #{issue_number} - {brief description}

- {detail 1}
- {detail 2}

Closes #{issue_number}"

# PRを作成
gh pr create --title "Fix: {issue title} (Closes #{issue_number})"
```

## Testing Requirements

### 必須

- [ ] 修正が問題を解決することを確認するテスト
- [ ] 既存のテストが引き続きパスすること

### 推奨

- [ ] 回帰を防ぐための新しいテストケース
- [ ] エッジケースのテスト

## Examples

### Example 1: エラーハンドリングのバグ修正

```
User: "Issue #123: トランザクション作成時にpanicが発生する"

Agent Actions:
1. gh issue view 123 でIssue内容を確認
2. docs/guidelines/core.md を読み込み（エラーハンドリング規約）
3. パニックが発生するコードパスを特定
4. 適切なエラーハンドリングを追加
5. テストを追加して回帰を防止
6. 検証コマンドを実行
7. PRを作成
```

### Example 2: データベース関連のバグ修正

```
User: "Issue #456: アドレス検索で重複結果が返される"

Agent Actions:
1. gh issue view 456 でIssue内容を確認
2. docs/guidelines/database.md を追加で読み込み
3. クエリロジックを調査
4. DISTINCT追加やクエリ修正
5. 単体テストとSQLCクエリテストを追加
6. 検証コマンドを実行
7. PRを作成
```

## Related Commands

- `/fix-issue #{issue_number}` - Issue修正の自動化コマンド
- `/fix-linter` - リンターエラー修正

## Related Documents

- [Fix Issue Command](../../../.claude/commands/fix-issue.md) - Issue修正コマンド詳細
- [Workflow Guidelines](../guidelines/workflow.md) - 検証ステップ詳細
- [Core Principles](../guidelines/core.md) - エラーハンドリング詳細
