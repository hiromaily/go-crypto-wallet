// SPDX-License-Identifier: MIT
pragma solidity ^0.8.34;

import "forge-std/Test.sol";
import "../contracts/HYT.sol";

contract HYTTest is Test {
    HYT token;
    address deployer = address(this);
    address recipient = address(0xBEEF);

    uint256 constant INITIAL_SUPPLY = 1_000_000 ether;

    function setUp() public {
        token = new HYT(INITIAL_SUPPLY);
    }

    function test_InitialSupplyMintedToDeployer() public view {
        assertEq(token.balanceOf(deployer), INITIAL_SUPPLY);
    }

    function test_TotalSupplyEqualsInitialSupply() public view {
        assertEq(token.totalSupply(), INITIAL_SUPPLY);
    }

    function test_Name() public view {
        assertEq(token.name(), "hiromaily Token");
    }

    function test_Symbol() public view {
        assertEq(token.symbol(), "HYT");
    }

    function test_Decimals() public view {
        assertEq(token.decimals(), 18);
    }

    function test_Transfer() public {
        uint256 amount = 100 ether;
        bool ok = token.transfer(recipient, amount);
        assertTrue(ok);
        assertEq(token.balanceOf(recipient), amount);
        assertEq(token.balanceOf(deployer), INITIAL_SUPPLY - amount);
    }

    function test_TransferToZeroAddressReverts() public {
        vm.expectRevert();
        token.transfer(address(0), 1 ether);
    }
}
