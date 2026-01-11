// SPDX-License-Identifier: AGPL-3.0
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "forge-std/StdInvariant.sol";
import "../../src/VaultFactory.sol";
import "../../src/IndividualVault.sol";

/**
 * @title VaultInvariantTest
 * @notice Invariant tests for VaultFactory and IndividualVault
 * @dev Tests critical properties that must always hold
 */
contract VaultInvariantTest is StdInvariant, Test {
    VaultFactory public factory;
    VaultHandler public handler;

    function setUp() public {
        // Deploy factory
        factory = new VaultFactory();

        // Deploy handler
        handler = new VaultHandler(factory);

        // Target handler for invariant testing
        targetContract(address(handler));

        // Only call specific handler functions
        bytes4[] memory selectors = new bytes4[](5);
        selectors[0] = VaultHandler.createVault.selector;
        selectors[1] = VaultHandler.depositToVault.selector;
        selectors[2] = VaultHandler.commitHeartbeat.selector;
        selectors[3] = VaultHandler.unlockVault.selector;
        selectors[4] = VaultHandler.approveInheritance.selector;

        targetSelector(
            FuzzSelector({addr: address(handler), selectors: selectors})
        );
    }

    /**
     * @notice Invariant: Total heir shares must always equal 100% (10000 basis points)
     */
    function invariant_HeirSharesAlwaysHundredPercent() public view {
        IndividualVault[] memory vaults = handler.getCreatedVaults();

        for (uint256 i = 0; i < vaults.length; i++) {
            IndividualVault.VaultConfig memory config = vaults[i].getConfig();
            
            uint256 totalShares = 0;
            for (uint256 j = 0; j < config.heirShares.length; j++) {
                totalShares += config.heirShares[j];
            }

            assertEq(
                totalShares,
                10000,
                "Total heir shares must equal 100%"
            );
        }
    }

    /**
     * @notice Invariant: Total claimed amount never exceeds total balance snapshot
     */
    function invariant_TotalClaimedNeverExceedsBalance() public view {
        IndividualVault[] memory vaults = handler.getCreatedVaults();

        for (uint256 i = 0; i < vaults.length; i++) {
            IndividualVault.VaultConfig memory config = vaults[i].getConfig();
            
            if (config.totalBalanceAtUnlock > 0) {
                uint256 totalClaimed = handler.getTotalClaimed(address(vaults[i]));
                
                assertTrue(
                    totalClaimed <= config.totalBalanceAtUnlock,
                    "Total claimed exceeds balance snapshot"
                );
            }
        }
    }

    /**
     * @notice Invariant: Locked vault cannot have approvals
     */
    function invariant_LockedVaultHasNoApprovals() public view {
        IndividualVault[] memory vaults = handler.getCreatedVaults();

        for (uint256 i = 0; i < vaults.length; i++) {
            IndividualVault.VaultConfig memory config = vaults[i].getConfig();
            
            if (config.isLocked) {
                assertEq(
                    config.approvalCount,
                    0,
                    "Locked vault should have no approvals"
                );
            }
        }
    }

    /**
     * @notice Invariant: Grace period active only when unlocked
     */
    function invariant_GracePeriodOnlyWhenUnlocked() public view {
        IndividualVault[] memory vaults = handler.getCreatedVaults();

        for (uint256 i = 0; i < vaults.length; i++) {
            IndividualVault.VaultConfig memory config = vaults[i].getConfig();
            
            if (config.gracePeriodActive) {
                assertFalse(
                    config.isLocked,
                    "Grace period only active when unlocked"
                );
            }
        }
    }

    /**
     * @notice Invariant: Unlock time always in the future when grace period active
     */
    function invariant_UnlockTimeInFuture() public view {
        IndividualVault[] memory vaults = handler.getCreatedVaults();

        for (uint256 i = 0; i < vaults.length; i++) {
            IndividualVault.VaultConfig memory config = vaults[i].getConfig();
            
            if (config.gracePeriodActive && config.unlockTime > 0) {
                assertTrue(
                    config.unlockTime > config.lastHeartbeat,
                    "Unlock time must be after last heartbeat"
                );
            }
        }
    }
}

/**
 * @title VaultHandler
 * @notice Handler contract for invariant testing
 * @dev Provides bounded random operations on vaults
 */
contract VaultHandler is Test {
    VaultFactory public factory;
    IndividualVault[] public createdVaults;
    
    mapping(address => uint256) public totalClaimed;

    address[] public actors;
    mapping(address => bytes32) public nonces;

    constructor(VaultFactory _factory) {
        factory = _factory;

        // Create actors
        actors.push(address(0x1));
        actors.push(address(0x2));
        actors.push(address(0x3));
        actors.push(address(0x4));
        actors.push(address(0x5));

        // Fund actors
        for (uint256 i = 0; i < actors.length; i++) {
            vm.deal(actors[i], 100 ether);
        }
    }

    function createVault(uint256 ownerSeed, uint256 heirCount) public {
        // Bound inputs
        ownerSeed = bound(ownerSeed, 0, actors.length - 1);
        heirCount = bound(heirCount, 2, 4); // 2-4 heirs

        address owner = actors[ownerSeed];

        // Create heirs (different from owner)
        address[] memory heirs = new address[](heirCount);
        uint256[] memory shares = new uint256[](heirCount);
        
        uint256 remainingShares = 10000;
        for (uint256 i = 0; i < heirCount; i++) {
            uint256 heirIndex = (ownerSeed + i + 1) % actors.length;
            heirs[i] = actors[heirIndex];
            
            if (i == heirCount - 1) {
                // Last heir gets remaining shares
                shares[i] = remainingShares;
            } else {
                // Random share between 1000-4000 (10%-40%)
                uint256 share = 1000 + (uint256(keccak256(abi.encodePacked(block.timestamp, i))) % 3000);
                share = share > remainingShares ? remainingShares : share;
                shares[i] = share;
                remainingShares -= share;
            }
        }

        vm.prank(owner);
        try factory.createVault(
            heirs,
            shares,
            7 days,
            30 days,
            (heirCount + 1) / 2 // Majority approval
        ) returns (address vaultAddr) {
            createdVaults.push(IndividualVault(payable(vaultAddr)));
        } catch {
            // Invalid configuration, skip
        }
    }

    function depositToVault(uint256 vaultIndex, uint256 amount) public {
        if (createdVaults.length == 0) return;

        vaultIndex = bound(vaultIndex, 0, createdVaults.length - 1);
        amount = bound(amount, 0.01 ether, 10 ether);

        IndividualVault vault = createdVaults[vaultIndex];
        
        (bool success, ) = address(vault).call{value: amount}("");
        require(success);
    }

    function commitHeartbeat(uint256 vaultIndex, uint256 actorSeed) public {
        if (createdVaults.length == 0) return;

        vaultIndex = bound(vaultIndex, 0, createdVaults.length - 1);
        actorSeed = bound(actorSeed, 0, actors.length - 1);

        IndividualVault vault = createdVaults[vaultIndex];
        IndividualVault.VaultConfig memory config = vault.getConfig();
        
        address actor = actors[actorSeed];

        if (actor == config.owner && !vault.paused()) {
            bytes32 nonce = keccak256(abi.encodePacked(block.timestamp, actorSeed));
            nonces[address(vault)] = nonce;
            bytes32 commitment = keccak256(abi.encodePacked(actor, nonce));

            vm.prank(actor);
            try vault.commitHeartbeat(commitment) {
                // Reveal immediately
                vm.prank(actor);
                try vault.revealHeartbeat(nonce) {} catch {}
            } catch {}
        }
    }

    function unlockVault(uint256 vaultIndex) public {
        if (createdVaults.length == 0) return;

        vaultIndex = bound(vaultIndex, 0, createdVaults.length - 1);

        IndividualVault vault = createdVaults[vaultIndex];
        IndividualVault.VaultConfig memory config = vault.getConfig();

        if (config.isLocked && block.timestamp > config.lastHeartbeat + config.heartbeatInterval) {
            try vault.checkAndUnlock() {} catch {}
        }
    }

    function approveInheritance(uint256 vaultIndex, uint256 actorSeed) public {
        if (createdVaults.length == 0) return;

        vaultIndex = bound(vaultIndex, 0, createdVaults.length - 1);
        actorSeed = bound(actorSeed, 0, actors.length - 1);

        IndividualVault vault = createdVaults[vaultIndex];
        IndividualVault.VaultConfig memory config = vault.getConfig();
        
        address actor = actors[actorSeed];

        if (!config.isLocked && vault.isHeir(actor) && !vault.heirApprovals(actor)) {
            vm.prank(actor);
            try vault.approveInheritance() {} catch {}
        }
    }

    function getCreatedVaults() external view returns (IndividualVault[] memory) {
        return createdVaults;
    }

    function getTotalClaimed(address vault) external view returns (uint256) {
        return totalClaimed[vault];
    }

    receive() external payable {}
}
