## Workflow Guidelines

This document describes development workflow, dependency management, and git operations for the go-crypto-wallet project.

### Refactoring Status

- Make changes incrementally without breaking existing functionality
- Run tests before and after refactoring
- Follow the phased approach outlined in refactoring documents

**Important**: This is an ongoing refactoring project moving toward Clean Architecture. When working on code:

1. Check if the code is part of a planned refactoring
2. Follow the refactoring plan if it exists
3. Don't introduce new patterns that conflict with the target architecture
4. When in doubt, ask about refactoring priorities
