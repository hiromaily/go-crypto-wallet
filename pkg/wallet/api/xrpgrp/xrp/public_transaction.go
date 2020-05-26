package xrp

// https://xrpl.org/transaction-methods.html
// 1. Connect to a Test Net Server
// 2. Prepare Transaction
//{
//  "TransactionType": "Payment",
//  "Account": "rPT1Sjq2YGrBMTttX4GZHjKu9dyfzbpAYe",
//  "Amount": "2000000",
//  "Destination": "rUCzEr6jrEyMpjhs4wSdQdz4g8Y382NxfM"
//}
// Transaction Common Fields
// - https://xrpl.org/transaction-common-fields.html
// - [Required]
//   - Account
//   - TransactionType
//   - Fee (auto-fillable)
//   - Sequence (auto-fillable)
//   - LastLedgerSequence (strongly recommended)

// Is there alternative way of prepareTransaction()??
// - LastLedgerSequence look important, how is it retrieved?
// - ledger command https://xrpl.org/ledger.html
// - Is `ledger_index` LastLedgerSequence??

// As one of workaround
// - create node.js server to run ripple-lib
// - golang client connects to node server by unix-domain socket
