# Cursor Rules

> **Note**: このディレクトリのルールは `.claude/rules/` から自動生成されます。
> 直接編集せず、ソースを編集後 `make sync-cursor-rules` を実行してください。

## 概要

Cursor Rules は AI Agent に対するシステムレベルの指示を提供します。
ルールの内容はモデルコンテキストの先頭に含まれ、コード生成や編集に一貫したガイダンスを与えます。

## ルールタイプ

| タイプ | frontmatter | 説明 |
|--------|-------------|------|
| **Always Apply** | `alwaysApply: true` | すべてのチャットセッションに適用 |
| **Apply Intelligently** | `alwaysApply: false` + `description` | Agent が関連性を判断して適用 |
| **Apply to Specific Files** | `globs: ["pattern"]` | ファイルパターンにマッチした場合に適用 |
| **Apply Manually** | (description のみ) | `@rule-name` でメンション時に適用 |

## ファイル形式

### RULE.md 形式 (推奨)

\`\`\`
.cursor/rules/
  my-rule/
    RULE.md           # メインルールファイル
    scripts/          # ヘルパースクリプト (オプション)
\`\`\`

### .mdc 形式 (レガシー、引き続きサポート)

\`\`\`yaml
---
description: ルールの説明
globs:
  - "**/*.go"
  - "**/*.ts"
alwaysApply: false
---

# ルール本文
...
\`\`\`

## frontmatter プロパティ

| プロパティ | 型 | 説明 |
|-----------|------|------|
| `description` | string | ルールの説明（Apply Intelligently 時に必須） |
| `globs` | string[] | 適用対象のファイルパターン |
| `alwaysApply` | boolean | `true`: 常に適用、`false`: 条件付き適用 |

## このプロジェクトでの使い方

### SSOT (Single Source of Truth)

- **ソース**: `.claude/rules/*.md`
- **生成先**: `.cursor/rules/*.mdc`
- **変換コマンド**: `make sync-cursor-rules`

### 変換ルール

| Claude (`paths:`) | Cursor 出力 |
|-------------------|-------------|
| なし | `alwaysApply: true` |
| あり | `globs: ...` + `alwaysApply: false` |
| 最初の `# 見出し` | `description: ...` |
| `.md` 拡張子 | `.mdc` 拡張子 |

### 同期コマンド

\`\`\`bash
# dry-run で確認
make sync-cursor-rules-dry

# 実際に変換（既存ファイルを上書き）
make sync-cursor-rules
\`\`\`

## 参考

- [Cursor Rules Documentation](https://cursor.com/docs/context/rules)
- [AGENTS.md](AGENTS.md) - プロジェクト全体のガイドライン
