const HyToken = artifacts.require('HyToken');

module.exports = async (deployer, _network, accounts) => {
  const deployAndMinter = accounts[0]; //TODO: when using geth, change address
  await deployer.deploy(HyToken);
  const hyToken = await HyToken.deployed();
  await hyToken.mint(deployAndMinter, 10000, { from: deployAndMinter });
};
// 2_all_contracts.js
// ==================

//    Deploying 'HyToken'
//    -------------------
//    > transaction hash:    0x6b4b3a1bb5c763e549ab00bd07e43dbc55dfff862201c640ae02732de3e24fac
//    > Blocks: 0            Seconds: 0
//    > contract address:    0x014C2061ba81a6Da4b8dD32b1322598D99B711D0
//    > block number:        3
//    > block timestamp:     1630828221
//    > account:             0xd0446b3eD62f23815bE724C17e34C0617A186e34
//    > balance:             99.9487734
//    > gas used:            2269975 (0x22a317)
//    > gas price:           20 gwei
//    > value sent:          0 ETH
//    > total cost:          0.0453995 ETH

//    > Saving migration to chain.
//    > Saving artifacts
//    -------------------------------------
//    > Total cost:           0.0453995 ETH
