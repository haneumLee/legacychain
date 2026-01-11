// SPDX-License-Identifier: AGPL-3.0
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../../src/VaultFactory.sol";
import "../../src/IndividualVault.sol";

/**
 * @title IndividualVaultTest
 * @notice Unit tests for IndividualVault contract
 */
contract IndividualVaultTest is Test {
    VaultFactory public factory;
    IndividualVault public vault;

    address public owner = address(1);
    address public heir1 = address(2);
    address public heir2 = address(3);
    address public heir3 = address(4);
    address public nonHeir = address(5);

    address[] public heirs;
    uint256[] public shares;

    uint256 constant HEARTBEAT_INTERVAL = 7 days;
    uint256 constant GRACE_PERIOD = 30 days;
    uint256 constant REQUIRED_APPROVALS = 2;

    event VaultCreated(address indexed vaultAddress, address indexed owner, uint256 vaultIndex, uint256 timestamp);
    event Heartbeat(uint256 timestamp, bytes32 commitment);
    event VaultUnlocked(uint256 unlockTime);
    event GracePeriodStarted(uint256 endTime);
    event UnlockCancelled(address indexed owner, uint256 timestamp);
    event InheritanceApproved(address indexed heir);
    event InheritanceClaimed(address indexed heir, uint256 amount);
    event EmergencyPaused(uint256 timestamp);
    event Deposited(address indexed from, uint256 amount);
    event Withdrawn(address indexed owner, uint256 amount);

    function setUp() public {
        // Deploy factory
        factory = new VaultFactory();

        // Setup heirs and shares
        heirs.push(heir1);
        heirs.push(heir2);
        heirs.push(heir3);

        shares.push(5000); // 50%
        shares.push(3000); // 30%
        shares.push(2000); // 20%

        // Fund owner
        vm.deal(owner, 100 ether);
    }

    function _createVault() internal returns (address) {
        vm.prank(owner);
        address vaultAddr = factory.createVault(
            heirs,
            shares,
            HEARTBEAT_INTERVAL,
            GRACE_PERIOD,
            REQUIRED_APPROVALS
        );
        vault = IndividualVault(payable(vaultAddr));
        return vaultAddr;
    }

    // ============ Factory Tests ============

    function test_FactoryCreatesVault() public {
        vm.expectEmit(false, true, false, false);
        emit VaultCreated(address(0), owner, 0, block.timestamp);

        vm.prank(owner);
        address vaultAddr = factory.createVault(
            heirs,
            shares,
            HEARTBEAT_INTERVAL,
            GRACE_PERIOD,
            REQUIRED_APPROVALS
        );

        assertTrue(vaultAddr != address(0), "Vault not created");
        assertEq(factory.totalVaults(), 1, "Total vaults incorrect");
        
        address[] memory ownerVaults = factory.getOwnerVaults(owner);
        assertEq(ownerVaults.length, 1, "Owner vaults count incorrect");
        assertEq(ownerVaults[0], vaultAddr, "Vault address mismatch");
    }

    function test_RevertWhen_NoHeirs() public {
        address[] memory emptyHeirs = new address[](0);
        uint256[] memory emptyShares = new uint256[](0);

        vm.prank(owner);
        vm.expectRevert("VaultFactory: no heirs");
        factory.createVault(
            emptyHeirs,
            emptyShares,
            HEARTBEAT_INTERVAL,
            GRACE_PERIOD,
            REQUIRED_APPROVALS
        );
    }

    function test_RevertWhen_SharesNotHundredPercent() public {
        shares[0] = 4000; // Total = 9000, not 10000

        vm.prank(owner);
        vm.expectRevert("VaultFactory: shares must be 100%");
        factory.createVault(
            heirs,
            shares,
            HEARTBEAT_INTERVAL,
            GRACE_PERIOD,
            REQUIRED_APPROVALS
        );
    }

    function test_RevertWhen_InvalidHeartbeatInterval() public {
        vm.prank(owner);
        vm.expectRevert("VaultFactory: invalid heartbeat interval");
        factory.createVault(
            heirs,
            shares,
            1 days, // Too short
            GRACE_PERIOD,
            REQUIRED_APPROVALS
        );
    }

    // ============ Initialization Tests ============

    function test_VaultInitialized() public {
        _createVault();

        IndividualVault.VaultConfig memory config = vault.getConfig();
        
        assertEq(config.owner, owner, "Owner incorrect");
        assertEq(config.heirs.length, 3, "Heirs count incorrect");
        assertEq(config.heirShares[0], 5000, "Heir1 share incorrect");
        assertEq(config.heartbeatInterval, HEARTBEAT_INTERVAL, "Heartbeat interval incorrect");
        assertEq(config.gracePeriod, GRACE_PERIOD, "Grace period incorrect");
        assertEq(config.requiredApprovals, REQUIRED_APPROVALS, "Required approvals incorrect");
        assertTrue(config.isLocked, "Vault should be locked");
        assertFalse(config.gracePeriodActive, "Grace period should not be active");
    }

    // ============ Commit-Reveal Heartbeat Tests ============

    function test_CommitRevealHeartbeat() public {
        _createVault();

        bytes32 nonce = keccak256("secret");
        bytes32 commitment = keccak256(abi.encodePacked(owner, nonce));

        // Commit phase
        vm.prank(owner);
        vault.commitHeartbeat(commitment);

        // Reveal phase
        vm.expectEmit(false, false, false, true);
        emit Heartbeat(block.timestamp, commitment);

        vm.prank(owner);
        vault.revealHeartbeat(nonce);

        IndividualVault.VaultConfig memory config = vault.getConfig();
        assertEq(config.lastHeartbeat, block.timestamp, "Heartbeat not updated");
    }

    function test_RevertWhen_CommitmentReused() public {
        _createVault();

        bytes32 nonce = keccak256("secret");
        bytes32 commitment = keccak256(abi.encodePacked(owner, nonce));

        vm.startPrank(owner);
        vault.commitHeartbeat(commitment);

        vm.expectRevert("IndividualVault: commitment already used");
        vault.commitHeartbeat(commitment);
        vm.stopPrank();
    }

    function test_RevertWhen_InvalidReveal() public {
        _createVault();

        bytes32 wrongNonce = keccak256("wrong");

        vm.prank(owner);
        vm.expectRevert("IndividualVault: invalid commitment");
        vault.revealHeartbeat(wrongNonce);
    }

    // ============ Unlock Tests ============

    function test_CheckAndUnlock() public {
        _createVault();

        // Fast forward past heartbeat interval
        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);

        vm.expectEmit(false, false, false, false);
        emit VaultUnlocked(block.timestamp + GRACE_PERIOD);

        vault.checkAndUnlock();

        IndividualVault.VaultConfig memory config = vault.getConfig();
        assertFalse(config.isLocked, "Vault should be unlocked");
        assertTrue(config.gracePeriodActive, "Grace period should be active");
        assertEq(config.unlockTime, block.timestamp + GRACE_PERIOD, "Unlock time incorrect");
    }

    function test_RevertWhen_HeartbeatNotExpired() public {
        _createVault();

        vm.expectRevert("IndividualVault: heartbeat not expired");
        vault.checkAndUnlock();
    }

    // ============ Grace Period Owner Return Tests ============

    function test_OwnerReturnsInGracePeriod() public {
        _createVault();

        // Unlock vault
        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        // Heirs approve
        vm.prank(heir1);
        vault.approveInheritance();
        
        IndividualVault.VaultConfig memory configBefore = vault.getConfig();
        assertEq(configBefore.approvalCount, 1, "Approval count should be 1");

        // Owner returns with heartbeat
        bytes32 nonce = keccak256("owner_returns");
        bytes32 commitment = keccak256(abi.encodePacked(owner, nonce));

        vm.startPrank(owner);
        vault.commitHeartbeat(commitment);
        
        vm.expectEmit(true, false, false, false);
        emit UnlockCancelled(owner, block.timestamp);
        
        vault.revealHeartbeat(nonce);
        vm.stopPrank();

        // Verify unlock cancelled
        IndividualVault.VaultConfig memory configAfter = vault.getConfig();
        assertTrue(configAfter.isLocked, "Vault should be locked again");
        assertFalse(configAfter.gracePeriodActive, "Grace period should be inactive");
        assertEq(configAfter.approvalCount, 0, "Approval count should be reset");
        assertFalse(vault.heirApprovals(heir1), "Heir approval should be reset");
    }

    // ============ Multi-sig Approval Tests ============

    function test_HeirApproval() public {
        _createVault();

        // Unlock vault
        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        vm.expectEmit(true, false, false, false);
        emit InheritanceApproved(heir1);

        vm.prank(heir1);
        vault.approveInheritance();

        assertTrue(vault.heirApprovals(heir1), "Heir1 approval not recorded");
        
        IndividualVault.VaultConfig memory config = vault.getConfig();
        assertEq(config.approvalCount, 1, "Approval count incorrect");
    }

    function test_RevertWhen_NotHeir() public {
        _createVault();

        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        vm.prank(nonHeir);
        vm.expectRevert("IndividualVault: not heir");
        vault.approveInheritance();
    }

    function test_RevertWhen_VaultLocked() public {
        _createVault();

        vm.prank(heir1);
        vm.expectRevert("IndividualVault: vault locked");
        vault.approveInheritance();
    }

    function test_RevertWhen_AlreadyApproved() public {
        _createVault();

        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        vm.startPrank(heir1);
        vault.approveInheritance();

        vm.expectRevert("IndividualVault: already approved");
        vault.approveInheritance();
        vm.stopPrank();
    }

    // ============ Claim Inheritance Tests ============

    function test_ClaimInheritance() public {
        _createVault();

        // Fund vault
        vm.deal(address(vault), 10 ether);

        // Unlock vault
        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        // Get approvals
        vm.prank(heir1);
        vault.approveInheritance();
        
        vm.prank(heir2);
        vault.approveInheritance();

        // Fast forward past grace period
        vm.warp(block.timestamp + GRACE_PERIOD + 1);

        // Heir1 claims (50%)
        uint256 balanceBefore = heir1.balance;
        
        vm.expectEmit(true, false, false, true);
        emit InheritanceClaimed(heir1, 5 ether);
        
        vm.prank(heir1);
        vault.claimInheritance();

        assertEq(heir1.balance - balanceBefore, 5 ether, "Heir1 claim amount incorrect");
        assertTrue(vault.heirClaimed(heir1), "Heir1 claimed flag not set");
    }

    function test_MultipleHeirsClaim() public {
        _createVault();

        vm.deal(address(vault), 10 ether);

        // Unlock and approve
        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        vm.prank(heir1);
        vault.approveInheritance();
        
        vm.prank(heir2);
        vault.approveInheritance();

        vm.warp(block.timestamp + GRACE_PERIOD + 1);

        // All heirs claim
        uint256 heir1Before = heir1.balance;
        vm.prank(heir1);
        vault.claimInheritance();
        assertEq(heir1.balance - heir1Before, 5 ether, "Heir1 incorrect"); // 50%

        uint256 heir2Before = heir2.balance;
        vm.prank(heir2);
        vault.claimInheritance();
        assertEq(heir2.balance - heir2Before, 3 ether, "Heir2 incorrect"); // 30%

        uint256 heir3Before = heir3.balance;
        vm.prank(heir3);
        vault.claimInheritance();
        assertEq(heir3.balance - heir3Before, 2 ether, "Heir3 incorrect"); // 20%
    }

    function test_RevertWhen_NotEnoughApprovals() public {
        _createVault();

        vm.deal(address(vault), 10 ether);

        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        // Only 1 approval (need 2)
        vm.prank(heir1);
        vault.approveInheritance();

        vm.warp(block.timestamp + GRACE_PERIOD + 1);

        vm.prank(heir1);
        vm.expectRevert("IndividualVault: not enough approvals");
        vault.claimInheritance();
    }

    function test_RevertWhen_GracePeriodNotEnded() public {
        _createVault();

        vm.deal(address(vault), 10 ether);

        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        vm.prank(heir1);
        vault.approveInheritance();
        
        vm.prank(heir2);
        vault.approveInheritance();

        // Don't wait for grace period to end

        vm.prank(heir1);
        vm.expectRevert("IndividualVault: grace period not ended");
        vault.claimInheritance();
    }

    function test_RevertWhen_AlreadyClaimed() public {
        _createVault();

        vm.deal(address(vault), 10 ether);

        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        vm.prank(heir1);
        vault.approveInheritance();
        
        vm.prank(heir2);
        vault.approveInheritance();

        vm.warp(block.timestamp + GRACE_PERIOD + 1);

        vm.startPrank(heir1);
        vault.claimInheritance();

        vm.expectRevert("IndividualVault: already claimed");
        vault.claimInheritance();
        vm.stopPrank();
    }

    // ============ Emergency Pause Tests ============

    function test_EmergencyPause() public {
        _createVault();

        vm.expectEmit(false, false, false, false);
        emit EmergencyPaused(block.timestamp);

        vm.prank(owner);
        vault.pause();

        assertTrue(vault.paused(), "Vault should be paused");
    }

    function test_PauseBlocksHeartbeat() public {
        _createVault();

        vm.prank(owner);
        vault.pause();

        bytes32 commitment = keccak256(abi.encodePacked(owner, keccak256("test")));

        vm.prank(owner);
        vm.expectRevert();
        vault.commitHeartbeat(commitment);
    }

    function test_PauseBlocksClaim() public {
        _createVault();

        vm.deal(address(vault), 10 ether);

        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        vm.prank(heir1);
        vault.approveInheritance();
        
        vm.prank(heir2);
        vault.approveInheritance();

        vm.warp(block.timestamp + GRACE_PERIOD + 1);

        // Pause before claim
        vm.prank(owner);
        vault.pause();

        vm.prank(heir1);
        vm.expectRevert();
        vault.claimInheritance();
    }

    function test_Unpause() public {
        _createVault();

        vm.startPrank(owner);
        vault.pause();
        vault.unpause();
        vm.stopPrank();

        assertFalse(vault.paused(), "Vault should be unpaused");
    }

    // ============ Owner Withdraw Tests ============

    function test_OwnerWithdraw() public {
        _createVault();

        vm.deal(address(vault), 10 ether);

        uint256 balanceBefore = owner.balance;

        vm.expectEmit(true, false, false, true);
        emit Withdrawn(owner, 5 ether);

        vm.prank(owner);
        vault.withdraw(5 ether);

        assertEq(owner.balance - balanceBefore, 5 ether, "Withdraw amount incorrect");
        assertEq(address(vault).balance, 5 ether, "Vault balance incorrect");
    }

    function test_RevertWhen_WithdrawUnlocked() public {
        _createVault();

        vm.deal(address(vault), 10 ether);

        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        vm.prank(owner);
        vm.expectRevert("IndividualVault: cannot withdraw when unlocked");
        vault.withdraw(5 ether);
    }

    // ============ Deposit Tests ============

    function test_Deposit() public {
        _createVault();

        vm.expectEmit(true, false, false, true);
        emit Deposited(address(this), 1 ether);

        (bool success, ) = address(vault).call{value: 1 ether}("");
        assertTrue(success, "Deposit failed");
        assertEq(address(vault).balance, 1 ether, "Balance incorrect");
    }

    // ============ View Function Tests ============

    function test_IsClaimable() public {
        _createVault();

        assertFalse(vault.isClaimable(), "Should not be claimable initially");

        vm.warp(block.timestamp + HEARTBEAT_INTERVAL + 1);
        vault.checkAndUnlock();

        vm.prank(heir1);
        vault.approveInheritance();
        
        vm.prank(heir2);
        vault.approveInheritance();

        assertFalse(vault.isClaimable(), "Should not be claimable in grace period");

        vm.warp(block.timestamp + GRACE_PERIOD + 1);

        assertTrue(vault.isClaimable(), "Should be claimable now");
    }

    function test_GetBalance() public {
        _createVault();

        vm.deal(address(vault), 5 ether);
        assertEq(vault.getBalance(), 5 ether, "Balance incorrect");
    }

    function test_IsHeir() public {
        _createVault();

        assertTrue(vault.isHeir(heir1), "Heir1 should be heir");
        assertTrue(vault.isHeir(heir2), "Heir2 should be heir");
        assertTrue(vault.isHeir(heir3), "Heir3 should be heir");
        assertFalse(vault.isHeir(nonHeir), "NonHeir should not be heir");
    }

    receive() external payable {}
}
