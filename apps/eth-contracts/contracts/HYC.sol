// SPDX-License-Identifier: MIT
pragma solidity ^0.8.34;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

/// @title HYC — hiromaily Coin ERC-20 token
/// @author hiromaily
/// @notice Standard ERC-20 token; mints the entire initial supply to the deployer on construction.
contract HYC is ERC20 {
    /// @notice Deploys the HYC token and mints the full supply to the caller.
    /// @param initialSupply Total token supply minted to the deployer (in wei units).
    constructor(uint256 initialSupply) ERC20("hiromaily Coin", "HYC") {
        _mint(msg.sender, initialSupply);
    }
}
