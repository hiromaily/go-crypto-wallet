# Bitcoin Core 0.19.0 RPC
[RPC 0.19.0](https://bitcoincore.org/en/doc/0.19.0/)

## Blockchain
#### [getbestblockhash](https://bitcoincore.org/en/doc/0.19.0/rpc/blockchain/getbestblockhash/)
`getbestblockhash`  

Returns the hash of the best (tip) block in the most-work fully-validated chain.

#### [getblock](https://bitcoincore.org/en/doc/0.19.0/rpc/blockchain/getblock/)
`getblock "blockhash" ( verbosity )`

If verbosity is 0, returns a string that is serialized, hex-encoded data for block 'hash'.  
If verbosity is 1, returns an Object with information about block <hash>.  
If verbosity is 2, returns an Object with information about block <hash> and information about each transaction.  

 
## Control

## Generating

## Mining

## Network

## RawTransactions

## Util

## Wallet

## ZMQ
