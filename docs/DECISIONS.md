# LegacyChain - Architecture Decision Records (ADR)

> **목적**: 중요한 기술적 의사결정 기록 및 근거  
> **형식**: [Title] / Context / Decision / Consequences  
> **작성일**: 2026년 1월 12일

---

## 📋 목차

1. [ADR-001: Factory 패턴 선택](#adr-001-factory-패턴-선택)
2. [ADR-002: Commit-Reveal Heartbeat](#adr-002-commit-reveal-heartbeat)
3. [ADR-003: Pausable Emergency Stop](#adr-003-pausable-emergency-stop)
4. [ADR-004: OpenZeppelin v5.5.0 사용](#adr-004-openzeppelin-v550-사용)

---

## ADR-001: Factory 패턴 선택

### Date
2026-01-12

### Status
✅ Accepted

### Context
초기 설계에서는 단일 컨트랙트에 모든 Vault를 저장하는 방식을 고려했습니다:

```solidity
// ❌ 초기 설계
contract LegacyVault {
    mapping(uint256 => Vault) public vaults;  // 모든 Vault가 한 곳에
}
```

**문제점**:
1. **보안 격리 부족**: 한 Vault의 취약점이 모든 Vault에 영향
2. **가스비 비효율**: Storage slot 충돌, 비효율적인 메모리 사용
3. **업그레이드 어려움**: 전체 컨트랙트를 업그레이드해야 함
4. **복잡도 증가**: 단일 컨트랙트에 모든 로직 집중

### Decision
**Factory + Clone 패턴** 채택 (EIP-1167: Minimal Proxy Contract)

```solidity
// ✅ 개선된 설계
contract VaultFactory {
    address public immutable vaultImplementation;
    
    function createVault(...) external returns (address) {
        address vault = vaultImplementation.clone();  // EIP-1167
        IndividualVault(payable(vault)).initialize(...);
        return vault;
    }
}

contract IndividualVault is Initializable {
    // 각 Vault가 독립된 컨트랙트
}
```

**장점**:
1. **보안 격리**: 각 Vault가 독립된 컨트랙트 주소
2. **가스비 최적화**: 
   - Clone 배포 비용: ~45,000 gas
   - 일반 배포 대비 95% 절감
3. **유연한 업그레이드**: 개별 Vault만 영향 받음
4. **명확한 소유권**: 1 Address = 1 Vault Contract

### Consequences

**Positive**:
- ✅ Cross-vault 공격 차단
- ✅ 가스비 대폭 절감 (45k vs 800k)
- ✅ 개별 Vault Pausable/Upgradeable
- ✅ 확장성 향상

**Negative**:
- ⚠️ 컨트랙트 복잡도 증가 (Factory + Implementation)
- ⚠️ 초기 구현 시간 추가 소요
- ⚠️ Initialize 패턴 필수 (Constructor 사용 불가)

**Mitigation**:
- OpenZeppelin Clones.sol 사용으로 안전성 확보
- Initializable.sol로 재초기화 방지
- 충분한 테스트 커버리지 (>95%)

### References
- [EIP-1167: Minimal Proxy Contract](https://eips.ethereum.org/EIPS/eip-1167)
- [OpenZeppelin Clones](https://docs.openzeppelin.com/contracts/5.x/api/proxy#Clones)
- [Gas Comparison: Clone vs Create](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/proxy/Clones.sol#L7-L18)

---

## ADR-002: Commit-Reveal Heartbeat

### Date
2026-01-12

### Status
✅ Accepted

### Context
Heartbeat 트랜잭션이 Public Mempool에 노출되면 **Front-running 공격** 가능:

```solidity
// ❌ 취약한 설계
function heartbeat(uint256 _vaultId) external {
    // Mempool에서 보임 → Attacker가 먼저 approveInheritance() 호출 가능
    vaults[_vaultId].lastHeartbeat = block.timestamp;
}
```

**공격 시나리오**:
1. Owner가 Heartbeat 트랜잭션 전송
2. Attacker가 Mempool에서 감지
3. 더 높은 Gas Price로 `approveInheritance()` 먼저 실행
4. Owner의 Heartbeat보다 먼저 상속 진행

### Decision
**Commit-Reveal 패턴** 도입

```solidity
// ✅ 보안 강화 설계
mapping(bytes32 => bool) private usedCommitments;

function commitHeartbeat(bytes32 _commitment) external onlyOwner {
    require(!usedCommitments[_commitment], "Already used");
    usedCommitments[_commitment] = true;
}

function revealHeartbeat(bytes32 _nonce) external onlyOwner {
    bytes32 commitment = keccak256(abi.encodePacked(msg.sender, _nonce));
    require(usedCommitments[commitment], "Invalid");
    
    config.lastHeartbeat = block.timestamp;
    emit Heartbeat(block.timestamp, commitment);
}
```

**동작 방식**:
1. **Commit Phase**: `keccak256(owner, nonce)` 해시 제출
2. **Reveal Phase**: `nonce` 공개하여 검증 + Heartbeat 실행

### Consequences

**Positive**:
- ✅ Front-running 공격 완전 차단
- ✅ MEV (Maximal Extractable Value) 공격 방어
- ✅ Privacy 향상 (트랜잭션 의도 숨김)

**Negative**:
- ⚠️ 2개 트랜잭션 필요 (가스비 2배)
- ⚠️ UX 복잡도 증가
- ⚠️ Nonce 관리 필요

**Mitigation**:
- Frontend에서 자동 Commit-Reveal 처리
- Nonce는 timestamp + random 조합 사용
- 실패 시 재시도 로직 구현

**Alternative Considered**:
- ❌ Flashbots Private Transaction: 중앙화 우려
- ❌ Time-lock만 사용: Front-running 여전히 가능
- ✅ **Commit-Reveal**: 분산화 + 보안

### References
- [Commit-Reveal Pattern](https://github.com/ethereum/wiki/wiki/Safety#commit-reveal)
- [Front-running Attacks](https://consensys.github.io/smart-contract-best-practices/attacks/frontrunning/)

---

## ADR-003: Pausable Emergency Stop

### Date
2026-01-12

### Status
✅ Accepted

### Context
Smart Contract 배포 후 Critical 버그 발견 시 대응 방안 필요:
- 자산 손실 위험
- 악용 가능한 취약점
- 긴급 패치 필요

**기존 방식의 문제**:
- Immutable Contract → 버그 수정 불가
- 긴급 상황 대응 불가

### Decision
**OpenZeppelin Pausable** 도입

```solidity
import "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";

contract IndividualVault is Pausable, ReentrancyGuard {
    function pause() external onlyOwner {
        _pause();
        emit EmergencyPaused(block.timestamp);
    }
    
    function unpause() external onlyOwner {
        _unpause();
    }
    
    // Critical functions respect pause
    function claimInheritance() external whenNotPaused nonReentrant {
        // ...
    }
}
```

**Pausable 적용 함수**:
- ✅ `commitHeartbeat()` - Front-running 방어
- ✅ `revealHeartbeat()` - Heartbeat 실행
- ✅ `approveInheritance()` - 상속 승인
- ✅ `claimInheritance()` - 자산 인출
- ❌ `getBalance()` - View 함수는 제외

### Consequences

**Positive**:
- ✅ Circuit Breaker 역할 (버그 발견 시 즉시 중지)
- ✅ 자산 손실 방지
- ✅ 패치 배포 시간 확보
- ✅ Owner 권한으로 제어 가능

**Negative**:
- ⚠️ 중앙화 우려 (Owner가 악의적으로 pause 가능)
- ⚠️ 가스비 약간 증가 (whenNotPaused modifier)

**Mitigation**:
- Timelock + Multi-sig Owner 고려 (Phase 2)
- Pause 이유를 Event로 명확히 기록
- 정기적인 Security Audit
- Community Governance 도입 (장기)

### References
- [OpenZeppelin Pausable](https://docs.openzeppelin.com/contracts/5.x/api/security#Pausable)
- [Circuit Breaker Pattern](https://consensys.github.io/smart-contract-best-practices/development-recommendations/general/external-calls/#circuit-breakers)

---

## ADR-004: OpenZeppelin v5.5.0 사용

### Date
2026-01-12

### Status
✅ Accepted

### Context
Smart Contract 개발 시 라이브러리 선택 필요:
- 직접 구현 vs 검증된 라이브러리
- 보안 vs 커스터마이징

### Decision
**OpenZeppelin Contracts v5.5.0** 채택

설치된 라이브러리:
```bash
openzeppelin-contracts v5.5.0
openzeppelin-contracts-upgradeable v5.5.0
```

**사용 모듈**:
- `Clones.sol` - Factory 패턴 (EIP-1167)
- `Initializable.sol` - 초기화 패턴
- `PausableUpgradeable.sol` - Emergency Stop
- `ReentrancyGuardUpgradeable.sol` - Reentrancy 방어
- `Ownable.sol` - 소유권 관리 (간단한 경우)

### Consequences

**Positive**:
- ✅ Battle-tested 코드 (수백 개 프로젝트 사용)
- ✅ 정기적인 Security Audit
- ✅ 커뮤니티 지원 활발
- ✅ Gas Optimized
- ✅ EIP 표준 준수

**Negative**:
- ⚠️ 추가 의존성
- ⚠️ 라이브러리 크기 (50MB+)
- ⚠️ 업그레이드 시 호환성 체크 필요

**Mitigation**:
- 특정 버전 고정 (v5.5.0)
- Submodule로 관리
- 사용하지 않는 모듈은 import 제외

### Alternatives Considered
- ❌ Solmate: 가벼우나 Audit 부족
- ❌ 직접 구현: 시간 소요 + 보안 리스크
- ✅ **OpenZeppelin**: 안정성 + 검증됨

### References
- [OpenZeppelin Contracts](https://github.com/OpenZeppelin/openzeppelin-contracts)
- [OpenZeppelin Security](https://www.openzeppelin.com/security-audits)

---

## 추가 예정 ADR

- ADR-005: DID Registry 다중 Oracle (Phase 1.5)
- ADR-006: Emergency Recovery Guardian 구조
- ADR-007: ERC-4337 Account Abstraction (Phase 2)
- ADR-008: Gas Optimization 전략
- ADR-009: Layer 2 Migration 계획

---

**Last Updated**: 2026-01-12  
**Status**: Active Development
