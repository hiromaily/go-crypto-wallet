###############################################################################
# AI Targets
###############################################################################
# create symlink from AGENTS.md for CLAUDE.md
PHONY: create-claudemd-symlink
create-claudemd-symlink:
	@find . -name "AGENTS.md" -type f | while read agents_file; do \
		dir=$$(dirname "$$agents_file"); \
		claude_file="$$dir/CLAUDE.md"; \
		ln -sf "AGENTS.md" "$$claude_file"; \
		echo "Created symlink: $$claude_file -> $$dir/AGENTS.md"; \
	done


# Install Codex custom prompts from .codex/prompts/ to ~/.codex/prompts/
#
# Usage workflow:
# 1. Create Claude custom slash command in .claude/commands/ (manual step)
# 2. Convert to Codex format using Claude's convert-custom-slash-for-codex command:
#    /convert-custom-slash-for-codex <command-name>
#    This generates Codex-optimized prompt files in .codex/prompts/
# 3. Install prompts to Codex directory:
#    make install-codex-custom-prompts
#    This copies all .md files from .codex/prompts/ to ~/.codex/prompts/
.PHONY: install-codex-custom-prompts
install-codex-custom-prompts:
	@mkdir -p $$HOME/.codex/prompts
	@for file in .codex/prompts/*.md; do \
		if [ -f "$$file" ]; then \
			filename=$$(basename "$$file"); \
			cp "$$file" "$$HOME/.codex/prompts/$$filename"; \
			echo "Copied: $$file -> $$HOME/.codex/prompts/$$filename"; \
		fi; \
	done
	@echo "Codex custom prompts installed to $$HOME/.codex/prompts/"

# add Notion MCP by Claude
.PHONY: add-notion-mcp-by-claude
add-notion-mcp-by-claude:
	claude mcp add notion --scope project --transport http https://mcp.notion.com/mcp

###############################################################################
# AI Agent SSOT Sync
###############################################################################
# Sync Claude rules to Cursor rules
# Source: .claude/rules/*.md -> Destination: .cursor/rules/*.mdc
#
# Conversion rules:
#   - paths: present -> globs: + alwaysApply: false
#   - paths: absent  -> alwaysApply: true (global rule)
#   - First # heading -> description:
#   - .md -> .mdc extension
.PHONY: sync-cursor-rules
sync-cursor-rules:
	@./scripts/ai-agent/sync-rule-claude-to-cursor.sh --force --verbose

# Dry-run: preview what would be converted
.PHONY: sync-cursor-rules-dry
sync-cursor-rules-dry:
	@./scripts/ai-agent/sync-rule-claude-to-cursor.sh --dry-run --verbose

# Sync all AI agent configurations (SSOT)
# - Cursor rules: auto-generated from Claude rules
# - Cursor skills: symlink to .claude/skills (manual setup required)
# - Cursor commands: auto-loaded from .claude/commands
.PHONY: sync-ai-agent
sync-ai-agent: sync-cursor-rules
	@echo ""
	@echo "AI Agent configurations synced."
	@echo "Note: .cursor/skills should be a symlink to ../.claude/skills"
	@echo "      Run: ln -sf ../.claude/skills .cursor/skills"
