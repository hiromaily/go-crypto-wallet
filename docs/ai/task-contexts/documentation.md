<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/ai/task-contexts/documentation.tpl.md · Run `make docs` to regenerate.
-->

---
task_type: documentation
description: ドキュメント整備タスク用のコンテキスト
version: 1.0.0
---

# Task: Documentation (ドキュメント整備)

## When to Use

- README、AGENTS.md などのドキュメント更新
- コメント追加・改善
- API仕様書の作成・更新
- 設定ファイルの説明追加
- タスクコンテキストファイルの作成・更新

## Required Context Documents

### Must Read (必須)

ドキュメント整備タスクでは、通常のコーディングガイドラインは不要です。

| Document | Path | Purpose |
|----------|------|---------|
| (なし) | - | ドキュメント編集は最小限のコンテキストで実行可能 |

### Conditional Read (条件付き)

条件に応じて追加で読み込むドキュメント：

| Condition | Document | Path |
|-----------|----------|------|
| AGENTS.md 系の更新 | Architecture | `docs/guidelines/architecture.md` |
| タスクコンテキスト作成 | Task-Oriented Context | `docs/task-contexts/task-oriented-context.md` |
| チェーン固有ドキュメント | Chain Specific | `docs/chains/` |

## Task-Specific Rules

### 1. Go Linting は不要

```
❌ make go-lint（ドキュメント編集時は不要）
❌ make check-build（ドキュメント編集時は不要）
❌ make go-test（ドキュメント編集時は不要）
```

### 2. Markdown フォーマット推奨

```
✅ 見出しレベルの一貫性（# → ## → ### の順）
✅ コードブロックの言語指定
✅ テーブルの整形
✅ リンクの有効性確認
```

### 3. 内容の整合性

```
✅ 他のドキュメントとの整合性を保つ
✅ 古い情報の更新
✅ 用語の統一（日本語/英語の混在に注意）
```

## Pre-Task Checklist

- [ ] 更新対象のドキュメントを特定したか
- [ ] 関連するドキュメントとの整合性を確認したか
- [ ] 用語や表記の一貫性を確認したか

## Verification Commands

```bash
# ドキュメント編集時は Go 関連のコマンドは不要

# オプション: markdownlint がインストールされている場合
# markdownlint docs/**/*.md

# オプション: リンク確認
# markdown-link-check docs/**/*.md
```

**重要**: ドキュメントのみの変更では `make go-lint`, `make check-build`, `make go-test` は実行しないでください。

## Documentation Workflow

### Step 1: 対象の特定

1. 更新が必要なドキュメントを特定
2. 関連するドキュメントを確認（リンク先など）

### Step 2: 編集

1. 内容を更新
2. フォーマットを整える
3. リンクの有効性を確認

### Step 3: レビュー

1. 誤字脱字のチェック
2. 他のドキュメントとの整合性確認
3. 必要に応じて関連ドキュメントも更新

### Step 4: コミット

```bash
git add .
git commit -m "docs: {what was updated}

- {detail 1}
- {detail 2}"
```

## Common Documentation Types

### Type 1: README / AGENTS.md

プロジェクト概要やAIエージェント向けガイドライン

```markdown
# 構成例
- Overview
- Quick Start
- Directory Structure
- Guidelines
- References
```

### Type 2: タスクコンテキスト

タスク種別ごとのコンテキスト定義

```markdown
# 構成例
---
task_type: xxx
description: xxx
version: 1.0.0
---

# Task: XXX
## When to Use
## Required Context Documents
## Task-Specific Rules
## Verification Commands
```

### Type 3: API / 技術ドキュメント

技術仕様の説明

```markdown
# 構成例
- Overview
- Architecture
- API Reference
- Examples
- Troubleshooting
```

## File Type Detection

このタスクは以下のファイルタイプに適用されます：

| Extension | Description |
|-----------|-------------|
| `*.md` | Markdown ドキュメント |
| `*.mdc` | Cursor Rules (Markdown) |
| `*.txt` | テキストファイル |

## Examples

### Example 1: README 更新

```
User: "README.mdにインストール手順を追加して"

Agent Actions:
1. README.md を読み込み
2. 既存の構成を確認
3. インストール手順セクションを追加
4. フォーマットを整える
5. コミット（Go lint は不要）
```

### Example 2: タスクコンテキスト作成

```
User: "security-audit タスク用のコンテキストファイルを作成して"

Agent Actions:
1. docs/task-contexts/task-oriented-context.md を読み込み（フォーマット確認）
2. 既存のタスクコンテキストを参照
3. security-audit.md を作成
4. README.md のテーブルを更新
5. コミット（Go lint は不要）
```

### Example 3: コメント追加

```
User: "この関数にGoDocコメントを追加して"

Agent Actions:
1. 対象ファイルを確認
2. GoDoc形式でコメントを追加
3. make go-lint で確認（Go ファイル編集のため必要）
```

**注意**: Go ファイル内のコメント編集は、Go ファイルの編集扱いとなるため、通常の Go Verification が必要です。

## Related Documents

- [Task-Oriented Context Management](./task-oriented-context.md) - タスクコンテキストの概念
- [Task Contexts README](./README.md) - タスクコンテキスト一覧
- [Verification Matrix](./verification.md) - ファイルタイプ別検証コマンド
