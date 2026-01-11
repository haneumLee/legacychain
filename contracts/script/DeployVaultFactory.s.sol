// SPDX-License-Identifier: AGPL-3.0
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/VaultFactory.sol";

/**
 * @title DeployVaultFactory
 * @notice Deployment script for VaultFactory contract
 * @dev Usage:
 *   Local: forge script script/DeployVaultFactory.s.sol --fork-url http://localhost:8545 --broadcast
 *   Testnet: forge script script/DeployVaultFactory.s.sol --rpc-url $RPC_URL --broadcast --verify
 */
contract DeployVaultFactory is Script {
    function run() external returns (VaultFactory) {
        // Get deployer private key from environment
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        // Start broadcasting transactions
        vm.startBroadcast(deployerPrivateKey);

        // Deploy VaultFactory
        VaultFactory factory = new VaultFactory();

        console.log("VaultFactory deployed at:", address(factory));
        console.log("Deployer:", vm.addr(deployerPrivateKey));
        console.log("Implementation vault:", factory.vaultImplementation());

        vm.stopBroadcast();

        return factory;
    }
}
