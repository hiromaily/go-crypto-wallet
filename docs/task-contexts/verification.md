---
task_type: verification
description: ファイルタイプ別の検証コマンドマトリックス
version: 2.0.0
---

# Verification Matrix

## SSOT Reference

ファイルタイプ別の検証コマンドは `.claude/rules/` に定義されています（SSOT）。

| File Type | Rule File | Required Commands |
|-----------|-----------|-------------------|
| `*.go` | `.claude/rules/go/` | `make go-lint`, `make check-build` |
| `*.ts`, `*.js` | `.claude/rules/typescript.md` | `yarn lint` / `npm run lint` |
| `*.sh` | `.claude/rules/shell-script.md` | `make shfmt` |
| `*.sql` | `.claude/rules/sql.md` | `make sqlc-validate` |
| `*.hcl` | `.claude/rules/hcl.md` | `make atlas-fmt`, `make atlas-lint` |
| `*.proto` | `.claude/rules/proto.md` | `make lint-proto` |
| `*.yaml` | `.claude/rules/yaml.md` | `make yaml-lint` |
| `*.md` | - | (none required) |

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

- [bug-fix.md](./bug-fix.md)
- [feature-add.md](./feature-add.md)
- [refactoring.md](./refactoring.md)
- [db-change.md](./db-change.md)
- [documentation.md](./documentation.md)

## Related Documents

- [Task Contexts README](./README.md) - タスクコンテキスト一覧
- [Coding Conventions](../guidelines/coding-conventions.md) - コーディング規約
