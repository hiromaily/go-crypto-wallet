// SPDX-License-Identifier: MIT
pragma solidity ^0.8.34;

import "forge-std/Script.sol";
import "../contracts/HYT.sol";

contract DeployHYT is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployerKey);

        new HYT(1_000_000 ether);

        vm.stopBroadcast();
    }
}
