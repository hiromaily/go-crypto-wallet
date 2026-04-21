# E2E Transaction Patterns Guide

This document explains transaction combination patterns for Bitcoin/Bitcoin Cash. Various E2E workflow patterns exist depending on key types and whether multisig is used.

## Table of Contents

1. [Overview](#overview)
2. [Supported Key Types](#supported-key-types)
3. [Signature Patterns](#signature-patterns)
4. [E2E Workflow Matrix](#e2e-workflow-matrix)
5. [Details of Each Pattern](#details-of-each-pattern)
6. [Account Types and Signing Requirements](#account-types-and-signing-requirements)
7. [Implementation Status](#implementation-status)
8. [E2E Script Reference](#e2e-script-reference)

---

## Overview

Bitcoin transactions can be classified along two main axes:

1. **Key Type (Address Type)** - Which BIP standard is used to generate addresses
2. **Signature Pattern** - Single-sig or multisig

The combination of these creates various E2E workflows.

---
