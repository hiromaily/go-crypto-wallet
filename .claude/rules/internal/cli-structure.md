---
paths:
  [
    "internal/interface-adapters/cli/watch/**/*.go",
    "internal/interface-adapters/cli/keygen/**/*.go",
    "internal/interface-adapters/cli/sign/**/*.go",
  ]
---

# CLI Subcommand Structure Rules

> **Current command tree and Command × Chain × UseCase matrix (SSOT)**: [`internal/interface-adapters/cli/README.md`](../../../internal/interface-adapters/cli/README.md)

## Rule 1: Command Hierarchy Order

CLI subcommands follow the order: **[wallet type] → [verb] → [target]**

```
<wallet>  <verb>     <target>
keygen    create     key
keygen    sign       tx
watch     import     address
watch     send       tx
```

**Verb Layer** contains action-oriented subcommands (e.g., `create`, `sign`, `import`, `send`).

**Target Layer** contains the object of the action (e.g., `key`, `tx`, `address`, `multisig`).

- `multisig` is a **target** — it lives at the target layer under a verb (e.g., `watch send multisig`), not as a top-level verb.
- Chain-specific API calls (e.g., `api btc`, `api eth`) are an **exception**: they are allowed as verbs to directly invoke chain node RPC commands. This exception applies only to low-level node API wrappers that do not fit naturally into the generic verb/target hierarchy.

```
# Standard pattern
keygen  create  key
keygen  create  multisig   ← multisig is a target

# Exception: chain API commands act as verbs
keygen  api  btc  <rpc-command>
watch   api  eth  <rpc-command>
```

## Rule 2: Unsupported Coin Type Handling

Some CLI commands are only valid for specific chains. When a command is called with an unsupported coin type, **do not silently skip** — output a human-readable message indicating the command is recognized but not supported for that chain.

```go
// ✅ GOOD: explicit unsupported message
switch coinType {
case coin.BTC, coin.BCH:
    cmd.AddCommand(newDescriptorCmd(btcClient))
default:
    fmt.Printf("[WARN] descriptor command is not supported for coin type: %s\n", coinType)
}
```

```go
// ❌ BAD: silent no-op
switch coinType {
case coin.BTC, coin.BCH:
    cmd.AddCommand(newDescriptorCmd(btcClient))
// other coins simply get nothing — no feedback
}
```

The message format should follow the pattern:
```
[WARN] <command> command is not supported for coin type: <coinType>
```
