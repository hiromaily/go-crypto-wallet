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
	keygen key create seed
	#seed: 00ySYFfazp+41jyOuLxFb2tWNfIGRmDpGFOBLrneuoQ=

# create hdkey for acounts
.PHONY: create-hdkey
create-hdkey:
	keygen create hdkey -account client -keynum 10
	keygen create hdkey -account receipt -keynum 10
	keygen create hdkey -account payment -keynum 10

###############################################################################
# import private key to keygen wallet
###############################################################################
.PHONY: import-privkey
import-privkey:
	keygen import privkey -account client
	keygen import privkey -account receipt
	keygen import privkey -account payment

###############################################################################
# export public key as csv file
###############################################################################
.PHONY: export-pubkey
export-pubkey:
	keygen export address -account client
	keygen export address -account receipt
	keygen export address -account payment

###############################################################################
# import multisig address from csv file
###############################################################################
.PHONY: import-multisig
import-multisig:
	keygen import multisig -account receipt
	keygen import multisig -account payment

###############################################################################
# sign on unsigned transaction as first signature
#  multisig requireds multiple signature
###############################################################################
#make filepath=./data/tx/receipt/receipt_8_unsigned_1534832793024491932 sign-unsignedtx
.PHONY: sign-unsignedtx
sign-unsignedtx:
	keygen sign file ${filepath}

# [coldwallet]出金用に未署名のトランザクションに署名する #出金時の署名は2回
#sign-payment1: bld
#	coldwallet1 -s -m 1 -i ./data/tx/payment/payment_3_unsigned_1534832966995082772
#
#sign-payment2: bld
#	coldwallet2 -s -m 1 -i ./data/tx/payment/payment_3_unsigned_1534832966995082772
