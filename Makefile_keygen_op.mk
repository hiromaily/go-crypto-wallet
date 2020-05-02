###############################################################################
# Keygen Wallet for generating keys
#  and sigunature for unsigned transaction for multisig
###############################################################################
###############################################################################
# create hd key
###############################################################################
# create seed
.PHONY: create-seed
create-seed:
	keygen create seed
	#seed: 00ySYFfazp+41jyOuLxFb2tWNfIGRmDpGFOBLrneuoQ=

# create hdkey for acounts
.PHONY: create-hdkey
create-hdkey:
	keygen create hdkey -account client -keynum 10
	keygen create hdkey -account deposit -keynum 10
	keygen create hdkey -account payment -keynum 10
	keygen create hdkey -account stored -keynum 10

###############################################################################
# import private key
###############################################################################
.PHONY: import-privkey
import-privkey:
	keygen import privkey -account client
	keygen import privkey -account deposit
	keygen import privkey -account payment
	keygen import privkey -account stored

###############################################################################
# import full-pubkey
###############################################################################
#make filepath=./data/pubkey/auth1_1588399093997165000.csv import-fullpubkey
.PHONY: import-fullpubkey
import-fullpubkey:
	keygen import fullpubkey -file ${filepath}
	keygen import fullpubkey -file ${filepath}
	keygen import fullpubkey -file ${filepath}
	keygen import fullpubkey -file ${filepath}
	keygen import fullpubkey -file ${filepath}

###############################################################################
# create multisig address
###############################################################################
.PHONY: create-multisig
create-multisig:
	keygen create multisig -account deposit
	keygen create multisig -account payment
	keygen create multisig -account stored

###############################################################################
# export address
###############################################################################
.PHONY: export-address
export-address:
	keygen export address -account client
	keygen export address -account deposit
	keygen export address -account payment
	keygen export address -account stored

###############################################################################
# sign on unsigned transaction as first signature
#  multisig requireds multiple signature
###############################################################################
#make filepath=./data/tx/deposit/deposit_8_unsigned_1534832793024491932 sign-unsignedtx
.PHONY: sign-unsignedtx
sign-unsignedtx:
	keygen sign file ${filepath}
