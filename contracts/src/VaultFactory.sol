// SPDX-License-Identifier: AGPL-3.0
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/proxy/Clones.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "./IndividualVault.sol";

/**
 * @title VaultFactory
 * @notice Factory contract for creating individual vault instances using EIP-1167 Clone pattern
 * @dev Uses minimal proxy pattern to reduce gas costs (~45k vs ~800k for regular deployment)
 * 
 * Architecture Decision: ADR-001 in docs/DECISIONS.md
 * - Each vault is an independent contract for security isolation
 * - Clone pattern for 95% gas savings
 * - Individual vault upgradability
 */
contract VaultFactory is Ownable {
    using Clones for address;

    /// @notice Implementation contract for cloning
    address public immutable vaultImplementation;

    /// @notice Mapping from owner address to their vaults
    mapping(address => address[]) public ownerVaults;

    /// @notice Total number of vaults created
    uint256 public totalVaults;

    /// @notice Emitted when a new vault is created
    event VaultCreated(
        address indexed vaultAddress,
        address indexed owner,
        uint256 vaultIndex,
        uint256 timestamp
    );

    /**
     * @notice Constructor deploys the implementation contract
     * @dev Implementation is deployed once and cloned for each vault
     */
    constructor() Ownable(msg.sender) {
        vaultImplementation = address(new IndividualVault());
    }

    /**
     * @notice Creates a new vault for the caller
     * @param _heirs Array of heir addresses
     * @param _shares Array of heir shares in basis points (10000 = 100%)
     * @param _heartbeatInterval Time in seconds between required heartbeats
     * @param _gracePeriod Time in seconds for grace period after unlock
     * @param _requiredApprovals Number of heir approvals required
     * @return vaultAddress Address of the newly created vault
     * 
     * @dev Uses EIP-1167 clone to minimize gas costs
     * Requirements:
     * - Heirs array must not be empty
     * - Heirs and shares arrays must have the same length
     * - Shares must sum to exactly 10000 (100%)
     * - Heartbeat interval between 7-90 days
     * - Grace period between 30-365 days
     */
    function createVault(
        address[] memory _heirs,
        uint256[] memory _shares,
        uint256 _heartbeatInterval,
        uint256 _gracePeriod,
        uint256 _requiredApprovals
    ) external returns (address vaultAddress) {
        // Input validation
        require(_heirs.length > 0, "VaultFactory: no heirs");
        require(_heirs.length == _shares.length, "VaultFactory: length mismatch");
        require(_heirs.length <= 10, "VaultFactory: too many heirs");
        
        // Validate shares sum to 100%
        uint256 totalShare = 0;
        for (uint256 i = 0; i < _shares.length; i++) {
            require(_shares[i] > 0, "VaultFactory: zero share");
            totalShare += _shares[i];
        }
        require(totalShare == 10000, "VaultFactory: shares must be 100%");

        // Validate time parameters
        require(
            _heartbeatInterval >= 7 days && _heartbeatInterval <= 90 days,
            "VaultFactory: invalid heartbeat interval"
        );
        require(
            _gracePeriod >= 30 days && _gracePeriod <= 365 days,
            "VaultFactory: invalid grace period"
        );
        require(
            _requiredApprovals > 0 && _requiredApprovals <= _heirs.length,
            "VaultFactory: invalid required approvals"
        );

        // Clone the implementation (EIP-1167)
        vaultAddress = vaultImplementation.clone();

        // Initialize the vault
        IndividualVault(payable(vaultAddress)).initialize(
            msg.sender,
            _heirs,
            _shares,
            _heartbeatInterval,
            _gracePeriod,
            _requiredApprovals
        );

        // Record ownership
        ownerVaults[msg.sender].push(vaultAddress);
        totalVaults++;

        emit VaultCreated(
            vaultAddress,
            msg.sender,
            ownerVaults[msg.sender].length - 1,
            block.timestamp
        );

        return vaultAddress;
    }

    /**
     * @notice Get all vaults owned by an address
     * @param _owner Address to query
     * @return Array of vault addresses
     */
    function getOwnerVaults(address _owner)
        external
        view
        returns (address[] memory)
    {
        return ownerVaults[_owner];
    }

    /**
     * @notice Get the number of vaults owned by an address
     * @param _owner Address to query
     * @return Number of vaults
     */
    function getOwnerVaultCount(address _owner)
        external
        view
        returns (uint256)
    {
        return ownerVaults[_owner].length;
    }

    /**
     * @notice Get a specific vault by owner and index
     * @param _owner Address of the owner
     * @param _index Index in the owner's vault array
     * @return Address of the vault
     */
    function getOwnerVaultAt(address _owner, uint256 _index)
        external
        view
        returns (address)
    {
        require(_index < ownerVaults[_owner].length, "VaultFactory: index out of bounds");
        return ownerVaults[_owner][_index];
    }
}
