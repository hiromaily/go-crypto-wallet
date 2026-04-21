## Why DevContainer?

### Benefits for All Developers

1. **Consistent Environment**: Everyone uses the same Go version, tools, and dependencies
2. **Isolated Development**: Project dependencies don't affect your host system
3. **Quick Setup**: New team members can start coding in minutes
4. **Clean Uninstall**: Remove the container to completely clean up the development environment

### Benefits for AI-Assisted Development

**DevContainer is particularly valuable when working with AI coding assistants like Claude Code, GitHub Copilot, or other AI tools:**

1. **Safety and Isolation**
   - AI tools may occasionally generate code that modifies or deletes files
   - DevContainer provides a sandboxed environment, protecting your host system
   - If something goes wrong, simply rebuild the container

2. **Consistent AI Context**
   - All team members and AI tools work with the same environment
   - AI-generated code works consistently across different machines
   - Reduces "works on my machine" issues with AI-suggested solutions

3. **Easy Rollback**
   - Quickly reset to a clean state if AI changes cause issues
   - Container snapshots allow experimentation without risk
   - Git combined with container isolation provides double safety

4. **Reproducible AI Sessions**
   - AI tools can rely on specific tool versions
   - Build tags and linter configurations are pre-configured
   - Reduces friction when AI suggests commands or configurations

---
