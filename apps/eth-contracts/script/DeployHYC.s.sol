// SPDX-License-Identifier: MIT
pragma solidity ^0.8.34;

import "forge-std/Script.sol";
import "../contracts/HYC.sol";

contract DeployHYC is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployerKey);

        new HYC(1_000_000 ether);

        vm.stopBroadcast();
    }
}
