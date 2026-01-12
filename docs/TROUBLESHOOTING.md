# LegacyChain - Troubleshooting Guide

> **목적**: 에러 해결 및 개선 히스토리 기록  
> **작성일**: 2026년 1월 12일

---

## 목차

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
    // 보안 격리
    // 가스비 최적화
    // Upgrade 유연성
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
- PRD 문서 v1.1 배포 (2026-01-12)
- 보안 취약점 사전 차단
- 개발 일정 현실화 (3주 → 4주+)
- 법적 리스크 대응 전략 수립
- 확장성 있는 아키텍처 설계

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

## 5. DevOps 이슈

### [2026-01-12] Besu Clique PoA 블록 생성 실패

#### Date
2026-01-12

#### Symptom
- Besu 노드 정상 시작되지만 블록 번호가 0에서 증가하지 않음
- RPC 응답: `{"jsonrpc":"2.0","id":1,"result":"0x0"}` (지속)
- 로그: `Unable to find sync target. Waiting for 5 peers minimum`
- VaultFactory 배포 시 "Known transaction" 에러 (mempool에 있지만 채굴 안됨)

#### Root Cause
1. **Full Sync Mode Peer Requirement**
   - `--sync-mode=FULL`로 실행 중
   - Full sync에서 `--sync-min-peers` 무시됨: `--sync-min-peers is ignored in FULL sync-mode`
   - 기본 5개 피어 대기 → 단일 노드에서 블록 생성 불가

2. **Network Configuration**
   - `--discovery-enabled=false` + 단일 노드 실행
   - 피어 없음 → Full Sync 시작 불가

3. **Clique Signer Setup**
   - Private key 파일 생성 
   - `--node-private-key-file` 설정  
   - But: Full Sync mode가 마이닝 차단 중 

#### Solution Options
- **Option 1**: `--sync-mode=FAST` + `--sync-min-peers=0`
- **Option 2**: `--sync-mode` 제거 (genesis부터 시작 시)
- **Option 3**: 멀티 노드 네트워크 구축 (4 nodes)

#### Next Steps
1. Option 2 시도 - sync-mode 제거
2. 실패 시 Option 1 - FAST sync
3. 성공 시 DEV_LOG.md 업데이트


### [2026-01-12] Besu Clique PoA 블록 생성 실패 - **RESOLVED**

#### Final Solution
**Option 2 적용**: `--sync-mode` 제거 + EVM 버전 설정

#### Implementation Steps
1. **docker-compose.yml**: `--sync-mode=FULL` 및 `--sync-min-peers=0` 제거
   - Single-node Clique는 sync mode 불필요
   
2. **foundry.toml**: `evm_version = "london"` 추가
   - Solidity 0.8.20+의 PUSH0 opcode 방지
   - Besu London 하드포크까지만 지원
   
3. **Data reset**: `docker volume rm docker_besu-node-1-data`
   - Genesis 변경 사항 적용 위해 필수

#### Deployment Result
VaultFactory deployed successfully
- Address: `0x5FbDB2315678afecb367f032d93F642f64180aa3`
- Implementation: `0xa16E02E87b7454126E5E10d957A927A7F5B5d2be`
- Block: 9
- Gas Used: 4,583,756

#### Key Learnings
1. Full Sync mode는 단일 노드 Clique에 부적합 (피어 대기)
2. Solidity 0.8.20+ PUSH0 opcode는 London EVM까지만 지원되지 않음
3. Shanghai hardfork는 Withdrawals 필요 → Clique PoA와 비호환
4. Genesis 변경 시 반드시 Docker volume 삭제 필요

