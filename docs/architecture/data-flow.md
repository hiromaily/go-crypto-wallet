<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/architecture/data-flow.tpl.md · Run `make docs` to regenerate.
-->

# Data Flow Examples

## Creating and Signing a Transaction

```
1. Watch Wallet (Online)
   └── Create unsigned transaction
   └── Export to file

2. Keygen Wallet (Offline)
   └── Import unsigned transaction
   └── Sign (first signature for multisig)
   └── Export partially signed transaction

3. Sign Wallet (Offline)
   └── Import partially signed transaction
   └── Sign (additional signatures)
   └── Export fully signed transaction

4. Watch Wallet (Online)
   └── Import signed transaction
   └── Broadcast to network
```
