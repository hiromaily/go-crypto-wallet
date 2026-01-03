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

# create Codex custom prompts from Claude custom slash commands
PHONY: create-codex-custom-prompts
create-codex-custom-prompts:
	@mkdir -p $$HOME/.codex/prompts
	@for file in .claude/commands/*.md; do \
		if [ -f "$$file" ]; then \
			filename=$$(basename "$$file"); \
			cp "$$file" "$$HOME/.codex/prompts/$$filename"; \
			echo "Copied: $$file -> $$HOME/.codex/prompts/$$filename"; \
		fi; \
	done
