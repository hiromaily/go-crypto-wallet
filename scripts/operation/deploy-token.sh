#!/bin/sh

set -eu

# https://github.com/hiromaily/erc20-token
git clone https://github.com/hiromaily/erc20-token.git
cd erc20-token

yarn install
yarn run deploy       # using 7545 port
#yarn run deploy-geth  # using 8545 port


# then takes note
# - contract address
# - account address would be account[0] in ganache which has tokens

# 2_all_contracts.js
# ==================
#
#    Deploying 'HyToken'
#    -------------------
#    > transaction hash:    0x047944b153bd8e9775594652f7c2a3d7537e9f82807b1797ef0836613f93434d
#    > Blocks: 0            Seconds: 0
#    > contract address:    0x9bAB891D3EE061395dE9a86210000508F644D89b
#    > block number:        3
#    > block timestamp:     1632368218
#    > account:             0xF7763dFB4eDeCd5854125c8dCa3531aed8077e55
#    > balance:             99.9487734
#    > gas used:            2269975 (0x22a317)
#    > gas price:           20 gwei
#    > value sent:          0 ETH
#    > total cost:          0.0453995 ETH