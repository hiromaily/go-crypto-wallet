---
paths:
  - ".claude/**"
  - ".cursor/**"
  - ".codex/**"
  - ".github/copilot-instructions.md"
---

# Agent Files

このルールは `.claude/`, `.cursor/`, `.codex/`, `.github/copilot-instructions.md` を編集する際に適用されます。

## Single Source of Truth (SSOT) 設計

AI Agent 設定ファイルは **Claude Code の `.claude/` ディレクトリを唯一のソース** として管理し、他のエージェント用設定は自動生成またはシンボリックリンクで対応します。

```
.claude/                    ← SSOT (Source of Truth)
├── commands/               ← Slash commands
├── rules/                  ← Rules (*.md)
└── skills/                 ← Skills (SKILL.md)

.cursor/                    ← 自動生成 / シンボリックリンク
├── commands/README.md      ← 参照のみ (.claude/commands を使用)
├── rules/*.mdc             ← 自動生成 (scripts/ai-agent/sync-rule-claude-to-cursor.sh)
└── skills/                 ← シンボリックリンク → ../.claude/skills

.codex/                     ← TODO (将来対応)
.github/copilot-instructions.md ← TODO (将来対応)
```

## 対応 AI Agent

| Agent | Version | Status | 設定場所 |
|-------|---------|--------|----------|
| Claude Code | v2 | ✅ Active | `.claude/` (SSOT) |
| Cursor | v2 | ✅ Active | `.cursor/` (自動生成) |
| Codex | v0.80 | 📋 TODO | `.codex/` |
| GitHub Copilot | 2026 | 📋 TODO | `.github/copilot-instructions.md` |

## ファイル参照形式

### Claude Code

```markdown
- Follow @AGENTS.md for guidelines
- Refer to @docs/standards/coding-conventions.md
```

### Cursor

同じ `@path` 形式をサポート。変換不要。

## ディレクトリ別ルール

### `.claude/commands/`

Slash commands の定義場所。Cursor からも自動的に読み込まれる。

**編集時の注意:**

- 新規コマンド追加は `.claude/commands/` に作成
- `.cursor/commands/` には README.md のみ配置

### `.claude/rules/`

Claude Code 用ルールファイル (`.md`)。

**フォーマット:**

```markdown
---
paths:
  - "**/*.go"
  - "**/*.ts"
---

# Rule Title

ルール内容...
```

- `paths:` がない場合 → すべての指示に適用 (グローバルルール)
- `paths:` がある場合 → 指定パターンにマッチするファイル編集時に適用

### `.claude/skills/`

Skills (MCP) の定義場所。各スキルはサブディレクトリに `SKILL.md` を配置。

```
.claude/skills/
├── go-development/SKILL.md
├── git-workflow/SKILL.md
└── db-migration/SKILL.md
```

### `.cursor/rules/`

**自動生成ファイル - 直接編集禁止**

Claude rules から自動生成される。変換ルール:

| Claude | Cursor |
|--------|--------|
| `paths:` なし | `alwaysApply: true` |
| `paths:` あり | `globs:` + `alwaysApply: false` |
| 最初の `# 見出し` | `description:` |
| `.md` | `.mdc` |

**同期コマンド:**

```bash
./scripts/ai-agent/sync-rule-claude-to-cursor.sh --force --verbose
```

## Model Context Protocol (MCP) / Skills

2024年末から提唱された MCP が2026年現在普及し、AI Agent が「Skills」としてリポジトリ内のスクリプトやローカルサーバーを道具として使用可能。

### Claude Code / Cursor

`.claude/skills/` または `.cursor/skills/` に `SKILL.md` を配置。

### Codex (将来対応)

CLI ベースの Codex agent では、Shell/Python スクリプトを Skill としてバインドし、`@database_tool` のように呼び出す形式を想定。

## 編集時のチェックリスト

`.claude/` または `.cursor/` を編集する場合:

- [ ] 変更は `.claude/` (SSOT) に対して行う
- [ ] `.cursor/rules/` を直接編集しない
- [ ] 新規 rules 追加後は同期スクリプトを実行
- [ ] README.md の更新が必要か確認

## 関連ドキュメント

- @.cursor/rules/README.md - Cursor rules の仕様
- @.cursor/commands/README.md - Commands の説明
- @AGENTS.md - プロジェクト全体のガイドライン
