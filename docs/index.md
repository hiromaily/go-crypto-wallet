---
layout: home

hero:
  name: go-crypto-wallet
  text: Multi-Chain Cryptocurrency Wallet
  tagline: A production-grade wallet supporting BTC, BCH, ETH, XRP, and Cosmos with Clean Architecture and threshold signing.
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started/installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/hiromaily/go-crypto-wallet

features:
  - title: Architecture
    details: Clean Architecture with strict layer separation. Domain layer has zero infrastructure dependencies. Designed for offline key generation and signing.
    link: /guidelines/architecture
    linkText: Learn more

  - title: Multi-Chain Support
    details: Bitcoin (BTC) with Taproot, Descriptor, PSBT, MuSig2. Bitcoin Cash (BCH), Ethereum (ETH) with ERC-20 and MPC-TSS. XRP and Cosmos.
    link: /chains/btc/
    linkText: Explore chains

  - title: Getting Started
    details: Install dependencies, configure wallets, and run CLI commands for key generation, signing, and transaction watching.
    link: /getting-started/installation
    linkText: Install now

  - title: Development Guidelines
    details: Coding conventions, testing strategy, security rules, workflow, and release process for contributors.
    link: /guidelines/
    linkText: Read guidelines

  - title: Database
    details: PostgreSQL and SQLite via Atlas migrations and SQLC code generation. Schema change workflow with migration flow documentation.
    link: /database/architecture
    linkText: View database docs

  - title: AI Agent & Dev Workflow
    details: Kiro spec-driven development, agent skills, task contexts, and Claude memory integration for AI-assisted development.
    link: /ai/agent-skills
    linkText: Explore AI workflow
---
