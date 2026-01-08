# Custom Slash Commands

This directory contains definition files for custom slash commands that can be used in Claude Desktop.

## Architecture

These commands follow a modular structure where:

- **Common workflow steps** are defined in [Workflow Guidelines](../../agents/workflow.md)
- **Command-specific logic** is defined in each command file
- This reduces duplication and makes maintenance easier

## Migration to Agent Skills

**Note**: The `create-github-issue` command has been migrated to Agent Skills format for better team collaboration.

### New Location (Agent Skills)

The GitHub issue creation functionality is now available as an Agent Skill:

- **Location**: `.claude/skills/github-issue-creation/SKILL.md`
- **Usage**: Use the Skill tool with `skill: "github-issue-creation"`
- **Benefits**: Automatic discovery, team-wide sharing, progressive disclosure

See [Agent Skills documentation](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview) for more information.

### Deprecated Commands

The following commands are being migrated to Agent Skills format:

#### `create-github-issue` (Deprecated - Use Agent Skill instead)

Creates a GitHub issue.

- **New location**: `.claude/skills/github-issue-creation/SKILL.md`
- Creates a well-structured issue with understanding of project context
- Suggests appropriate labels
- Includes detailed template with security and architecture considerations
- See [create-github-issue.md](create-github-issue.md) for legacy reference

## Available Commands

### Task-Oriented Commands (タスク指向コマンド)

タスクタイプに応じて適切なコンテキストを自動的にロードするコマンドです。

#### `task-bug-fix`

バグ修正タスク用のコンテキストをロードして作業を開始します。

- 必要なコンテキストドキュメントを自動ロード
- チェーン指定（BTC, BCH, ETH, XRP）に対応
- Bug Fix ワークフローに従って実装
- See [task-bug-fix.md](task-bug-fix.md) for details

**Usage**: `/task-bug-fix {description} Chain: {chain}`

#### `task-feature-add`

新機能追加タスク用のコンテキストをロードして作業を開始します。

- Clean Architecture に準拠した実装ガイド
- チェーン固有の実装パターンを提供
- レイヤー別の実装順序を案内
- See [task-feature-add.md](task-feature-add.md) for details

**Usage**: `/task-feature-add {description} Chain: {chain}`

#### `task-refactoring`

リファクタリングタスク用のコンテキストをロードして作業を開始します。

- 機能変更なしの原則を徹底
- 段階的な変更アプローチ
- 既存テストの維持を確認
- See [task-refactoring.md](task-refactoring.md) for details

**Usage**: `/task-refactoring {description} Chain: {chain}`

#### `task-db-change`

データベース変更タスク用のコンテキストをロードして作業を開始します。

- Atlas/SQLC の手順を自動案内
- スキーマ変更からコード生成まで一貫したワークフロー
- Docker 環境での検証手順を含む
- See [task-db-change.md](task-db-change.md) for details

**Usage**: `/task-db-change {description}`

### Issue/PR Commands (従来のコマンド)

### `fix-issue`

Resolves a specified GitHub issue.

- Fetches and analyzes issue content
- Creates a feature branch
- Implements code fixes
- Tests and verifies changes
- Commits and creates a pull request
- See [fix-issue.md](fix-issue.md) for details

### `fix-linter`

Fixes linter errors.

- Analyzes errors detected by `make lint-fix`
- Prioritizes and fixes errors
- Uses a step-by-step approach to fixes
- See [fix-linter.md](fix-linter.md) for details

### `fix-pr-review`

Addresses pull request review comments.

- Fetches PR information and review comments
- Categorizes and prioritizes comments
- Implements fixes and tests
- Commits and pushes changes
- See [fix-pr-review.md](fix-pr-review.md) for details

### `convert-custom-slash-for-codex`

Converts Claude custom slash commands to Codex-optimized prompt format.

- Reads Claude commands from `.claude/commands/`
- Extracts metadata (description, argument hints)
- Transforms parameter placeholders to Codex syntax
- Generates YAML frontmatter
- Outputs to `.codex/prompts/`
- See [convert-custom-slash-for-codex.md](convert-custom-slash-for-codex.md) for details

## Usage

In Claude Desktop or Cursor, you can use these commands as slash commands.

### Task-Oriented Commands (推奨)

タスクタイプとチェーンを指定してコンテキストを自動ロード：

```
# バグ修正
/task-bug-fix Issue #123 Chain: BTC
/task-bug-fix GetAddressInfo error Chain: BCH
/task-bug-fix Database timeout

# 新機能追加
/task-feature-add MuSig2署名集約 Chain: BTC
/task-feature-add ERC-20残高確認 Chain: ETH
/task-feature-add HD Wallet改善

# リファクタリング
/task-refactoring Use Case層への移動 Chain: BTC
/task-refactoring インターフェース抽出

# DB変更
/task-db-change labelカラム追加
/task-db-change 新テーブル作成
```

### Issue/PR Commands

従来のIssue/PR対応コマンド：

- `/fix-issue #123` - Fix issue #123
- `/fix-linter` - Fix linter errors
- `/fix-pr-review #456` - Address review comments for PR #456

### Other Commands

- `/create-github-issue` - Create a new GitHub issue (deprecated, use Agent Skill)
- `/convert-custom-slash-for-codex fix-issue` - Convert fix-issue command to Codex format

## Common Workflow Steps

All commands follow common workflow steps defined in [Workflow Guidelines](../../agents/workflow.md):

- **Required Tools and Versions**: See [Required Tools and Versions](../../agents/requirements.md) - Essential tools (Git, GitHub CLI, Go 1.25.5) and development tools (Atlas v1.0.0, golangci-lint v2.7.2) with version requirements
- **Pre-Flight Checks**: Git status, branch verification, project guidelines review
- **Safety Rules**: Critical rules for security and code quality
- **Verification Steps**: Commands to run before committing (`make lint-fix`, `make tidy`, `make check-build`, `make gotest`)
- **Special Considerations**: Security-sensitive changes, auto-generated files, breaking changes

Each command file links to these common steps to avoid duplication and ensure consistency.

**Important**: Always verify tool versions before starting work. Using incorrect versions (especially Atlas v1.0.0) may cause compatibility issues.

## Cursor Rules Integration

Cursor を使用している場合、`.cursor/rules/task-context-loading.mdc` により自動的にタスクタイプとチェーンが判定され、適切なコンテキストがロードされます。

### 自動判定されるキーワード

| キーワード | Task Type |
|-----------|-----------|
| バグ, 修正, fix, error, Issue # | bug-fix |
| 追加, 実装, 機能, feature, add | feature-add |
| リファクタリング, 整理, 移動 | refactoring |
| スキーマ, DB, テーブル, カラム | db-change |

| キーワード | Chain |
|-----------|-------|
| Bitcoin, BTC, Taproot, Descriptor | BTC |
| Bitcoin Cash, BCH, CashAddr | BCH |
| Ethereum, ETH, ERC-20, Gas | ETH |
| Ripple, XRP, Destination Tag | XRP |

See [.cursor/rules/task-context-loading.mdc](../../.cursor/rules/task-context-loading.mdc) for details.

## Project Guidelines

Each command operates in accordance with the project guidelines ([AGENTS.md](../../AGENTS.md)), following:

- Clean Architecture principles
- Security best practices
- Code quality standards
- Testing requirements

## Related Documentation

- [Task-Oriented Context Management](../../docs/ai-agents/task-oriented-context.md) - コンテキスト管理の概念と戦略
- [Task Contexts](../../docs/ai-agents/task-contexts/README.md) - タスクコンテキスト一覧
- [Agent Skills](../../docs/ai-agents/agent-skills.md) - Agent Skills使用ガイド
