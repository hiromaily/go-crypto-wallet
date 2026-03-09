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
    '**/issues/**',         // Internal refactoring checklists/plans
    '**/overview-ja.md',    // Japanese-language overview (English version exists)
    '**/archive/**',        // Archived/historical BTC operation docs
  ],

  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [],

    sidebar: {},

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
