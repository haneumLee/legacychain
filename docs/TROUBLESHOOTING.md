# LegacyChain - Troubleshooting Guide

> **목적**: 에러 해결 및 개선 히스토리 기록  
> **작성일**: 2026년 1월 12일

---

## 📋 목차

1. [PRD 문서 개선 히스토리](#1-prd-문서-개선-히스토리)
2. [Smart Contract 이슈](#2-smart-contract-이슈)
3. [Backend API 이슈](#3-backend-api-이슈)
4. [Frontend 이슈](#4-frontend-이슈)
5. [DevOps 이슈](#5-devops-이슈)

---

## 1. PRD 문서 개선 히스토리

### [2026-01-12] PRD 보안 강화 업데이트

#### Date
2026-01-12

#### Changes
PRD 문서(v1.0 → v1.1) 보안 및 설계 개선

#### Root Cause Analysis
초기 PRD 문서가 다음 보안 및 설계 이슈를 간과함:
1. **Smart Contract 설계**: 단일 컨트랙트에 모든 Vault 저장 → 보안 격리 부족, 가스비 비효율
2. **Front-running 공격**: Heartbeat 트랜잭션이 Mempool에서 노출 → Attacker가 먼저 상속 승인 가능
3. **Grace Period Owner 복귀**: Owner가 Grace Period 중 돌아올 경우 처리 로직 누락
4. **Database 스키마**: vault_id가 INTEGER (2^31-1 제한), Soft Delete 미지원
5. **Emergency Stop**: Critical 버그 발견 시 긴급 중지 메커니즘 없음
6. **Oracle 의존성**: NICE API 단일 의존 → 중단 시 전체 인증 불가
7. **법적 리스크**: Smart Contract의 법적 효력 불확실성 미대응
8. **가스비 리스크**: ETH 가스비 폭등 시 사용자 이탈 리스크

#### Solution Implemented

##### 1. Smart Contract Factory 패턴 도입
```solidity
// Before: 단일 컨트랙트
contract LegacyVault {
    mapping(uint256 => Vault) public vaults;  // 모든 Vault가 한 컨트랙트에
}

// After: Factory + Clone 패턴
contract VaultFactory {
    function createVault(...) returns (address) {
        return vaultImplementation.clone();  // EIP-1167
    }
}

contract IndividualVault {
    // 각 Vault가 독립된 컨트랙트
    // ✅ 보안 격리
    // ✅ 가스비 최적화
    // ✅ Upgrade 유연성
}
```

##### 2. Front-running 방지 (Commit-Reveal)
```solidity
// Commit Phase
function commitHeartbeat(bytes32 _commitment) external {
    usedCommitments[_commitment] = true;
}

// Reveal Phase
function revealHeartbeat(bytes32 _nonce) external {
    bytes32 commitment = keccak256(abi.encodePacked(msg.sender, _nonce));
    require(usedCommitments[commitment], "Invalid");
    // Execute heartbeat
}
```

##### 3. Grace Period Owner 복귀 처리
```solidity
function revealHeartbeat(bytes32 _nonce) external {
    // ...
    if (config.gracePeriodActive) {
        config.gracePeriodActive = false;
        config.approvalCount = 0;
        // 모든 Heir 승인 초기화
        emit UnlockCancelled(msg.sender, block.timestamp);
    }
}
```

##### 4. Database 스키마 개선
```sql
-- BIGINT로 확장
vault_id BIGINT UNIQUE NOT NULL,

-- Soft Delete 지원
deleted_at TIMESTAMP,

-- 인덱스에 Soft Delete 조건 추가
CREATE INDEX idx_vaults_owner ON vaults(owner_id) WHERE deleted_at IS NULL;
```

##### 5. Emergency Stop (Pausable)
```solidity
import "@openzeppelin/contracts/security/Pausable.sol";

contract IndividualVault is Pausable {
    function pause() external onlyOwner {
        _pause();
    }
    
    function claimInheritance() external whenNotPaused {
        // Critical functions respect pause
    }
}
```

##### 6. 다중 Oracle 지원
```solidity
enum VerificationProvider {
    NICE, PASS, OIDC, CHAINLINK
}

struct Attestation {
    VerificationProvider provider;
    bytes32 identityHash;
    address attestor;
}

// 최소 2개 이상 Attestation 필요
function addAttestation(...) external onlyRole(ORACLE_ROLE) {
    if (doc.attestations.length >= 2) {
        doc.verified = true;
    }
}
```

##### 7. 법적 리스크 대응
- 서비스 약관에 명확한 면책 조항 추가
- "법적 유언장의 기술적 보조 도구" 포지셔닝
- 법무법인 협업 및 공증 서비스 연동 계획
- 법원 명령 기반 Emergency Recovery 메커니즘

##### 8. 가스비 최적화 전략
- Layer 2 마이그레이션 로드맵 (Arbitrum/Optimism)
- Paymaster로 가스비 선지급
- EIP-1559 Base Fee 모니터링
- Batch Processing (가스비 낮은 시간대)

##### 9. Invariant Test 추가
```solidity
contract VaultInvariantTest {
    // Heir shares 합 = 100%
    function invariant_heirSharesSum() public;
    
    // 출금액 <= 입금액
    function invariant_balanceConsistency() public;
    
    // Locked 상태에서 Claim 불가
    function invariant_lockedVaultNoClaim() public;
}
```

##### 10. Phase 재조정
```
Phase 1 (2주): MVP 핵심 기능 (ETH만, Factory 패턴)
Phase 1.5 (1주): DID + Emergency Recovery
Phase 2 (2주): Account Abstraction
Phase 3 (4주): 토큰 지원 및 고도화
```

#### Result
- ✅ PRD 문서 v1.1 배포 (2026-01-12)
- ✅ 보안 취약점 사전 차단
- ✅ 개발 일정 현실화 (3주 → 4주+)
- ✅ 법적 리스크 대응 전략 수립
- ✅ 확장성 있는 아키텍처 설계

#### References
- [EIP-1167: Minimal Proxy Contract](https://eips.ethereum.org/EIPS/eip-1167)
- [OpenZeppelin Pausable](https://docs.openzeppelin.com/contracts/4.x/api/security#Pausable)
- [Commit-Reveal Pattern](https://github.com/ethereum/wiki/wiki/Safety#commit-reveal)
- [Foundry Invariant Testing](https://book.getfoundry.sh/forge/invariant-testing)

---

## 2. Smart Contract 이슈

_작성 예정_

---

## 3. Backend API 이슈

_작성 예정_

---

## 4. Frontend 이슈

_작성 예정_

---

## 5. DevOps 이슈

_작성 예정_

---

**Last Updated**: 2026-01-12
