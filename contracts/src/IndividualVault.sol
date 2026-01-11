// SPDX-License-Identifier: AGPL-3.0
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

/**
 * @title IndividualVault
 * @notice Individual vault contract for digital asset inheritance
 * @dev Each vault is a separate contract instance created via VaultFactory
 * 
 * Security Features:
 * - Commit-Reveal Heartbeat (ADR-002): Prevents front-running attacks
 * - Pausable (ADR-003): Emergency stop mechanism
 * - ReentrancyGuard: Prevents reentrancy attacks
 * - Grace Period: Allows owner to cancel unlock if still alive
 * 
 * Architecture: Factory Pattern (ADR-001)
 */
contract IndividualVault is
    Initializable,
    PausableUpgradeable,
    ReentrancyGuard
{
    /// @notice Vault configuration struct
    struct VaultConfig {
        address owner;                  // Vault owner
        address[] heirs;                // List of heirs
        uint256[] heirShares;           // Heir shares in basis points (10000 = 100%)
        uint256 heartbeatInterval;      // Required heartbeat frequency (seconds)
        uint256 lastHeartbeat;          // Timestamp of last heartbeat
        uint256 unlockTime;             // When vault can be claimed
        uint256 gracePeriod;            // Grace period duration (seconds)
        uint256 requiredApprovals;      // Number of approvals needed
        uint256 approvalCount;          // Current number of approvals
        uint256 totalBalanceAtUnlock;   // Balance snapshot when grace period ends
        bool isLocked;                  // Vault lock status
        bool gracePeriodActive;         // Whether in grace period (owner can return)
    }

    /// @notice Vault configuration
    VaultConfig public config;

    /// @notice Mapping of heir approvals
    mapping(address => bool) public heirApprovals;

    /// @notice Mapping of used commitments (for Commit-Reveal pattern)
    /// @dev Prevents replay attacks and front-running
    mapping(bytes32 => bool) private usedCommitments;

    /// @notice Mapping to track if heir has claimed
    mapping(address => bool) public heirClaimed;

    // ============ Events ============

    event Heartbeat(uint256 timestamp, bytes32 commitment);
    event VaultUnlocked(uint256 unlockTime);
    event GracePeriodStarted(uint256 endTime);
    event UnlockCancelled(address indexed owner, uint256 timestamp);
    event InheritanceApproved(address indexed heir);
    event InheritanceClaimed(address indexed heir, uint256 amount);
    event EmergencyPaused(uint256 timestamp);
    event Deposited(address indexed from, uint256 amount);
    event Withdrawn(address indexed owner, uint256 amount);

    // ============ Modifiers ============

    modifier onlyOwner() {
        require(msg.sender == config.owner, "IndividualVault: not owner");
        _;
    }

    modifier onlyHeir() {
        bool isHeir = false;
        for (uint256 i = 0; i < config.heirs.length; i++) {
            if (config.heirs[i] == msg.sender) {
                isHeir = true;
                break;
            }
        }
        require(isHeir, "IndividualVault: not heir");
        _;
    }

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initializes the vault (replaces constructor for clones)
     * @param _owner Vault owner address
     * @param _heirs Array of heir addresses
     * @param _shares Array of heir shares in basis points
     * @param _heartbeatInterval Heartbeat frequency in seconds
     * @param _gracePeriod Grace period duration in seconds
     * @param _requiredApprovals Number of approvals required
     * 
     * @dev Called by VaultFactory after cloning
     * Requirements enforced by VaultFactory:
     * - Heirs length > 0 and <= 10
     * - Shares sum = 10000
     * - Valid time parameters
     */
    function initialize(
        address _owner,
        address[] memory _heirs,
        uint256[] memory _shares,
        uint256 _heartbeatInterval,
        uint256 _gracePeriod,
        uint256 _requiredApprovals
    ) external initializer {
        __Pausable_init();
        // ReentrancyGuard doesn't need initialization in OpenZeppelin v5

        config.owner = _owner;
        config.heirs = _heirs;
        config.heirShares = _shares;
        config.heartbeatInterval = _heartbeatInterval;
        config.gracePeriod = _gracePeriod;
        config.requiredApprovals = _requiredApprovals;
        config.lastHeartbeat = block.timestamp;
        config.isLocked = true;
        config.gracePeriodActive = false;
    }

    // ============ Heartbeat Functions (Commit-Reveal Pattern) ============

    /**
     * @notice Commit phase of heartbeat (Step 1/2)
     * @param _commitment Hash of (owner address + nonce)
     * 
     * @dev Commit-Reveal pattern prevents front-running (ADR-002)
     * Generate commitment: keccak256(abi.encodePacked(msg.sender, nonce))
     */
    function commitHeartbeat(bytes32 _commitment)
        external
        onlyOwner
        whenNotPaused
    {
        require(!usedCommitments[_commitment], "IndividualVault: commitment already used");
        usedCommitments[_commitment] = true;
    }

    /**
     * @notice Reveal phase of heartbeat (Step 2/2)
     * @param _nonce Random nonce used in commitment
     * 
     * @dev Verifies commitment and executes heartbeat
     * If in grace period, cancels unlock and resets approvals
     */
    function revealHeartbeat(bytes32 _nonce)
        external
        onlyOwner
        whenNotPaused
    {
        bytes32 commitment = keccak256(abi.encodePacked(msg.sender, _nonce));
        require(usedCommitments[commitment], "IndividualVault: invalid commitment");

        // Update heartbeat
        config.lastHeartbeat = block.timestamp;
        config.isLocked = true;

        // If in grace period, owner has returned - cancel unlock
        if (config.gracePeriodActive) {
            config.gracePeriodActive = false;
            config.approvalCount = 0;
            config.totalBalanceAtUnlock = 0; // Reset balance snapshot

            // Reset all heir approvals
            for (uint256 i = 0; i < config.heirs.length; i++) {
                heirApprovals[config.heirs[i]] = false;
            }

            emit UnlockCancelled(msg.sender, block.timestamp);
        }

        emit Heartbeat(block.timestamp, commitment);
    }

    // ============ Inheritance Functions ============

    /**
     * @notice Check if vault should be unlocked and start grace period
     * @dev Anyone can call this to trigger unlock check
     */
    function checkAndUnlock() external {
        require(
            block.timestamp >= config.lastHeartbeat + config.heartbeatInterval,
            "IndividualVault: heartbeat not expired"
        );
        require(config.isLocked, "IndividualVault: already unlocked");

        config.isLocked = false;
        config.unlockTime = block.timestamp + config.gracePeriod;
        config.gracePeriodActive = true;

        emit VaultUnlocked(config.unlockTime);
        emit GracePeriodStarted(config.unlockTime);
    }

    /**
     * @notice Heir approves inheritance
     * @dev Part of multi-sig approval process
     */
    function approveInheritance()
        external
        onlyHeir
        whenNotPaused
    {
        require(!config.isLocked, "IndividualVault: vault locked");
        require(!heirApprovals[msg.sender], "IndividualVault: already approved");

        heirApprovals[msg.sender] = true;
        config.approvalCount++;

        emit InheritanceApproved(msg.sender);
    }

    /**
     * @notice Heir claims their inheritance share
     * @dev Uses CEI pattern and ReentrancyGuard for security
     * 
     * Requirements:
     * - Vault must be unlocked
     * - Required approvals must be met
     * - Grace period must have ended
     * - Caller must be an heir
     * - Heir must not have claimed yet
     */
    function claimInheritance()
        external
        onlyHeir
        nonReentrant
        whenNotPaused
    {
        // Checks
        require(!config.isLocked, "IndividualVault: vault locked");
        require(
            config.approvalCount >= config.requiredApprovals,
            "IndividualVault: not enough approvals"
        );
        require(
            block.timestamp >= config.unlockTime,
            "IndividualVault: grace period not ended"
        );
        require(!heirClaimed[msg.sender], "IndividualVault: already claimed");

        // Find heir index and calculate share
        uint256 heirIndex = type(uint256).max;
        for (uint256 i = 0; i < config.heirs.length; i++) {
            if (config.heirs[i] == msg.sender) {
                heirIndex = i;
                break;
            }
        }
        require(heirIndex != type(uint256).max, "IndividualVault: not heir");

        // Effects
        heirClaimed[msg.sender] = true;
        
        // Snapshot balance on first claim after grace period ends
        if (config.totalBalanceAtUnlock == 0) {
            config.totalBalanceAtUnlock = address(this).balance;
        }
        
        // Use snapshotted balance for fair distribution
        uint256 amount = (config.totalBalanceAtUnlock * config.heirShares[heirIndex]) / 10000;

        // Interactions (CEI pattern)
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "IndividualVault: transfer failed");

        emit InheritanceClaimed(msg.sender, amount);
    }

    // ============ Emergency Functions ============

    /**
     * @notice Pause the vault in case of emergency
     * @dev Only owner can pause. Stops all critical functions.
     */
    function pause() external onlyOwner {
        _pause();
        emit EmergencyPaused(block.timestamp);
    }

    /**
     * @notice Unpause the vault
     * @dev Only owner can unpause
     */
    function unpause() external onlyOwner {
        _unpause();
    }

    // ============ Owner Functions ============

    /**
     * @notice Owner withdraws funds (only when locked)
     * @param _amount Amount to withdraw
     * 
     * @dev Owner can only withdraw when vault is locked (not in inheritance mode)
     */
    function withdraw(uint256 _amount)
        external
        onlyOwner
        whenNotPaused
    {
        require(config.isLocked, "IndividualVault: cannot withdraw when unlocked");
        require(address(this).balance >= _amount, "IndividualVault: insufficient balance");

        (bool success, ) = msg.sender.call{value: _amount}("");
        require(success, "IndividualVault: transfer failed");

        emit Withdrawn(msg.sender, _amount);
    }

    /**
     * @notice Deposit ETH to vault
     * @dev Anyone can deposit to help the vault
     */
    receive() external payable {
        emit Deposited(msg.sender, msg.value);
    }

    // ============ View Functions ============

    /**
     * @notice Check if vault is currently claimable
     * @return true if all conditions are met for claiming
     */
    function isClaimable() external view returns (bool) {
        return
            !config.isLocked &&
            config.approvalCount >= config.requiredApprovals &&
            block.timestamp >= config.unlockTime;
    }

    /**
     * @notice Get vault balance
     * @return Current balance in wei
     */
    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }

    /**
     * @notice Get vault configuration
     * @return VaultConfig struct
     */
    function getConfig() external view returns (VaultConfig memory) {
        return config;
    }

    /**
     * @notice Get heir at index
     * @param _index Heir index
     * @return Heir address
     */
    function getHeir(uint256 _index) external view returns (address) {
        require(_index < config.heirs.length, "IndividualVault: index out of bounds");
        return config.heirs[_index];
    }

    /**
     * @notice Get heir share at index
     * @param _index Heir index
     * @return Share in basis points
     */
    function getHeirShare(uint256 _index) external view returns (uint256) {
        require(_index < config.heirs.length, "IndividualVault: index out of bounds");
        return config.heirShares[_index];
    }

    /**
     * @notice Get number of heirs
     * @return Number of heirs
     */
    function getHeirCount() external view returns (uint256) {
        return config.heirs.length;
    }

    /**
     * @notice Check if address is an heir
     * @param _address Address to check
     * @return true if address is an heir
     */
    function isHeir(address _address) external view returns (bool) {
        for (uint256 i = 0; i < config.heirs.length; i++) {
            if (config.heirs[i] == _address) {
                return true;
            }
        }
        return false;
    }
}
