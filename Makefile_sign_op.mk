
###############################################################################
# Sign Wallet for generating key of only authorization account
#  and multisig address,
#  and signature for multisig addresses
###############################################################################
###############################################################################
# create hd key
###############################################################################
# create seed
.PHONY: create-seed-signwallet
create-seed-signwallet:
	sign create seed

# create hdkey for acounts
.PHONY: create-hdkey-auth
create-hdkey-auth:
	sign -wallet sign1 create hdkey
	sign2 -wallet sign2 create hdkey
	sign3 -wallet sign3 create hdkey
	sign4 -wallet sign4 create hdkey
	sign5 -wallet sign5 create hdkey

###############################################################################
# import private key to sign wallet
###############################################################################
.PHONY: import-privkey-auth
import-privkey-auth:
	sign -wallet sign1 import privkey
	sign2 -wallet sign2 import privkey
	sign3 -wallet sign3 import privkey
	sign4 -wallet sign4 import privkey
	sign5 -wallet sign5 import privkey

###############################################################################
# export full-pubkey as csv
###############################################################################
.PHONY: export-multisig
export-multisig:
	sign -wallet sign1 export fullpubkey
	sign2 -wallet sign2 export fullpubkey
	sign3 -wallet sign3 export fullpubkey
	sign4 -wallet sign4 export fullpubkey
	sign5 -wallet sign5 export fullpubkey
