package xrp

// MaxLedgerVersionOffset is the maximum ledger version offset for XRP transactions.
// 1 ledger version takes approximately 4 seconds, so 15 means 60 seconds.
const MaxLedgerVersionOffset uint64 = 15
