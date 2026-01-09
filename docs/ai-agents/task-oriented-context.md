# Task-Oriented Context Management

## Overview

タスク指向のコンテキスト管理とは、AI Agentに依頼するタスクの種類に応じて、必要なドキュメントとコンテキストを最適化して提供する手法です。これにより、トークン効率の向上、回答精度の改善、一貫したワークフローの実現が可能になります。

## Problem Statement

### 現状の課題

1. **コンテキストの肥大化**: 全てのガイドラインを常にロードすると、トークン消費が増大
2. **関連性の低下**: 不要な情報がノイズとなり、重要な指示が埋もれる
3. **一貫性の欠如**: タスクによって必要なドキュメントが異なるが、明示的なマッピングがない
4. **オンボーディングコスト**: 新しいAgentツールへの移行時に再構成が必要

### 現在のドキュメント構造

```
go-crypto-wallet/
├── AGENTS.md                    # エントリーポイント（常にロード）
├── docs/ai-agents/              # AI Agent関連ドキュメント
│   ├── guidelines/              # トピック別ガイドライン
│   │   ├── architecture.md
│   │   ├── coding-standards.md
│   │   ├── database.md
│   │   ├── testing.md
│   │   └── ...
│   ├── task-contexts/           # タスク別コンテキスト
│   └── agent-skills.md          # Skills使用ガイド
├── .claude/
│   ├── commands/                # カスタムコマンド
│   └── skills/                  # Agent Skills
└── internal/AGENTS.md           # internal固有ガイドライン
```

## Proposed Solution

### Task Type Classification

タスクを以下のカテゴリに分類し、各タスクタイプに必要なコンテキストを定義します：

| Task Type | Description | Primary Use Cases |
|-----------|-------------|-------------------|
| `bug-fix` | バグ修正 | Issue修正、エラー対応 |
| `feature-add` | 新機能追加 | Use Case追加、新エンドポイント |
| `refactoring` | リファクタリング | コード整理、アーキテクチャ改善 |
| `db-change` | データベース変更 | スキーマ変更、マイグレーション |
| `code-review` | コードレビュー | PR確認、品質チェック |
| `security-audit` | セキュリティ監査 | 脆弱性チェック、鍵管理確認 |
| `test-add` | テスト追加 | 単体テスト、統合テスト作成 |
| `documentation` | ドキュメント更新 | README、API仕様更新 |
| `code-generation` | コード生成 | SQLC、protobuf、mock生成 |
| `multi-chain` | マルチチェーン対応 | BTC/ETH/XRP固有の実装 |

### Directory Structure

```
docs/ai-agents/
├── README.md                     # 概要とクイックスタート
├── agent-skills.md               # 既存: Skills使用ガイド
├── task-oriented-context.md      # 本ドキュメント
└── task-contexts/                # タスクコンテキスト定義
    ├── README.md                 # タスクタイプ一覧
    ├── bug-fix.md
    ├── feature-add.md
    ├── refactoring.md
    ├── db-change.md
    ├── code-review.md
    ├── security-audit.md
    ├── test-add.md
    ├── documentation.md
    ├── code-generation.md
    └── multi-chain.md
```

## Task Context File Format

各タスクコンテキストファイルは以下の構造に従います：

```markdown
---
task_type: <task-type-name>
description: タスクの説明
version: 1.0.0
---

# Task: <Task Name>

## When to Use
- このタスクタイプを使用する条件

## Required Context Documents

### Must Read (必須)
| Document | Path | Purpose |
|----------|------|---------|

### Conditional Read (条件付き)
| Condition | Document | Path |
|-----------|----------|------|

## Task-Specific Rules
タスク固有のルール

## Pre-Task Checklist
- [ ] チェック項目

## Verification Commands
検証コマンド

## Examples
使用例
```

## Integration Strategy

### Phase 1: ドキュメント作成

1. `docs/ai-agents/task-contexts/` ディレクトリ作成
2. 主要タスクタイプのコンテキストファイル作成（優先度順）:
   - `bug-fix.md` - 最も頻繁
   - `feature-add.md` - 開発の中心
   - `refactoring.md` - 継続的改善
   - `db-change.md` - 特殊手順が必要

### Phase 2: 既存システムとの統合

1. **AGENTS.mdの更新**
   - タスクコンテキストへのリンク追加
   - 「When to Use Each Document」テーブルの拡張

2. **カスタムコマンドの更新**
   - `fix-issue.md` にタスクコンテキスト参照を追加
   - 新規コマンドでタスクコンテキストを自動ロード

3. **Agent Skillsとの連携**
   - Skillの `description` にタスクタイプを明記
   - Skill内でタスクコンテキストを参照

### Phase 3: 運用と改善

1. **使用状況の収集**
   - どのタスクタイプが頻繁に使われるか
   - 追加のコンテキストが必要なケース

2. **フィードバックループ**
   - タスクコンテキストの改善
   - 新しいタスクタイプの追加

## Usage Patterns

### Pattern 1: 明示的なタスクタイプ指定

```
User: "Task Type: bug-fix. Issue #123を修正して"

Agent:
1. docs/ai-agents/task-contexts/bug-fix.md を読み込み
2. 必須ドキュメントをロード
3. タスク実行
```

### Pattern 2: Agent自動判定

```
User: "Use CaseにMuSig2署名機能を追加して"

Agent:
1. 内容から "feature-add" タイプと判定
2. docs/ai-agents/task-contexts/feature-add.md を読み込み
3. タスク実行
```

### Pattern 3: カスタムコマンド経由

```
User: "/fix-issue #123"

Command Behavior:
1. Issue内容を分析
2. 適切なタスクタイプを判定（bug-fix, feature-add等）
3. 対応するタスクコンテキストをロード
4. タスク実行
```

## Benefits

| Aspect | Benefit |
|--------|---------|
| **トークン効率** | 必要なドキュメントのみロード（最大50%削減可能） |
| **精度向上** | 焦点を絞ったコンテキストで関連性の高い回答 |
| **一貫性** | 同じタスクタイプは同じプロセスで処理 |
| **保守性** | タスクタイプ毎にドキュメント更新が容易 |
| **拡張性** | 新しいタスクタイプの追加が容易 |
| **ツール非依存** | Claude/Cursor/Windsurf等で同じコンテキスト使用可能 |

## Implementation Priority

| Priority | Task Type | Reason |
|----------|-----------|--------|
| P0 | `bug-fix` | 最頻出、既存コマンドあり |
| P0 | `feature-add` | 開発の中心タスク |
| P1 | `refactoring` | リファクタリング進行中 |
| P1 | `db-change` | 複雑な手順が必要 |
| P2 | `code-generation` | 自動生成ファイルの扱い |
| P2 | `security-audit` | 重要だが頻度低 |
| P3 | その他 | 必要に応じて追加 |

## Migration Path

### From Current State

```
現在: AGENTS.md → agents/*.md (全てフラット)
     ↓
移行: AGENTS.md → Task Type判定 → task-contexts/*.md → agents/*.md (必要なもののみ)
```

### Backward Compatibility

- 既存のガイドラインは `docs/ai-agents/guidelines/` に統合済み
- タスクコンテキストは追加レイヤーとして機能
- 明示的なタスクタイプ指定がなくても従来通り動作

## Related Documents

- [AGENTS.md](../../AGENTS.md) - プロジェクトガイドライン
- [Agent Skills](./agent-skills.md) - Agent Skills使用ガイド
- [Task Contexts](./task-contexts/README.md) - タスクコンテキスト一覧
- [Task Analysis](./task-analysis.md) - Issue/Commit パターン分析
- [Verification Matrix](./task-contexts/verification.md) - ファイルタイプ別検証コマンド
- [Workflow Guidelines](./guidelines/workflow.md) - 共通ワークフロー
