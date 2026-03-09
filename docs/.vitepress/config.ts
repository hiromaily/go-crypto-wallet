import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: 'go-crypto-wallet',
  description:
    'A multi-chain cryptocurrency wallet supporting BTC, BCH, ETH, XRP, and Cosmos with Clean Architecture',
  base: '/go-crypto-wallet/',
  cleanUrls: true,

  // Existing docs contain cross-references to files outside docs/ (AGENTS.md, CLAUDE.md,
  // .claude/rules/, scripts/, internal/, etc.). These will be cleaned up in Task 5
  // (SSOT enforcement). Until then, suppress dead-link build failures.
  ignoreDeadLinks: true,

  // Exclude internal investigation/archive documents that are not suitable for the
  // public-facing docs site, or that contain syntax (e.g. ${{ }}) that conflicts
  // with Vue template processing.
  srcExclude: [
    '**/github-actions/**', // Japanese investigation docs with GitHub Actions ${{ }} syntax
    '**/issues/**', // Internal refactoring checklists/plans
    '**/overview-ja.md', // Japanese-language overview (English version exists)
    '**/archive/**', // Archived/historical BTC operation docs
    // Note: chains/btc/keygen/improvements-2025-ja.md (Japanese) is included intentionally
    // alongside the English version as supplementary reading for Japanese-speaking contributors.
  ],

  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config

    // ─── Navigation bar ────────────────────────────────────────────────────────
    nav: [
      { text: 'Overview', link: '/overview' },
      { text: 'Getting Started', link: '/Installation' },
      { text: 'Architecture', link: '/guidelines/architecture' },
      {
        text: 'Chains',
        items: [
          { text: 'Bitcoin (BTC)', link: '/chains/btc/' },
          { text: 'Bitcoin Cash (BCH)', link: '/chains/bch/' },
          { text: 'Ethereum (ETH)', link: '/chains/eth/' },
          { text: 'XRP', link: '/chains/xrp/' },
          { text: 'Cosmos', link: '/chains/cosmos/' },
        ],
      },
      { text: 'Database', link: '/database/architecture' },
      { text: 'Guidelines', link: '/guidelines/' },
      { text: 'AI Agent', link: '/agent-skills' },
    ],

    // ─── Sidebars (path-keyed; longest-prefix match wins) ──────────────────────
    sidebar: {
      // ── Default sidebar: overview, getting started, AI agent (root-level pages) ──
      '/': [
        {
          text: 'Overview',
          items: [
            { text: 'Project Overview', link: '/overview' },
            { text: 'Transaction Flow', link: '/transaction-flow' },
            { text: 'Directory Structure', link: '/directory_structure' },
          ],
        },
        {
          text: 'Getting Started',
          items: [
            { text: 'Installation', link: '/Installation' },
            { text: 'CLI Commands', link: '/commands' },
            { text: 'Dev Container', link: '/devcontainer' },
            { text: 'Protocol Buffers', link: '/proto' },
          ],
        },
        {
          text: 'AI Agent & Dev Workflow',
          items: [
            { text: 'Agent Skills', link: '/agent-skills' },
            { text: 'AI Agent Design', link: '/design/ai-agents-instruction' },
            { text: 'Task Contexts', link: '/task-contexts/' },
          ],
        },
        {
          text: 'Tools',
          items: [{ text: 'golangci-lint', link: '/tools/golangci-lint' }],
        },
      ],

      // ── Guidelines ────────────────────────────────────────────────────────────
      '/guidelines/': [
        {
          text: 'Guidelines',
          items: [
            { text: 'Overview', link: '/guidelines/' },
            { text: 'Architecture', link: '/guidelines/architecture' },
            { text: 'Core Principles', link: '/guidelines/core' },
            { text: 'Multi-Chain Design', link: '/guidelines/multi-chain' },
            { text: 'Coding Conventions', link: '/guidelines/coding-conventions' },
            { text: 'Testing', link: '/guidelines/testing' },
            { text: 'Workflow', link: '/guidelines/workflow' },
            { text: 'Security', link: '/guidelines/security' },
            { text: 'Release', link: '/guidelines/release' },
            { text: 'Code Generation', link: '/guidelines/code-generation' },
            { text: 'Task Classification', link: '/guidelines/task-classification' },
            { text: 'Requirements', link: '/guidelines/requirements' },
            { text: 'Claude Memory', link: '/guidelines/claude-mem' },
          ],
        },
      ],

      // ── Bitcoin (BTC) ─────────────────────────────────────────────────────────
      '/chains/btc/': [
        {
          text: 'Bitcoin (BTC)',
          items: [
            { text: 'README', link: '/chains/btc/' },
            { text: 'Architecture', link: '/chains/btc/architecture' },
          ],
        },
        {
          text: 'Overview',
          collapsed: false,
          items: [
            { text: 'Overview', link: '/chains/btc/overview/' },
            { text: 'Address Types', link: '/chains/btc/overview/address-types' },
            { text: 'Technical Reference', link: '/chains/btc/overview/technical-reference' },
          ],
        },
        {
          text: 'Taproot',
          collapsed: false,
          items: [
            { text: 'README', link: '/chains/btc/taproot/' },
            { text: 'User Guide', link: '/chains/btc/taproot/user-guide' },
            { text: 'Testing', link: '/chains/btc/taproot/testing' },
          ],
        },
        {
          text: 'Descriptor Wallets',
          collapsed: false,
          items: [
            { text: 'README', link: '/chains/btc/descriptor/' },
            { text: 'Architecture', link: '/chains/btc/descriptor/architecture' },
            { text: 'API', link: '/chains/btc/descriptor/api' },
            { text: 'User Guide', link: '/chains/btc/descriptor/user-guide' },
            { text: 'Development', link: '/chains/btc/descriptor/development' },
            { text: 'Examples', link: '/chains/btc/descriptor/examples' },
            { text: 'Migration', link: '/chains/btc/descriptor/migration' },
            { text: 'Compatibility', link: '/chains/btc/descriptor/compatibility' },
          ],
        },
        {
          text: 'PSBT',
          collapsed: false,
          items: [
            { text: 'README', link: '/chains/btc/psbt/' },
            { text: 'Implementation', link: '/chains/btc/psbt/implementation' },
            { text: 'User Guide', link: '/chains/btc/psbt/user-guide' },
            { text: 'Developer Guide', link: '/chains/btc/psbt/developer-guide' },
            { text: 'Migration', link: '/chains/btc/psbt/migration' },
            { text: 'PoC Example', link: '/chains/btc/psbt/poc-example' },
          ],
        },
        {
          text: 'MuSig2',
          collapsed: false,
          items: [
            { text: 'README', link: '/chains/btc/musig2/' },
            { text: 'Architecture', link: '/chains/btc/musig2/architecture' },
            { text: 'User Guide', link: '/chains/btc/musig2/user-guide' },
            { text: 'Security', link: '/chains/btc/musig2/security' },
            { text: 'Migration from Traditional', link: '/chains/btc/musig2/migration-from-traditional' },
          ],
        },
        {
          text: 'Key Generation',
          collapsed: false,
          items: [
            { text: 'README', link: '/chains/btc/keygen/' },
            { text: 'Interface Design', link: '/chains/btc/keygen/interface-design' },
            { text: 'Improvements 2025', link: '/chains/btc/keygen/improvements-2025' },
          ],
        },
        {
          text: 'Operations',
          collapsed: false,
          items: [
            { text: 'README', link: '/chains/btc/operations/' },
            { text: 'Wallet Flow', link: '/chains/btc/operations/wallet-flow' },
            { text: 'E2E Transaction Patterns', link: '/chains/btc/operations/e2e-transaction-patterns' },
            { text: 'Wallet Flow Improvements 2025', link: '/chains/btc/operations/wallet-flow-improvements-2025' },
          ],
        },
        {
          text: 'Testing',
          collapsed: false,
          items: [
            { text: 'README', link: '/chains/btc/testing/' },
            { text: 'Pattern 3 Verification', link: '/chains/btc/testing/pattern3-verification' },
          ],
        },
      ],

      // ── Bitcoin Cash (BCH) ───────────────────────────────────────────────────
      '/chains/bch/': [
        {
          text: 'Bitcoin Cash (BCH)',
          items: [
            { text: 'README', link: '/chains/bch/' },
            { text: 'Interface Separation', link: '/chains/bch/interface-separation' },
          ],
        },
      ],

      // ── Ethereum (ETH) ───────────────────────────────────────────────────────
      '/chains/eth/': [
        {
          text: 'Ethereum (ETH)',
          items: [
            { text: 'README', link: '/chains/eth/' },
            { text: 'Architecture', link: '/chains/eth/architecture' },
            { text: 'Transaction Patterns', link: '/chains/eth/transaction-patterns' },
            { text: 'Solidity Development', link: '/chains/eth/solidity-development' },
            { text: 'JSON Schema', link: '/chains/eth/json-schema' },
            { text: 'Multisig', link: '/chains/eth/multisig' },
          ],
        },
      ],

      // ── XRP ──────────────────────────────────────────────────────────────────
      '/chains/xrp/': [
        {
          text: 'XRP (Ripple)',
          items: [
            { text: 'README', link: '/chains/xrp/' },
            { text: 'Architecture 2026', link: '/chains/xrp/architecture-2026' },
            { text: 'Architecture (gRPC Server)', link: '/chains/xrp/architecture-xrpl-grpc-server-version' },
            { text: 'Transaction Flow', link: '/chains/xrp/transaction-flow' },
            { text: 'Library Selection', link: '/chains/xrp/library-selection' },
            { text: 'Network (Devnet)', link: '/chains/xrp/network-devnet' },
            { text: 'Testing Strategy', link: '/chains/xrp/testing-strategy' },
            { text: 'Docker Compose Setup', link: '/chains/xrp/setup-docker-compose-standalone-xrpl' },
            { text: 'xrpl-go Library', link: '/chains/xrp/xrpl-go' },
          ],
        },
      ],

      // ── Cosmos ───────────────────────────────────────────────────────────────
      '/chains/cosmos/': [
        {
          text: 'Cosmos',
          items: [{ text: 'README', link: '/chains/cosmos/' }],
        },
      ],

      // ── Database ─────────────────────────────────────────────────────────────
      '/database/': [
        {
          text: 'Database',
          items: [
            { text: 'Architecture', link: '/database/architecture' },
            { text: 'DB Management', link: '/database/db-management' },
            { text: 'Atlas Migration Flow', link: '/database/atlas-migration-flow' },
            { text: 'SQLC Code Generation', link: '/database/sqlc-code-generation-flow' },
            { text: 'Schema Changes', link: '/database/schema-changes' },
            { text: 'Quick Reference', link: '/database/quick-reference' },
          ],
        },
      ],

      // ── Design Notes ─────────────────────────────────────────────────────────
      '/design/': [
        {
          text: 'Design Notes',
          items: [
            { text: 'AI Agents Instruction', link: '/design/ai-agents-instruction' },
            { text: 'BTC Network Mode Switching', link: '/design/btc-network-mode-switching' },
            { text: 'Revise DB Atlas SQLC Flow', link: '/design/revise-db-atlas-sqlc-flow' },
            { text: 'Superpowers Integration', link: '/design/superpowers-integration' },
            { text: 'PG2SQLite Alter Table', link: '/design/pg2sqlite-alter-table-support' },
            { text: 'Claude Best Practice', link: '/design/claude-best-practice' },
          ],
        },
      ],

      // ── Task Contexts ─────────────────────────────────────────────────────────
      '/task-contexts/': [
        {
          text: 'AI Agent & Dev Workflow',
          items: [
            { text: 'Agent Skills', link: '/agent-skills' },
            { text: 'AI Agent Design', link: '/design/ai-agents-instruction' },
          ],
        },
        {
          text: 'Task Contexts',
          items: [
            { text: 'Overview', link: '/task-contexts/' },
            { text: 'Task-Oriented Context', link: '/task-contexts/task-oriented-context' },
            { text: 'Task Analysis', link: '/task-contexts/task-analysis' },
            { text: 'Bug Fix', link: '/task-contexts/bug-fix' },
            { text: 'Feature Add', link: '/task-contexts/feature-add' },
            { text: 'Refactoring', link: '/task-contexts/refactoring' },
            { text: 'DB Change', link: '/task-contexts/db-change' },
            { text: 'Documentation', link: '/task-contexts/documentation' },
            { text: 'Testing', link: '/task-contexts/test' },
            { text: 'Chain-Specific', link: '/task-contexts/chain-specific' },
            { text: 'Verification', link: '/task-contexts/verification' },
          ],
        },
      ],

      // ── Tools ─────────────────────────────────────────────────────────────────
      '/tools/': [
        {
          text: 'Tools',
          items: [{ text: 'golangci-lint', link: '/tools/golangci-lint' }],
        },
      ],
    },

    search: {
      provider: 'local',
    },

    socialLinks: [
      {
        icon: 'github',
        link: 'https://github.com/hiromaily/go-crypto-wallet',
      },
    ],
  },
})
