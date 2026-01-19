# Migration Guide: ripple-lib 1.x to xrpl.js 4.5.0

This guide provides detailed migration instructions for creating a new `xrpl-grpc-server` using xrpl.js 4.5.0 to replace the deprecated `ripple-lib-server`.

## Official Reference

- [xrpl.org Migration Guide](https://xrpl.org/docs/references/xrpljs2-migration-guide)
- [xrpl.js GitHub Releases](https://github.com/XRPLF/xrpl.js/releases)

## Overview

### Tech Stack Changes

| Component | Old (ripple-lib-server) | New (xrpl-grpc-server) |
|-----------|-------------------------|------------------------|
| Runtime | Node.js + ts-node | **Bun** |
| XRP Library | ripple-lib 1.10.1 | **xrpl.js 4.5.0** |
| gRPC | grpc (deprecated) | **ConnectRPC** (@connectrpc/connect) |
| Linter/Formatter | ESLint + Prettier | **Biome** (Rust-based) |
| TypeScript | 4.7.4 | **5.9.3** |
| Proto codegen | grpc_tools_node_protoc_ts | **@bufbuild/protoc-gen-es** + **@connectrpc/protoc-gen-connect-es** |

### Why These Changes?

- **Bun**: Faster runtime, native TypeScript support, better developer experience
- **xrpl.js 4.5.0**: Active development, modern API, better TypeScript types
- **ConnectRPC**: Better Bun compatibility, supports gRPC/gRPC-Web/Connect protocols
- **Biome**: Single tool for linting and formatting, faster (Rust-based)

## Security Warning

> **CRITICAL**: xrpl.js versions **4.2.1-4.2.4** and **2.14.2** were compromised in a supply chain attack. Always use version **4.5.0** or newer.

## API Migration Reference

### Class and Module Changes

| ripple-lib 1.x | xrpl.js 4.5.0 | Notes |
|----------------|---------------|-------|
| `new RippleAPI({server})` | `new Client(server)` | Class renamed |
| `api.connect()` | `client.connect()` | Same |
| `api.disconnect()` | `client.disconnect()` | Same |
| `api.isConnected()` | `client.isConnected()` | Same |
| `api.getLedgerVersion()` | `client.getLedgerIndex()` | Method renamed |
| `api.generateAddress()` | `Wallet.generate()` | Static method |
| `api.isValidAddress()` | `xrpl.isValidAddress()` | Module function |
| `api.sign(txJSON, secret)` | `wallet.sign(tx)` | Use Wallet instance |
| `api.submit(txBlob)` | `client.submit(txBlob)` | Same signature |
| `api.prepareTransaction()` | `client.autofill()` | Different return format |
| `api.getAccountInfo()` | `client.request({command: "account_info"})` | Raw WebSocket API |
| `api.getTransaction()` | `client.request({command: "tx"})` | Raw WebSocket API |
| `api.combine()` | `xrpl.multisign()` | Module function |
| `api.on("ledger", cb)` | `client.on("ledgerClosed", cb)` | Event renamed, requires subscribe |
| `api.xrpToDrops()` | `xrpl.xrpToDrops()` | Module function |
| `api.dropsToXrp()` | `xrpl.dropsToXrp()` | Module function |

---

## 1. Boilerplate / Server Setup

### Before (ripple-lib 1.x + grpc)

```typescript
import * as grpc from 'grpc';
import * as ripple from 'ripple-lib';

const rippleAPI = new ripple.RippleAPI({server: wsURL});

rippleAPI.on('error', (errorCode, errorMessage) => {
  console.log(errorCode + ': ' + errorMessage);
});

rippleAPI.on('disconnected', (code) => {
  console.log('disconnected, code:', code);
});

await rippleAPI.connect();

const server = new grpc.Server();
server.addService(service, implementation);
server.bindAsync(`0.0.0.0:${port}`, grpc.ServerCredentials.createInsecure(), callback);
server.start();
```

### After (xrpl.js 4.5.0 + ConnectRPC + Bun)

```typescript
import { Client } from "xrpl";
import { createConnectRouter } from "@connectrpc/connect";
import { createBunServeHandler } from "@connectrpc/connect/bun";

const client = new Client(wsURL);

client.on("error", (errorCode, errorMessage) => {
  console.log(`${errorCode}: ${errorMessage}`);
});

client.on("disconnected", (code) => {
  console.log("disconnected, code:", code);
});

await client.connect();

const router = createConnectRouter()
  .service(RippleAccountAPI, accountService)
  .service(RippleAddressAPI, addressService)
  .service(RippleTransactionAPI, transactionService);

Bun.serve({
  port: 50051,
  fetch: createBunServeHandler(router),
});
```

---

## 2. RippleAddressAPI Service

### GenerateAddress

**Before:**

```typescript
generateAddress = (call, callback) => {
  const generated = this.rippleAPI.generateAddress();
  const res = new pb.ResponseGenerateAddress();
  res.setXaddress(generated.xAddress);
  res.setClassicaddress(generated.classicAddress);
  res.setAddress(generated.address);
  res.setSecret(generated.secret);
  callback(null, res);
}
```

**After:**

```typescript
import { Wallet } from "xrpl";

generateAddress: async (req) => {
  const wallet = Wallet.generate();
  return {
    xAddress: wallet.getXAddress(),
    classicAddress: wallet.classicAddress,
    address: wallet.address,
    secret: wallet.seed ?? "",
  };
}
```

### GenerateXAddress

**Before:**

```typescript
generateXAddress = (call, callback) => {
  const generated = this.rippleAPI.generateXAddress();
  const res = new pb.ResponseGenerateXAddress();
  res.setXaddress(generated.xAddress);
  res.setSecret(generated.secret);
  callback(null, res);
}
```

**After:**

```typescript
generateXAddress: async (req) => {
  const wallet = Wallet.generate();
  return {
    xAddress: wallet.getXAddress(),
    secret: wallet.seed ?? "",
  };
}
```

### IsValidAddress

**Before:**

```typescript
isValidAddress = (call, callback) => {
  const address = call.request.getAddress();
  const isValid = this.rippleAPI.isValidAddress(address);
  const res = new pb.ResponseIsValidAddress();
  res.setIsvalid(isValid);
  callback(null, res);
}
```

**After:**

```typescript
import { isValidAddress } from "xrpl";

isValidAddress: async (req) => {
  return {
    isValid: isValidAddress(req.address),
  };
}
```

---

## 3. RippleAccountAPI Service

### GetAccountInfo

**Before:**

```typescript
getAccountInfo = (call, callback) => {
  const address = call.request.getAddress();
  this.rippleAPI.getAccountInfo(address)
    .then(info => {
      const res = new pb.ResponseGetAccountInfo();
      res.setSequence(info.sequence);
      res.setXrpbalance(info.xrpBalance);
      res.setOwnercount(info.ownerCount);
      res.setPreviousaffectingtransactionid(info.previousAffectingTransactionID);
      res.setPreviousaffectingtransactionledgerversion(info.previousAffectingTransactionLedgerVersion);
      callback(null, res);
    })
    .catch(error => callback(error, null));
}
```

**After:**

```typescript
import { Client, dropsToXrp } from "xrpl";

getAccountInfo: async (req) => {
  const response = await client.request({
    command: "account_info",
    account: req.address,
    ledger_index: "validated",
  });

  const accountData = response.result.account_data;
  return {
    sequence: BigInt(accountData.Sequence),
    xrpBalance: dropsToXrp(accountData.Balance),
    ownerCount: BigInt(accountData.OwnerCount),
    previousAffectingTransactionId: accountData.PreviousTxnID ?? "",
    previousAffectingTransactionLedgerVersion: BigInt(accountData.PreviousTxnLgrSeq ?? 0),
  };
}
```

---

## 4. RippleTransactionAPI Service

### PrepareTransaction

**Before:**

```typescript
private async _prepareTransaction(call) {
  const txType = call.request.getTxType();
  const instructions = call.request.getInstructions();

  const paramInst = {};
  if (instructions?.getMaxledgerversionoffset()) {
    paramInst.maxLedgerVersionOffset = instructions.getMaxledgerversionoffset();
  }
  // ... more instruction mapping

  const preparedTx = await this.rippleAPI.prepareTransaction({
    "TransactionType": enumTransactionTypeString[txType],
    "Account": call.request.getSenderaccount(),
    "Amount": this.rippleAPI.xrpToDrops(call.request.getAmount().toString()),
    "Destination": call.request.getReceiveraccount(),
  }, paramInst);

  return preparedTx.txJSON;
}
```

**After:**

```typescript
import { Client, xrpToDrops, type Transaction } from "xrpl";

prepareTransaction: async (req) => {
  const tx = {
    TransactionType: enumTransactionTypeString[req.txType],
    Account: req.senderAccount,
    Amount: xrpToDrops(req.amount.toString()),
    Destination: req.receiverAccount,
  } as Transaction;

  // Handle LastLedgerSequence
  if (req.instructions?.maxLedgerVersionOffset) {
    const vli = await client.getLedgerIndex();
    tx.LastLedgerSequence = vli + Number(req.instructions.maxLedgerVersionOffset);
  }

  const prepared = await client.autofill(tx);

  return {
    txJson: JSON.stringify(prepared),
    instructions: req.instructions,
  };
}
```

### SignTransaction

**Before:**

```typescript
signTransaction = (call, callback) => {
  const signed = this.rippleAPI.sign(call.request.getTxjson(), call.request.getSecret());
  console.log("txID:", signed.id);           // SECURITY: Don't log in production!
  console.log("txBlob:", signed.signedTransaction);

  const res = new pb.ResponseSignTransaction();
  res.setTxid(signed.id);
  res.setTxblob(signed.signedTransaction);
  callback(null, res);
}
```

**After:**

```typescript
import { Wallet, type Transaction } from "xrpl";

signTransaction: async (req) => {
  const wallet = Wallet.fromSeed(req.secret);
  const tx = JSON.parse(req.txJson) as Transaction;
  const signed = wallet.sign(tx);

  return {
    txId: signed.hash,
    txBlob: signed.tx_blob,
  };
}
```

### SubmitTransaction

**Before:**

```typescript
private async _submitTransaction(call) {
  const latestLedgerVersion = await this.rippleAPI.getLedgerVersion();
  const txBlob = call.request.getTxblob();
  const resJSON = await this.rippleAPI.submit(txBlob);
  return { resJSON, earlistLedgerVersion: latestLedgerVersion + 1 };
}
```

**After:**

```typescript
submitTransaction: async (req) => {
  const latestLedgerVersion = await client.getLedgerIndex();
  const response = await client.submit(req.txBlob);

  return {
    resultJsonString: JSON.stringify(response.result),
    earliestLedgerVersion: BigInt(latestLedgerVersion + 1),
  };
}
```

### GetTransaction

**Before:**

```typescript
getTransaction = (call, callback) => {
  const txID = call.request.getTxid();
  const earliestLedgerVersion = call.request.getMinledgerversion();

  this.rippleAPI.getTransaction(txID, {minLedgerVersion: earliestLedgerVersion})
    .then(tx => {
      const res = new pb.ResponseGetTransaction();
      res.setResultjsonstring(JSON.stringify(tx));
      callback(null, res);
    })
    .catch(error => callback(error, null));
}
```

**After:**

```typescript
getTransaction: async (req) => {
  const response = await client.request({
    command: "tx",
    transaction: req.txId,
    min_ledger: Number(req.minLedgerVersion),
  });

  // Check if validated
  if (!response.result.validated) {
    throw new Error("Transaction not yet validated");
  }

  return {
    resultJsonString: JSON.stringify(response.result),
  };
}
```

### CombineTransaction (Multi-signature)

**Before:**

```typescript
combineTransaction = (call, callback) => {
  const signedObj = this.rippleAPI.combine(call.request.getSignedtransactionsList());
  const res = new pb.ResponseCombineTransaction();
  res.setSignedtransaction(signedObj.signedTransaction);
  res.setTxid(signedObj.txJSON);
  callback(null, res);
}
```

**After:**

```typescript
import { multisign, hashes } from "xrpl";

combineTransaction: async (req) => {
  const combined = multisign(req.signedTransactions);
  const txId = hashes.hashSignedTx(combined);

  return {
    signedTransaction: combined,
    txId: txId,
  };
}
```

### WaitValidation (Server Streaming)

**Before:**

```typescript
waitValidation = (call) => {
  const ledgerHandler = (ledger) => {
    if (call.cancelled) {
      call.end();
      this.rippleAPI.removeListener('ledger', ledgerHandler);
      return;
    }
    const res = new pb.ResponseWaitValidation();
    res.setLedgerversion(ledger.ledgerVersion);
    call.write(res);
  };
  this.rippleAPI.on('ledger', ledgerHandler);
}
```

**After:**

```typescript
waitValidation: async function* (req, context) {
  // Must explicitly subscribe to ledger stream
  await client.request({ command: "subscribe", streams: ["ledger"] });

  const ledgerPromise = () => new Promise<number>((resolve) => {
    client.once("ledgerClosed", (ledger) => {
      resolve(ledger.ledger_index);
    });
  });

  try {
    while (!context.signal.aborted) {
      const ledgerIndex = await ledgerPromise();
      yield { ledgerVersion: BigInt(ledgerIndex) };
    }
  } finally {
    await client.request({ command: "unsubscribe", streams: ["ledger"] });
  }
}
```

---

## 5. Important Behavior Changes

### Event Handling

In xrpl.js, you must explicitly subscribe to streams:

```typescript
// Subscribe to ledger events
await client.request({ command: "subscribe", streams: ["ledger"] });

// Now you can listen
client.on("ledgerClosed", (ledger) => {
  console.log(`Ledger #${ledger.ledger_index} closed`);
});

// Unsubscribe when done
await client.request({ command: "unsubscribe", streams: ["ledger"] });
```

### Validated Results

Unlike ripple-lib 1.x which defaulted to validated results, xrpl.js returns current (pending) data by default. Always specify `ledger_index: "validated"` for finalized data:

```typescript
const response = await client.request({
  command: "account_info",
  account: address,
  ledger_index: "validated",  // Important!
});
```

### Offline Operations

The `Wallet.fromSeed()` and `wallet.sign()` operations work offline (no network required). This is important for the sign wallet security model.

---

## 6. New Project Structure

```
apps/xrpl-grpc-server/
├── biome.json              # Biome config (linter + formatter)
├── bun.lock                # Bun lockfile
├── package.json            # Bun-compatible dependencies
├── tsconfig.json           # TypeScript 5.9.3 config
├── buf.gen.yaml            # Buf code generation config
├── buf.yaml                # Buf module config
├── README.md               # Documentation
├── docs/
│   └── MIGRATION-GUIDE.md  # This file
├── src/
│   ├── index.ts            # Entry point
│   ├── server.ts           # ConnectRPC server setup
│   ├── config.ts           # Environment configuration
│   ├── xrpl/
│   │   ├── client.ts       # xrpl.js Client wrapper
│   │   └── wallet.ts       # Wallet utilities
│   ├── services/
│   │   ├── account.ts      # RippleAccountAPI implementation
│   │   ├── address.ts      # RippleAddressAPI implementation
│   │   └── transaction.ts  # RippleTransactionAPI implementation
│   └── gen/                # Generated protobuf/connect code
└── Makefile                # Build/dev commands
```

---

## 7. Dependencies

### package.json

```json
{
  "name": "xrpl-grpc-server",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "bun --hot src/index.ts",
    "build": "bun build src/index.ts --outdir=dist --target=bun",
    "lint": "biome check .",
    "lint:fix": "biome check --write .",
    "format": "biome format --write .",
    "typecheck": "tsc --noEmit",
    "proto": "buf generate"
  },
  "dependencies": {
    "@bufbuild/protobuf": "^2.0.0",
    "@connectrpc/connect": "^2.0.0",
    "xrpl": "4.5.0"
  },
  "devDependencies": {
    "@biomejs/biome": "^1.9.0",
    "@bufbuild/buf": "^1.40.0",
    "@bufbuild/protoc-gen-es": "^2.0.0",
    "@connectrpc/protoc-gen-connect-es": "^2.0.0",
    "typescript": "^5.9.3"
  }
}
```

---

## 8. Verification Commands

```bash
bun run lint      # biome check
bun run format    # biome format
bun run typecheck # tsc --noEmit
bun run dev       # bun --hot src/index.ts
bun run build     # bun build src/index.ts
bun run proto     # buf generate
```

---

## 9. Relationship with ripple-lib-server

The old `apps/ripple-lib-server` is preserved for reference. Once `xrpl-grpc-server` is fully tested and integrated, the old server can be deprecated and eventually removed.

| Directory | Status | Description |
|-----------|--------|-------------|
| `apps/ripple-lib-server` | **Preserved** | Legacy server using ripple-lib 1.x |
| `apps/xrpl-grpc-server` | **New** | Modern server using xrpl.js 4.5.0 |
