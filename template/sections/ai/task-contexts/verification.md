---
task_type: verification
description: ファイルタイプ別の検証コマンドマトリックス
version: 2.0.0
---

# Verification Matrix

## SSOT Reference

ファイルタイプ別の検証コマンドは `.claude/rules/` に定義されています（SSOT）。  
完全なルールファイルと検証コマンドの一覧は [Coding Conventions — Language-Specific Rules](../../../../docs/guidelines/coding-conventions.md#language-specific-rules) を参照してください。

## Task Type × File Type Matrix

タスクタイプとファイルタイプの組み合わせによる検証コマンド:

| Task Type | Go Files | MD Files | SQL/HCL | Config |
|-----------|----------|----------|---------|--------|
| bug-fix | lint, build, test | (none) | atlas-fmt | (none) |
| feature-add | lint, build, test | (none) | atlas-* | (none) |
| refactoring | lint, build, test | (none) | atlas-fmt | (none) |
| db-change | lint, build | (none) | atlas-*, sqlc | (none) |
| documentation | (none) | (optional) | (none) | (none) |

## Skip Verification Scenarios

以下のシナリオでは検証をスキップまたは最小化できます：

| Scenario | Skippable Commands | Reason |
|----------|-------------------|--------|
| ドキュメントのみ変更 | 全ての Go 関連 | コード変更なし |
| Config のみ変更 | 全ての Go 関連 | コード変更なし |
| コメントのみ追加 | `gotest` | 機能変更なし |
| typo 修正 | `gotest` | 機能変更なし |

## Integration with Task Contexts

各タスクコンテキストファイルは、この Verification Matrix と rules を参照しています：

- [bug-fix.md](../../../../docs/ai/task-contexts/bug-fix.md)
- [feature-add.md](../../../../docs/ai/task-contexts/feature-add.md)
- [refactoring.md](../../../../docs/ai/task-contexts/refactoring.md)
- [db-change.md](../../../../docs/ai/task-contexts/db-change.md)
- [documentation.md](../../../../docs/ai/task-contexts/documentation.md)

## Related Documents

- [Task Contexts README](../../../../docs/ai/task-contexts/README.md) - タスクコンテキスト一覧
- [Coding Conventions](../../../../docs/guidelines/coding-conventions.md) - コーディング規約
