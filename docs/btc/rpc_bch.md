# Bitcoin ABC 0.21.0 RPC
[RPC 0.21.0](https://github.com/Bitcoin-ABC/bitcoin-abc/tree/master/doc)

```
$ bitcoin-cli help

== Avalanche ==
addavalanchepeer nodeid "publickey"
getavalanchekey

== Blockchain ==
getbestblockhash
getblock "blockhash" ( verbosity )
getblockchaininfo
getblockcount
getblockhash height
getblockheader "blockhash" ( verbose )
getblockstats hash_or_height ( stats )
getchaintips
getchaintxstats ( nblocks "blockhash" )
getdifficulty
getmempoolancestors "txid" ( verbose )
getmempooldescendants "txid" ( verbose )
getmempoolentry "txid"
getmempoolinfo
getrawmempool ( verbose )
gettxout "txid" n ( include_mempool )
gettxoutproof ["txid",...] ( "blockhash" )
gettxoutsetinfo
preciousblock "blockhash"
pruneblockchain height
savemempool
scantxoutset "action" [scanobjects,...]
verifychain ( checklevel nblocks )
verifytxoutproof "proof"

== Control ==
getmemoryinfo ( "mode" )
getrpcinfo
help ( "command" )
logging ( "include" "exclude" )
stop
uptime

== Generating ==
generate nblocks ( maxtries )
generatetoaddress nblocks "address" ( maxtries )

== Mining ==
getblocktemplate ( "template_request" )
getmininginfo
getnetworkhashps ( nblocks height )
prioritisetransaction "txid" dummy fee_delta
submitblock "hexdata" ( "dummy" )
submitheader "hexdata"

== Network ==
addnode "node" "command"
clearbanned
disconnectnode ( "address" nodeid )
getaddednodeinfo ( "node" )
getconnectioncount
getexcessiveblock
getnettotals
getnetworkinfo
getnodeaddresses ( count )
getpeerinfo
listbanned
ping
setban "subnet" "command" ( bantime absolute )
setexcessiveblock blockSize
setnetworkactive state

== Rawtransactions ==
combinepsbt ["psbt",...]
combinerawtransaction ["hexstring",...]
converttopsbt "hexstring" ( permitsigdata )
createpsbt [{"txid":"hex","vout":n,"sequence":n},...] [{"address":amount},{"data":"hex"},...] ( locktime )
createrawtransaction [{"txid":"id","vout":n},...] [{"address":amount},{"data":"hex"},...] ( locktime )
decodepsbt "psbt"
decoderawtransaction "hexstring"
decodescript "hexstring"
finalizepsbt "psbt" ( extract )
fundrawtransaction "hexstring" ( options )
getrawtransaction "txid" ( verbose "blockhash" )
sendrawtransaction "hexstring" ( allowhighfees )
signrawtransactionwithkey "hexstring" ["privatekey",...] ( [{"txid":"hex","vout":n,"scriptPubKey":"hex","redeemScript":"hex","amount":amount},...] "sighashtype" )
testmempoolaccept ["rawtxs"] ( allowhighfees )

== Util ==
createmultisig nrequired ["key",...]
estimatefee
signmessagewithprivkey "privkey" "message"
validateaddress "address"
verifymessage "address" "signature" "message"

== Wallet ==
abandontransaction "txid"
abortrescan
addmultisigaddress nrequired ["key",...] ( "label" )
backupwallet "destination"
createwallet "wallet_name" ( disable_private_keys blank )
dumpprivkey "address"
dumpwallet "filename"
encryptwallet "passphrase"
getaddressesbylabel "label"
getaddressinfo "address"
getbalance ( "dummy" minconf include_watchonly )
getnewaddress ( "label" )
getrawchangeaddress
getreceivedbyaddress "address" ( minconf )
getreceivedbylabel "label" ( minconf )
gettransaction "txid" ( include_watchonly )
getunconfirmedbalance
getwalletinfo
importaddress "address" ( "label" rescan p2sh )
importmulti "requests" ( "options" )
importprivkey "privkey" ( "label" rescan )
importprunedfunds "rawtransaction" "txoutproof"
importpubkey "pubkey" ( "label" rescan )
importwallet "filename"
keypoolrefill ( newsize )
listaddressgroupings
listlabels ( "purpose" )
listlockunspent
listreceivedbyaddress ( minconf include_empty include_watchonly "address_filter" )
listreceivedbylabel ( minconf include_empty include_watchonly )
listsinceblock ( "blockhash" target_confirmations include_watchonly include_removed )
listtransactions ( "label" count skip include_watchonly )
listunspent ( minconf maxconf ["address",...] include_unsafe query_options )
listwalletdir
listwallets
loadwallet "filename"
lockunspent unlock ( [{"txid":"hex","vout":n},...] )
removeprunedfunds "txid"
rescanblockchain ( start_height stop_height )
sendmany "" {"address":amount} ( minconf "comment" ["address",...] )
sendtoaddress "address" amount ( "comment" "comment_to" subtractfeefromamount )
sethdseed ( newkeypool "seed" )
setlabel "address" "label"
settxfee amount
signmessage "address" "message"
signrawtransactionwithwallet "hexstring" ( [{"txid":"hex","vout":n,"scriptPubKey":"hex","redeemScript":"hex","amount":amount},...] "sighashtype" )
unloadwallet ( "wallet_name" )
walletcreatefundedpsbt [{"txid":"hex","vout":n,"sequence":n},...] [{"address":amount},{"data":"hex"},...] ( locktime {"changeAddress":"str","changePosition":n,"includeWatching":bool,"lockUnspents":bool,"feeRate":amount,"subtractFeeFromOutputs":[int,...]} bip32derivs )
walletprocesspsbt "psbt" ( sign "sighashtype" bip32derivs )

== Zmq ==
getzmqnotifications
```