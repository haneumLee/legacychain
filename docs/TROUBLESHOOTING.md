# LegacyChain - Troubleshooting Guide

> **목적**: 에러 해결 및 개선 히스토리 기록  
> **작성일**: 2026년 1월 12일

---

## 목차

1. [DevOps 이슈](#1-devops-이슈)
2. [PRD 문서 개선 히스토리](#2-prd-문서-개선-히스토리)
3. [Smart Contract 이슈](#3-smart-contract-이슈)
4. [Backend API 이슈](#4-backend-api-이슈)
5. [Frontend 이슈](#5-frontend-이슈)

---

## 1. DevOps 이슈

### [2026-01-12] Besu 네트워크 가스 가격 불일치로 인한 트랜잭션 Pending

#### Date
2026-01-12

#### Error
Forge script로 Besu 네트워크에 컨트랙트 배포 시 트랜잭션이 무한정 pending 상태로 고착

#### Symptom
```bash
forge script script/DeployVaultFactory.s.sol --rpc-url http://localhost:8545 --broadcast

# Output
Waiting for transaction receipt... (60+ seconds)
[00:01:00] [----] 0/1 receipts

# 트랜잭션 상태
eth_getTransactionByHash:
  blockHash: null
  gasPrice: 0xf (15 wei)
  maxFeePerGas: 0xf (15 wei)
```

블록은 계속 생성되지만 (0xebf → 0xf60, +161 blocks) 트랜잭션은 블록에 포함되지 않음

#### Cause
Forge가 자동으로 설정한 가스 가격 (15 wei)이 Besu Clique PoA 네트워크의 최소 가스 가격 (1000 wei)보다 낮음

```bash
# Besu 네트워크 최소 가스 가격 확인
curl -X POST http://localhost:8545 -d '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}'
# Result: 0x3e8 (1000 wei)

# Forge 트랜잭션 가스 가격
gasPrice: 0xf (15 wei) ← 너무 낮음
```

**Besu Clique PoA 특성:**
- 고정된 최소 가스 가격 (min-gas-price) 존재
- Anvil/Hardhat 개발 네트워크와 달리 자동 가스 가격 조정 없음
- 낮은 가스 가격의 트랜잭션은 트랜잭션 풀에서 무한정 대기하거나 거부됨

#### Fix
명시적으로 가스 가격을 Besu 네트워크에 맞게 설정:

```bash
PRIVATE_KEY=0xac09... forge script script/DeployVaultFactory.s.sol \
  --rpc-url http://localhost:8545 \
  --broadcast \
  --legacy \
  --gas-price 1000  # Besu 최소 가스 가격 (1000 wei)
```

**결과:**
```
Contract Address: 0x5FbDB2315678afecb367f032d93F642f64180aa3
Block: 4378
Gas Used: 4,583,756
Gas Price: 1000 wei
Total Cost: 0.000000004583756 ETH
```

#### Reference
- Besu Min Gas Price 설정: `docker/besu/genesis.json` → `config.clique.minGasPrice`
- Foundry 가스 가격 설정: https://book.getfoundry.sh/reference/forge/forge-script#general-options
- Besu 트랜잭션 풀: https://besu.hyperledger.org/public-networks/concepts/transactions/pool

---

### [2026-01-12] Besu Node Healthcheck Failure

#### Date
2026-01-12

#### Symptom
besu-node-1 컨테이너가 `unhealthy` 상태로 표시되며, 다른 노드들은 `healthy` 상태

#### Root Cause Analysis

1. **Healthcheck 설정 문제**: docker-compose.yml의 healthcheck가 `curl`을 사용하도록 설정됨
2. **Besu 이미지 특성**: hyperledger/besu:24.12.0 이미지에 `curl` 도구가 포함되지 않음
3. **Healthcheck 실패 로그**:
   ```
   OCI runtime exec failed: exec failed: unable to start container process: 
   exec: "curl": executable file not found in $PATH
   ```
4. **Docker 캐싱 이슈**: `docker compose restart` 명령으로는 healthcheck 설정 변경이 적용되지 않음

#### Solution Implemented

1. **Healthcheck 명령 변경**: `curl -f http://localhost:8545` → `pgrep -f besu`
   ```yaml
   healthcheck:
     test: ["CMD", "pgrep", "-f", "besu"]
     interval: 10s
     timeout: 5s
     retries: 5
     start_period: 10s
   ```

2. **컨테이너 재생성**: 설정 변경 적용을 위해 `--force-recreate` 사용
   ```bash
   docker compose up -d --force-recreate besu-node-1
   ```

3. **검증**:
   ```bash
   docker exec legacychain-besu-node-1 pgrep -f besu
   # Output: 1 (성공)
   
   docker compose ps
   # besu-node-1: Up 24 seconds (healthy)
   ```

#### Impact
- 문제: 모니터링 시스템에서 node-1이 unhealthy로 잘못 표시
- 실제: 네트워크는 정상적으로 블록 생성 중 (Block #1,100+)
- 해결: 모든 노드 healthy 상태 확인

#### Lessons Learned
1. Docker 이미지의 포함된 도구를 가정하지 말고 확인할 것
2. Healthcheck는 최소한의 의존성만 사용 (프로세스 확인이 가장 안전)
3. docker-compose.yml 변경 후 `restart`가 아닌 `up -d --force-recreate` 사용
4. start_period를 추가하여 초기 시작 시간 고려

---

### [2026-01-12] Deleted File Mysterious Reappearance

#### Date
2026-01-12

#### Symptom
환경 변수 기반 Private Key 관리로 마이그레이션한 후 삭제한 `docker/besu/node1-key` 파일이 복구됨

#### Timeline
- 13:55: `docker/.env` 생성
- 14:16: `node1-key` 파일 생성 (65 bytes)
- 14:30: `entrypoint.sh` 생성
- 15:xx: 사용자가 파일 복구 발견

#### Root Cause Analysis

1. **Git 히스토리**: 파일이 이전 커밋 `684c760`에 포함되어 있음
2. **가능한 원인**:
   - VSCode/Editor의 자동 복구 기능
   - Git 명령 실행 (예: `git checkout`)
   - Docker 볼륨 마운트 설정 잔여
   - 백그라운드 프로세스의 파일 재생성

3. **확인 사항**:
   - `docker-compose.yml`에 node1-key 참조 없음
   - `.gitignore`에 `docker/.env` 설정 완료
   - Git status에 추적되지 않음

#### Solution Implemented

1. **파일 재삭제**:
   ```bash
   rm -f docker/besu/node1-key
   ```

2. **Git 히스토리에서 완전 제거** (권장):
   ```bash
   git filter-branch --force --index-filter \
     'git rm --cached --ignore-unmatch docker/besu/node*-key' \
     --prune-empty --tag-name-filter cat -- --all
   ```

3. **환경 기반 시스템 검증**:
   - entrypoint.sh가 정상 실행
   - 모든 노드가 환경 변수에서 키 로드
   - 네트워크 정상 작동 (Block #1,100+)

#### Prevention
1. 민감한 파일은 처음부터 `.gitignore`에 추가
2. 이미 커밋된 파일은 Git 히스토리에서 제거
3. 정기적으로 `git status --porcelain`으로 추적 파일 확인

---

## 2. PRD 문서 개선 히스토리

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

Phase 1 (2주): MVP 핵심 기능 (ETH만, Factory 패턴)
Phase 1.5 (1주): DID + Emergency Recovery
Phase 2 (2주): Account Abstraction
Phase 3 (4주): 토큰 지원 및 고도화

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

## 3. Smart Contract 이슈

### [2026-01-12] Solidity PUSH0 Opcode 호환성 이슈

#### Date
2026-01-12

#### Symptom
VaultFactory 배포 시 트랜잭션 실패 (status: 0)
```
Error:
  status: 0 (failed)
  gas used: 4,583,756
Warning: EIP-3855 is not supported by the detected EVM version
```

#### Root Cause Analysis

1. **Solidity 버전 충돌**:
   - VaultFactory 및 IndividualVault: Solidity 0.8.20+ 사용
   - OpenZeppelin Contracts v5.5.0 최소 요구사항: ^0.8.20

2. **PUSH0 Opcode 이슈**:
   - Solidity 0.8.20부터 PUSH0 opcode(EIP-3855) 기본 사용
   - PUSH0는 Shanghai 하드포크(2023년 4월)에서 도입
   - Hyperledger Besu 24.12.0은 Shanghai까지 지원하지만 Clique PoA는 Withdrawals 비호환

3. **Shanghai + Clique 비호환성**:
   ```
   Invalid genesis block. withdrawals must not be null when Withdrawals are activated
   ```
   - Shanghai 하드포크는 `withdrawals` 필드 필수
   - Clique PoA는 PoS 전용 Withdrawals 미지원
   - Genesis JSON에 `withdrawals: null` 설정해도 실패

#### Solution Implemented

**Step 1**: EVM 버전을 London으로 다운그레이드
```toml
# contracts/foundry.toml
[profile.default]
evm_version = "london"  # Shanghai 대신 London 사용
```

**Step 2**: 컨트랙트 재빌드 및 배포
```bash
cd contracts
forge clean
forge build
forge script script/DeployVaultFactory.s.sol --rpc-url http://localhost:8545 --broadcast
```

**Step 3**: Deployment 성공 확인
```
VaultFactory: 0x5FbDB2315678afecb367f032d93F642f64180aa3
Implementation: 0xa16E02E87b7454126E5E10d957A927A7F5B5d2be
Block: 9
Gas Used: 4,583,756
Status: Success (1)
```

#### Impact

**Before (Shanghai EVM + PUSH0)**:
- Deployment 실패 (status: 0)
- Clique PoA와 비호환
- 트랜잭션 가스만 소비

**After (London EVM)**:
- Deployment 성공
- Clique PoA와 호환
- 기능 손실 없음 (PUSH0는 최적화 기능)
- Gas 미세 증가 (~1-2%)

#### Alternative Approaches Considered

1. **Besu 네트워크를 Shanghai로 업그레이드**:
   - Clique PoA가 Withdrawals 미지원 → 실행 불가
   - 대안: PoS로 마이그레이션 (프로덕션 단계)

2. **Solidity 버전 다운그레이드 (0.8.19 이하)**:
   - OpenZeppelin Contracts v5.5.0이 ^0.8.20 요구
   - 보안 패치 및 신규 기능 손실

3. **Custom EVM 설정**:
   - London EVM 사용 (최종 선택)
   - Solidity 0.8.20+ 유지 가능
   - OpenZeppelin v5.5.0 호환 유지
   - PUSH0만 비활성화

#### Lessons Learned

1. **EVM 버전 호환성 확인**: Private network 구축 시 Solidity 버전과 EVM 버전 사전 매칭
2. **Consensus 메커니즘 제약**: PoA(Clique)는 PoS 전용 기능(Withdrawals) 미지원
3. **Foundry 설정 우선순위**: `foundry.toml`의 `evm_version` 설정으로 호환성 제어
4. **Gas 최적화 vs 호환성**: PUSH0 opcode는 ~1-2% gas 절감이지만 호환성이 더 중요

#### References

- [EIP-3855: PUSH0 Instruction](https://eips.ethereum.org/EIPS/eip-3855)
- [Solidity 0.8.20 Release Notes](https://soliditylang.org/blog/2023/02/22/solidity-0.8.20-release-announcement/)
- [Besu EVM Compatibility](https://besu.hyperledger.org/en/stable/public-networks/concepts/the-merge/)
- [Foundry EVM Version Configuration](https://book.getfoundry.sh/reference/config/solidity-compiler#evm_version)

---

### [2026-01-12] Balance Snapshot 미구현으로 인한 불공정 분배

#### Date
2026-01-12 (초기 구현 단계에서 발견 및 수정)

#### Symptom
Heir들이 지분율대로 상속받지 못하는 문제

**시나리오**:
- Vault 잔액: 10 ETH
- Heir1(50%), Heir2(30%), Heir3(20%)
- Heir1이 먼저 claim → 5 ETH 인출 (잔액 5 ETH 남음)
- Heir2가 claim → 5 ETH * 30% = 1.5 ETH (원래 3 ETH여야 함)
- Heir3가 claim → 3.5 ETH * 20% = 0.7 ETH (원래 2 ETH여야 함)

#### Root Cause

초기 구현에서 현재 잔액 기준으로 비율 계산:
```solidity
// Before (잘못된 로직)
function claimInheritance() external {
    uint256 balance = address(this).balance;  // 매번 변하는 현재 잔액
    uint256 amount = (balance * share) / 10000;
    payable(msg.sender).transfer(amount);
}
```

**문제점**:
1. 먼저 claim하는 Heir가 정확한 금액 수령
2. 나중에 claim하는 Heir가 줄어든 잔액 기준으로 계산
3. 총 분배 금액이 100%를 초과하거나 미달 가능

#### Solution Implemented

Grace Period 종료 시점(첫 claim) 잔액을 스냅샷:
```solidity
// After (올바른 로직)
uint256 public totalBalanceAtUnlock;  // 스냅샷 저장

function claimInheritance() external {
    require(config.unlocked, "Not unlocked");
    
    // 첫 번째 claim 시 잔액 스냅샷
    if (config.totalBalanceAtUnlock == 0) {
        config.totalBalanceAtUnlock = address(this).balance;
        emit BalanceSnapshotTaken(config.totalBalanceAtUnlock);
    }
    
    uint16 share = heirShares[msg.sender];
    uint256 amount = (config.totalBalanceAtUnlock * share) / 10000;
    
    require(!hasClaimed[msg.sender], "Already claimed");
    hasClaimed[msg.sender] = true;
    
    payable(msg.sender).transfer(amount);
}
```

#### Validation

**Unit Test 추가**:
```solidity
function test_FairDistribution() public {
    // 10 ETH Vault
    vm.deal(address(vault), 10 ether);
    
    // Grace Period 종료 후 unlock
    vm.warp(block.timestamp + 40 days);
    vault.checkAndUnlock();
    
    // Heir1 claim (50% = 5 ETH)
    vm.prank(heir1);
    vault.claimInheritance();
    assertEq(heir1.balance, 5 ether);
    
    // Heir2 claim (30% = 3 ETH)
    vm.prank(heir2);
    vault.claimInheritance();
    assertEq(heir2.balance, 3 ether);
    
    // Heir3 claim (20% = 2 ETH)
    vm.prank(heir3);
    vault.claimInheritance();
    assertEq(heir3.balance, 2 ether);
    
    // Vault 잔액 0
    assertEq(address(vault).balance, 0);
}
```

**Invariant Test 추가**:
```solidity
// 청구액 <= 스냅샷 잔액
function invariant_balanceConsistency() public view {
    if (vault.config().totalBalanceAtUnlock > 0) {
        uint256 totalClaimed = /* calculate from events */;
        assert(totalClaimed <= vault.config().totalBalanceAtUnlock);
    }
}
```

#### Impact

**Before**:
- Heir1: 5.0 ETH (정확)
- Heir2: 1.5 ETH (3.0 ETH여야 함)
- Heir3: 0.7 ETH (2.0 ETH여야 함)
- Total: 7.2 ETH (2.8 ETH 손실)

**After**:
- Heir1: 5.0 ETH
- Heir2: 3.0 ETH
- Heir3: 2.0 ETH
- Total: 10.0 ETH (공정 분배)

#### Lessons Learned

1. **State Mutation 주의**: 잔액처럼 변하는 값을 기준으로 계산 금지
2. **Critical Point Snapshot**: 중요한 시점(unlock)의 상태를 저장
3. **Unit Test 필수**: Edge case (순차 claim) 반드시 테스트
4. **Invariant 검증**: 총 분배액 = 100% 불변성 확인

#### References

- [IndividualVault.sol](../contracts/src/IndividualVault.sol#L234-L249)
- [Unit Test - Fair Distribution](../contracts/test/unit/IndividualVault.t.sol#L412)

---

## 4. Backend API 이슈

### [2026-01-12] Blockchain Service ABI 구조 불일치

#### Date
2026-01-12

#### Symptom
- `not enough arguments in call to s.vaultFactory.CreateVault`
- 3개 파라미터로 호출했지만 실제 함수는 6개 필요

#### Root Cause
Initial implementation에서 Solidity contract signature를 완전히 파악하지 않음:
```solidity
// 실제 VaultFactory.sol
function createVault(
    address[] memory _heirs,
    uint16[] memory _heirShares,
    uint256 _heartbeatInterval,
    uint256 _gracePeriod,
    uint8 _requiredApprovals
) external returns (address)
```

#### Solution
Foundry 출력 확인 후 함수 시그니처 수정:
```go
// Before
vault, tx, err := s.vaultFactory.CreateVault(
    opts,
    heirs,
    shares,
    big.NewInt(int64(heartbeatInterval)),
)

// After
vault, tx, err := s.vaultFactory.CreateVault(
    opts,
    heirs,
    shares,
    big.NewInt(int64(heartbeatInterval)),
    big.NewInt(int64(gracePeriod)),
    uint8(requiredApprovals),
)
```

#### References
- [VaultFactory.sol](../contracts/src/VaultFactory.sol)
- `forge build --extra-output abi`

---

### [2026-01-12] Config 구조체 필드 누락 - Heirs/HeirShares

#### Date
2026-01-12

#### Symptom
```
config.Heirs undefined (type bindings.IndividualVaultVaultConfig has no field or method Heirs)
config.HeirShares undefined
```

#### Root Cause
Solidity `VaultConfig` struct에 heirs/shares 배열이 없음:
```solidity
struct VaultConfig {
    address owner;
    uint256 heartbeatInterval;
    uint256 gracePeriod;
    uint256 lastHeartbeat;
    bool unlocked;
    bool gracePeriodActive;
    uint8 requiredApprovals;
    uint8 approvalCount;
}
```

Heirs와 Shares는 별도 mapping으로 관리:
```solidity
mapping(address => bool) public isHeir;
mapping(address => uint16) public heirShares;
```

#### Solution
GetVaultConfig에서 Heirs/HeirShares를 nil로 설정:
```go
return &model.VaultConfig{
    Owner:              config.Owner.Hex(),
    HeartbeatInterval:  config.HeartbeatInterval.Uint64(),
    GracePeriod:        config.GracePeriod.Uint64(),
    LastHeartbeat:      config.LastHeartbeat.Uint64(),
    Unlocked:           config.Unlocked,
    GracePeriodActive:  config.GracePeriodActive,
    RequiredApprovals:  config.RequiredApprovals,
    ApprovalCount:      config.ApprovalCount,
    Heirs:              nil, // Separate contract call needed
    HeirShares:         nil, // Separate contract call needed
}, nil
```

#### Next Steps
추후 필요 시 별도 함수로 Heir 목록 조회 추가

---

### [2026-01-12] RevealHeartbeat Nonce 타입 불일치

#### Date
2026-01-12

#### Symptom
```
cannot use nonce (variable of type *big.Int) as [32]byte value in argument to vault.RevealHeartbeat
```

#### Root Cause
Solidity 함수는 `bytes32` 타입 요구:
```solidity
function revealHeartbeat(bytes32 _nonce) external
```

Go에서 `big.Int`로 변환한 후 [32]byte로 캐스팅 실패

#### Solution
Hex string → bytes → [32]byte 고정 배열 변환:
```go
// Nonce를 hex string에서 bytes로 변환
nonceBytes, err := hexutil.Decode(heartbeat.Nonce)
if err != nil {
    return "", fmt.Errorf("invalid nonce format: %w", err)
}

// [32]byte 고정 배열로 변환
var nonceArray [32]byte
copy(nonceArray[:], nonceBytes)

// Blockchain 호출
tx, err := vault.RevealHeartbeat(opts, nonceArray)
```

#### References
- [go-ethereum/common/hexutil](https://pkg.go.dev/github.com/ethereum/go-ethereum/common/hexutil)
- Solidity bytes32 mapping to Go [32]byte

---

### [2026-01-12] Heir Approval 함수명 변경

#### Date
2026-01-12

#### Symptom
```
vault.ApproveHeirClaim undefined (type *bindings.IndividualVault has no field or method ApproveHeirClaim)
```

#### Root Cause
Contract 함수명을 추측하여 `ApproveHeirClaim()` 사용했지만, 실제 Solidity는:
```solidity
function approveInheritance() external
```

#### Solution
ABI binding 확인 후 정확한 함수명 사용:
```go
// Before
tx, err := vault.ApproveHeirClaim(opts)

// After
tx, err := vault.ApproveInheritance(opts)
```

#### Lesson Learned
- abigen 생성 Go 파일 먼저 확인
- Solidity 원본 소스 review
- 함수명 추측 금지

---

### [2026-01-12] Heir Approval Logic 오해

#### Date
2026-01-12

#### Symptom
Initially designed approval flow as "Heir A approves Heir B to claim"

#### Root Cause
Multi-sig contract 설계 오해:
- **잘못된 이해**: 각 Heir가 다른 Heir를 승인
- **올바른 설계**: 각 Heir가 자신의 승인을 등록 (multi-sig 투표)

```solidity
// 실제 로직
mapping(address => bool) public hasApproved;
uint8 public approvalCount;

function approveInheritance() external {
    require(isHeir[msg.sender], "Not an heir");
    require(!hasApproved[msg.sender], "Already approved");
    
    hasApproved[msg.sender] = true;
    approvalCount++;
}
```

#### Solution
API에서 target heir 파라미터 제거, 본인 승인만 가능하도록 수정:
```go
// Before
type ApproveRequest struct {
    VaultID    string `json:"vault_id"`
    TargetHeir string `json:"target_heir"` // 불필요
}

// After
type ApproveRequest struct {
    VaultID string `json:"vault_id"` // 본인 승인만
}

// Handler
func (h *HeirHandler) ApproveInheritance(c *fiber.Ctx) error {
    // JWT에서 heir address 추출
    heirAddress := c.Locals("address").(string)
    
    // Blockchain 호출 (msg.sender = heirAddress)
    txHash, err := h.blockchain.ApproveInheritance(vaultAddress, heirAddress)
    // ...
}
```

#### Result
- Multi-sig pattern 정확히 구현
- 과반수(n/2 + 1) 승인 필요
- 각 Heir는 1번만 승인 가능

---

---

## 5. Frontend 이슈

작성 예정

---

## 6. DevOps 이슈 (Legacy)

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

---

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

---

**Last Updated**: 2026-01-12
---

### [2026-01-12] Multi-node Besu Peer Discovery 실패 - **RESOLVED**

#### Date
2026-01-12

#### Context
4-node Clique PoA 네트워크 구축 중 peer 연결 실패

#### Symptom 1: Bootnode with Docker Hostname
```
Invalid enode URL syntax 'enode://...@besu-node-1:30303'.
Enode URL should have format 'enode://<node_id>@<ip>:<listening_port>'.
Invalid ip address.
```

**Root Cause**: Besu enode parser가 hostname 미지원 (IP 주소만 허용)

**시도한 해결책**: Docker Compose의 `--bootnodes` 파라미터에 hostname 사용
```yaml
command: >
  --bootnodes=enode://...@besu-node-1:30303
```
**실패**: Besu가 hostname을 IP로 변환하지 못함

---

#### Symptom 2: Discovery Enabled (No Peers)
```bash
$ curl -X POST --data '{"jsonrpc":"2.0","method":"admin_peers","params":[],"id":1}' http://localhost:8545
{"result": []}  # 0 peers after 15 seconds
```

**Root Cause**: Discovery protocol은 public network 용으로 설계됨. Private network에서는:
- Discovery 시간 소요 (수분)
- Explicit peer list 없으면 찾기 어려움

**시도한 해결책**:
```yaml
--discovery-enabled=true
```
**실패**: 15초 대기 후에도 0 peers

---

#### Final Solution: Static Nodes with Fixed IP Addresses

**Step 1**: Docker network에 subnet 정의
```yaml
networks:
  legacychain-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

**Step 2**: 각 노드에 고정 IP 할당
```yaml
besu-node-1:
  networks:
    legacychain-network:
      ipv4_address: 172.20.0.11

besu-node-2:
  networks:
    legacychain-network:
      ipv4_address: 172.20.0.12

besu-node-3:
  networks:
    legacychain-network:
      ipv4_address: 172.20.0.13

besu-node-4:
  networks:
    legacychain-network:
      ipv4_address: 172.20.0.14
```

**Step 3**: static-nodes.json 생성 (IP 주소 사용)
```json
[
  "enode://8318535b54105d4a...@172.20.0.11:30303",
  "enode://4f5500f9c469ebf1...@172.20.0.12:30303",
  "enode://a5bbe97420f19da5...@172.20.0.13:30303",
  "enode://fba9dceba6fcb998...@172.20.0.14:30303"
]
```

**Step 4**: static-nodes.json 마운트 + discovery 비활성화
```yaml
volumes:
  - ./besu/static-nodes.json:/data/static-nodes.json
command: >
  --discovery-enabled=false
```

#### Validation Results

**Peer Connectivity** (즉시 연결):
```bash
Node 1: 3 peers
Node 2: 3 peers
Node 3: 3 peers
Node 4: 3 peers
```

**Block Production** (10초 후):
```bash
Node 1: Block 9 (0x9)
Node 2: Block 9 (0x9)
Node 3: Block 9 (0x9)
Node 4: Block 9 (0x9)
```

**Validator Registration**:
```bash
$ curl -X POST --data '{"jsonrpc":"2.0","method":"clique_getSigners","params":["latest"],"id":1}' http://localhost:8545
{
  "result": [
    "0x06b520312ae3b741d00664701c06e88345d93068",
    "0x30d4601b89e08b10b519f9fffead36751f714b1c",
    "0x75364fd8b8529409af49cc38273312147d46dd95",
    "0xdc7db81ff06fbfa3ffede0bc643957e6584f11ec"
  ]
}
```

#### Key Learnings

1. **Besu enode URL은 IP만 허용**: Hostname 사용 불가
2. **Discovery는 private network에 부적합**: Static nodes가 표준
3. **Genesis 변경 시 volume 삭제 필수**: `docker compose down -v`
4. **Validator에도 초기 잔액 할당**: 컨트랙트 배포 위해 필수

#### Files Modified
- `docker/docker-compose.yml`: Subnet + fixed IPs + static-nodes mount
- `docker/besu/static-nodes.json`: 4 enode URLs with IPs
- `docker/besu/genesis.json`: 4 validators in extraData + alloc

#### References
- [Besu - Static Nodes](https://besu.hyperledger.org/en/stable/private-networks/how-to/configure/static-nodes/)
- [Clique PoA Specification](https://eips.ethereum.org/EIPS/eip-225)
- [Docker Network Configuration](https://docs.docker.com/compose/networking/)

---

**Last Updated**: 2026-01-12
## 7. 블록체인 연동 이슈

### [2026-01-12] abigen 바인딩 파일명 충돌 (Case-insensitive)

#### Date
2026-01-12

#### Symptom
Backend 빌드 실패:
```
case-insensitive file name collision:
"IndividualVault.go" and "individualvault.go"
```

#### Root Cause Analysis

1. **이전 수동 생성 파일**: 대문자로 시작하는 파일명
   - `pkg/bindings/IndividualVault.go`
   - `pkg/bindings/VaultFactory.go`

2. **abigen 자동 생성**: 소문자로 시작하는 파일명
   - `pkg/bindings/individualvault.go`
   - `pkg/bindings/vaultfactory.go`

3. **운영체제 특성**:
   - Linux: Case-sensitive 파일 시스템
   - 그러나 Go 빌드 시스템은 package 충돌 감지

#### Solution Implemented

```bash
cd /root/legacychain/backend/pkg/bindings
rm -f IndividualVault.* VaultFactory.*
ls -la
# Only individualvault.go and vaultfactory.go remain
```

#### Impact

**Before**:
- 빌드 실패 (충돌 에러)
- 4개 파일 존재 (중복)

**After**:
- 빌드 성공
- 2개 파일만 유지 (abigen 생성본)
- 깔끔한 바인딩 디렉토리

#### Lessons Learned

1. **일관된 네이밍**: abigen 기본 설정 사용 (소문자 시작)
2. **자동 생성 파일 우선**: 수동 작성 파일 피하기
3. **Git에 생성 파일 커밋**: 코드 리뷰 가능하도록

#### Prevention

- `.gitignore`에 수동 파일 추가
- Makefile로 자동화:
  ```makefile
  clean-bindings:
      rm -f pkg/bindings/*.go
  
  generate-bindings: clean-bindings
      # abigen commands...
  ```
- CI/CD에서 자동 검증

---

### [2026-01-12] Foundry ABI JSON 형식 이슈

#### Date
2026-01-12

#### Symptom
abigen 실행 실패:
```
Fatal: Failed to generate ABI binding: 
json: cannot unmarshal object into Go value of type []struct {...}
```

#### Root Cause

Foundry의 `out/` 디렉토리 JSON 파일은 전체 메타데이터 포함:
```json
{
  "abi": [...],
  "bytecode": {...},
  "deployedBytecode": {...},
  "methodIdentifiers": {...},
  ...
}
```

abigen은 순수 ABI 배열만 필요:
```json
[
  {"type": "constructor", ...},
  {"type": "function", "name": "createVault", ...},
  ...
]
```

#### Solution

`jq`로 ABI 필드만 추출:
```bash
cat out/VaultFactory.sol/VaultFactory.json | jq '.abi' > /tmp/VaultFactory.abi
abigen --abi /tmp/VaultFactory.abi --pkg bindings --type VaultFactory --out vaultfactory.go
```

#### Alternative

Foundry의 `--extra-output abi` 옵션 사용:
```bash
forge build --extra-output abi
# out/VaultFactory.sol/VaultFactory.abi 생성됨
```

#### References
- [Foundry Artifacts](https://book.getfoundry.sh/reference/forge/forge-build#description)
- [abigen Usage](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings)

---

**Last Updated**: 2026-01-12
