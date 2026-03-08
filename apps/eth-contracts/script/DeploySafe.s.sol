// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// Deployment script for Safe v1.4.1 proxy on Anvil (local) for E2E testing.
//
// Environment variables:
//   SAFE_OWNERS    Comma-separated list of owner addresses (e.g. "0xaaa...,0xbbb...")
//   SAFE_THRESHOLD Number of required confirmations (e.g. "2")
//   DEPLOYER_KEY   Private key of the broadcaster account (Anvil default account is fine)
//
// Output:
//   Prints "SAFE_ADDRESS=0x..." to stdout for shell script parsing.
//
// Usage:
//   forge script script/DeploySafe.s.sol \
//     --rpc-url http://localhost:8546 \
//     --broadcast \
//     -vvv

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {Safe} from "@safe-global/safe-contracts/contracts/Safe.sol";
import {SafeProxyFactory} from "@safe-global/safe-contracts/contracts/proxies/SafeProxyFactory.sol";

contract DeploySafe is Script {
    /// @notice Deploys a Safe v1.4.1 proxy with the given owners and threshold.
    /// @dev Reads SAFE_OWNERS (comma-separated addresses), SAFE_THRESHOLD,
    ///      and DEPLOYER_KEY from the environment.  Prints the proxy address
    ///      as "SAFE_ADDRESS=0x..." so callers can parse it with grep/sed.
    function run() external {
        uint256 deployerKey = vm.envUint("DEPLOYER_KEY");
        address[] memory owners = vm.envAddress("SAFE_OWNERS", ",");
        uint256 threshold = vm.envUint("SAFE_THRESHOLD");

        require(owners.length > 0, "DeploySafe: no owners provided");
        require(threshold > 0 && threshold <= owners.length, "DeploySafe: invalid threshold");

        vm.startBroadcast(deployerKey);

        // Deploy Safe singleton (master copy).
        Safe singleton = new Safe();

        // Deploy SafeProxyFactory.
        SafeProxyFactory factory = new SafeProxyFactory();

        // Encode Safe.setup() initialiser call.
        bytes memory initData = abi.encodeCall(
            Safe.setup,
            (
                owners,
                threshold,
                address(0), // optional delegate-call target (none)
                new bytes(0), // optional delegate-call data (none)
                address(0), // fallback handler (none for simplicity)
                address(0), // payment token (ETH, i.e. zero address)
                0, // payment amount
                payable(address(0)) // payment receiver
            )
        );

        // Create proxy.  Salt nonce = 0 produces a deterministic address for
        // the given singleton + initData combination.
        address safeProxy = address(factory.createProxyWithNonce(address(singleton), initData, 0));

        vm.stopBroadcast();

        // Print in "KEY=VALUE" format so the E2E shell script can parse with:
        //   SAFE_ADDR=$(forge script ... | grep 'SAFE_ADDRESS=' | cut -d= -f2)
        // vm.toString() is required — console.log(string, address) does NOT do
        // printf-style substitution; it would output "SAFE_ADDRESS=%s 0xAddr".
        console.log(string.concat("SAFE_ADDRESS=", vm.toString(safeProxy)));
    }
}
