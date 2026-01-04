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
