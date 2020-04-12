###############################################################################
# Keygen Wallet
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
	keygen key create hdkey -account client -keynum 10
	keygen key create hdkey -account receipt -keynum 10
	keygen key create hdkey -account payment -keynum 10


###############################################################################
# import private key to keygen wallet
###############################################################################
.PHONY: import-priv-key
import-priv-key:
	keygen key import privkey


# Clientのprivate keyをcoldwalletに登録する
add-client-priv-key:
	coldwallet1 -k -m 20

# Receiptのprivate keyをcoldwalletに登録する
add-receipt-priv-key:
	coldwallet1 -k -m 21

# Paymentのprivate keyをcoldwalletに登録する
add-payment-priv-key:
	coldwallet1 -k -m 22




# Clientのpubアドレスをexportする
export-client-pub-key:
	coldwallet1 -k -m 30

# Receiptのpubアドレスをexportする
export-receipt-pub-key:
	coldwallet1 -k -m 31

# Paymentのpubアドレスをexportする
export-payment-pub-key:
	coldwallet1 -k -m 32


# Receiptのmultisigアドレスをimportする
import-receipt-multisig-address:
	coldwallet1 -k -m 40

# Paymentのmultisigアドレスをimportする
import-payment-multisig-address:
	coldwallet1 -k -m 41











#

# [coldwallet] 未署名のトランザクションに署名する
sign: bld
	coldwallet1 -w 1 -s -m 1 -i ./data/tx/receipt/receipt_8_unsigned_1534832793024491932

# [coldwallet]出金用に未署名のトランザクションに署名する #出金時の署名は2回
sign-payment1: bld
	coldwallet1 -s -m 1 -i ./data/tx/payment/payment_3_unsigned_1534832966995082772

sign-payment2: bld
	coldwallet2 -s -m 1 -i ./data/tx/payment/payment_3_unsigned_1534832966995082772





