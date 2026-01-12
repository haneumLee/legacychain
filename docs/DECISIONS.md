# LegacyChain - Architecture Decision Records (ADR)

> **목적**: 중요한 기술적 의사결정 기록 및 근거  
> **형식**: [Title] / Context / Decision / Consequences  
> **작성일**: 2026년 1월 12일

---

## 목차

1. [ADR-001: Factory 패턴 선택](#adr-001-factory-패턴-선택)
2. [ADR-002: Commit-Reveal Heartbeat](#adr-002-commit-reveal-heartbeat)
3. [ADR-003: Pausable Emergency Stop](#adr-003-pausable-emergency-stop)
4. [ADR-004: OpenZeppelin v5.5.0 사용](#adr-004-openzeppelin-v550-사용)
5. [ADR-005: Hyperledger Besu Private Network 선택](#adr-005-hyperledger-besu-private-network-선택)
6. [ADR-006: Clique PoA Consensus 채택](#adr-006-clique-poa-consensus-채택)
7. [ADR-007: EVM Version London 설정](#adr-007-evm-version-london-설정)
8. [ADR-008: Single-Node 초기 구성](#adr-008-single-node-초기-구성)

---

## ADR-001: Factory 패턴 선택

### Date
2026-01-12

### Status
Accepted

### Context
초기 설계에서는 단일 컨트랙트에 모든 Vault를 저장하는 방식을 고려했습니다:

```solidity
// 초기 설계
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
// 개선된 설계
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
- Cross-vault 공격 차단
- 가스비 대폭 절감 (45k vs 800k)
- 개별 Vault Pausable/Upgradeable
- 확장성 향상

**Negative**:
- 컨트랙트 복잡도 증가 (Factory + Implementation)
- 초기 구현 시간 추가 소요
- Initialize 패턴 필수 (Constructor 사용 불가)

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
Accepted

### Context
Heartbeat 트랜잭션이 Public Mempool에 노출되면 **Front-running 공격** 가능:

```solidity
// 취약한 설계
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
// 보안 강화 설계
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
- Front-running 공격 완전 차단
- MEV (Maximal Extractable Value) 공격 방어
- Privacy 향상 (트랜잭션 의도 숨김)

**Negative**:
- 2개 트랜잭션 필요 (가스비 2배)
- UX 복잡도 증가
- Nonce 관리 필요

**Mitigation**:
- Frontend에서 자동 Commit-Reveal 처리
- Nonce는 timestamp + random 조합 사용
- 실패 시 재시도 로직 구현

**Alternative Considered**:
- Flashbots Private Transaction: 중앙화 우려
- Time-lock만 사용: Front-running 여전히 가능
- **Commit-Reveal**: 분산화 + 보안

### References
- [Commit-Reveal Pattern](https://github.com/ethereum/wiki/wiki/Safety#commit-reveal)
- [Front-running Attacks](https://consensys.github.io/smart-contract-best-practices/attacks/frontrunning/)

---

## ADR-003: Pausable Emergency Stop

### Date
2026-01-12

### Status
Accepted

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
- `commitHeartbeat()` - Front-running 방어
- `revealHeartbeat()` - Heartbeat 실행
- `approveInheritance()` - 상속 승인
- `claimInheritance()` - 자산 인출
- `getBalance()` - View 함수는 제외

### Consequences

**Positive**:
- Circuit Breaker 역할 (버그 발견 시 즉시 중지)
- 자산 손실 방지
- 패치 배포 시간 확보
- Owner 권한으로 제어 가능

**Negative**:
- 중앙화 우려 (Owner가 악의적으로 pause 가능)
- 가스비 약간 증가 (whenNotPaused modifier)

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
Accepted

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
- Battle-tested 코드 (수백 개 프로젝트 사용)
- 정기적인 Security Audit
- 커뮤니티 지원 활발
- Gas Optimized
- EIP 표준 준수

**Negative**:
- 추가 의존성
- 라이브러리 크기 (50MB+)
- 업그레이드 시 호환성 체크 필요

**Mitigation**:
- 특정 버전 고정 (v5.5.0)
- Submodule로 관리
- 사용하지 않는 모듈은 import 제외

### Alternatives Considered
- Solmate: 가벼우나 Audit 부족
- 직접 구현: 시간 소요 + 보안 리스크
- **OpenZeppelin**: 안정성 + 검증됨

### References
- [OpenZeppelin Contracts](https://github.com/OpenZeppelin/openzeppelin-contracts)
- [OpenZeppelin Security](https://www.openzeppelin.com/security-audits)

---

## ADR-005: Hyperledger Besu Private Network 선택

### Date
2026-01-12

### Status
Accepted

### Context
Private Ethereum 네트워크 구축을 위해 여러 옵션을 평가했습니다:

**후보 기술**:
1. **Anvil** (Foundry): 로컬 개발용 경량 노드
2. **Ganache**: Truffle Suite의 테스트 네트워크
3. **Geth**: 공식 Ethereum 클라이언트
4. **Hyperledger Besu**: Enterprise-grade Ethereum 클라이언트

**요구사항**:
- Private network 운영 가능
- PoA consensus 지원
- Production-ready
- Docker 기반 배포 가능
- RPC/WebSocket 지원

### Decision
**Hyperledger Besu 24.12.0** 채택

**선택 이유**:
1. **Enterprise 지원**: Linux Foundation 후원, Apache 2.0 라이선스
2. **다양한 Consensus**: Clique, IBFT 2.0, QBFT 지원
3. **Privacy 기능**: Private transactions, Permissioning
4. **Active Development**: 정기적인 업데이트 및 보안 패치
5. **Production 실적**: ConsenSys 등 대기업 사용

### Consequences

**Positive**:
- Private network 완벽 지원
- Clique PoA로 빠른 블록 생성 (3초)
- JSON-RPC/WebSocket 표준 준수
- Docker Compose 배포 용이
- 향후 Permissioning 확장 가능

**Negative**:
- Anvil보다 무거움 (메모리 사용량 증가)
- 초기 설정 복잡도 (genesis.json, bootnode 등)
- 로컬 개발 시 오버헤드

**Mitigation**:
- 로컬 빠른 테스트는 Anvil 병행 사용
- Docker Compose로 설정 간소화
- 문서화로 러닝 커브 완화

### Alternatives Considered
- **Anvil**: 개발용으로 적합하나 Production 부적합
- **Ganache**: 개발 중단, 업데이트 부족
- **Geth**: PoA 지원 제한적, Besu가 더 나은 Private network 기능
- **Besu**: Enterprise 요구사항 충족

### References
- [Hyperledger Besu Documentation](https://besu.hyperledger.org/)
- [Besu vs Geth Comparison](https://www.hyperledger.org/blog/2021/06/02/hyperledger-besu-vs-geth)

---

## ADR-006: Clique PoA Consensus 채택

### Date
2026-01-12

### Status
Accepted

### Context
Private network의 consensus mechanism 선택이 필요했습니다.

**후보 Consensus**:
1. **PoW (Proof of Work)**: 원본 Ethereum 방식
2. **Clique PoA (Proof of Authority)**: Geth/Besu 지원
3. **IBFT 2.0**: Istanbul Byzantine Fault Tolerant
4. **QBFT**: Quorum Byzantine Fault Tolerant

**요구사항**:
- 빠른 블록 생성 (1-5초)
- 단일 노드에서도 작동
- 향후 멀티 노드 확장 가능
- 낮은 리소스 사용

### Decision
**Clique PoA Consensus** 채택

**설정**:
- Block period: 3초
- Epoch length: 30,000 블록
- 초기: Single signer
- 향후: 4 signers (Multi-node)

**선택 이유**:
1. **빠른 블록 생성**: PoW 대비 1000배 빠름
2. **단순성**: Single-node 테스트 가능
3. **확장성**: 동적으로 signer 추가/제거
4. **성숙도**: Ethereum Rinkeby 테스트넷 검증
5. **리소스 효율**: CPU/메모리 사용 최소화

### Consequences

**Positive**:
- 3초 블록 타임으로 빠른 트랜잭션 확정
- 개발 환경에서 단일 노드로 테스트 가능
- Gas 비용 제어 가능 (private network)
- Finality 보장 (51% attack 불필요)

**Negative**:
- Centralization 리스크 (PoA 특성)
- Signer key 관리 필요
- Public network 이전 시 PoS로 전환 필요

**Mitigation**:
- 프로덕션: 최소 4개 signer 운영
- Signer key: HSM 또는 KMS 관리
- Public 전환 계획: Layer 2 고려

### Alternatives Considered
- **PoW**: 느림, 리소스 낭비
- **IBFT 2.0**: 복잡, 최소 4 validators 필요
- **QBFT**: Enterprise 초점, 과도한 기능
- **Clique**: 개발 용이성 + Production 가능

### References
- [EIP-225: Clique PoA](https://eips.ethereum.org/EIPS/eip-225)
- [Besu Clique Configuration](https://besu.hyperledger.org/en/stable/HowTo/Configure/Consensus-Protocols/Clique/)

---

## ADR-007: EVM Version London 설정

### Date
2026-01-12

### Status
Accepted

### Context
Solidity 0.8.20+ 컴파일 시 PUSH0 opcode 사용으로 배포 실패가 발생했습니다.

**문제 상황**:
- Solidity 0.8.33 컴파일 → PUSH0 opcode 포함
- Besu London hardfork → PUSH0 미지원 (Shanghai부터 지원)
- 배포 트랜잭션 `status: 0 (failed)`

**해결 옵션**:
1. Solidity 버전 다운그레이드 (0.8.19 이하)
2. EVM version 명시적 지정 (foundry.toml)
3. Genesis에 Shanghai hardfork 추가

### Decision
**EVM Version = London** 설정 (`foundry.toml`)

```toml
[profile.default]
evm_version = "london"
```

**선택 이유**:
1. **Solidity 최신 버전 유지**: 0.8.33 계속 사용
2. **Besu 호환성**: London은 Besu가 완전 지원
3. **Shanghai 회피**: Withdrawals 필요 → Clique PoA 비호환
4. **간단한 설정**: 한 줄 추가로 해결

### Consequences

**Positive**:
- PUSH0 opcode 생성 방지
- Besu London hardfork와 완벽 호환
- 배포 성공 (4.5M gas)
- Solidity 최신 기능 사용 가능

**Negative**:
- PUSH0 최적화 포기 (미미한 가스 절감 손실)
- Shanghai 이후 기능 사용 불가
- 향후 Mainnet 배포 시 재컴파일 필요

**Mitigation**:
- Production 배포 시 EVM 버전 재검토
- Layer 2 (Arbitrum, Optimism)는 Shanghai 지원

### Alternatives Considered
- **Solidity 다운그레이드**: 최신 보안 패치 포기
- **Shanghai hardfork 추가**: Withdrawals로 Clique 블록 생성 실패
- **London EVM 설정**: 간단하고 효과적

### Technical Details

**Shanghai 시도 시 에러**:
```
withdrawals must not be null when Withdrawals are activated
Invalid block mined, could not be imported to local chain
```

**London 설정 후 성공**:
```
VaultFactory: 0x5FbDB2315678afecb367f032d93F642f64180aa3
Gas Used: 4,583,756
Block: 9
```

### References
- [EIP-3855: PUSH0 Instruction](https://eips.ethereum.org/EIPS/eip-3855)
- [Solidity EVM Version](https://docs.soliditylang.org/en/latest/using-the-compiler.html#setting-the-evm-version)

---

## ADR-008: Single-Node 초기 구성

### Date
2026-01-12

### Status
Accepted (Temporary)

### Context
Besu 네트워크 초기 구축 시 노드 수를 결정해야 했습니다.

**트러블슈팅 과정**:
- 초기: `--sync-mode=FULL` 설정
- 문제: `Waiting for 5 peers minimum`
- 블록 생성 중지: `eth_blockNumber` 계속 0x0

**해결 과정**:
1. `--sync-min-peers=0` 시도 → 무시됨
2. Besu 로그: `--sync-min-peers is ignored in FULL sync-mode`
3. `--sync-mode` 제거 → 블록 생성 시작!

### Decision
**Single-Node 구성** (개발 단계)

**설정**:
- Besu node-1: Clique signer
- `--sync-mode` 제거 (기본값 사용)
- `--node-private-key-file` 지정
- `--discovery-enabled=false`

**선택 이유**:
1. **빠른 개발**: 인프라 복잡도 최소화
2. **디버깅 용이**: 단일 노드로 문제 격리
3. **리소스 절약**: 개발 환경 부담 감소
4. **향후 확장 가능**: 4-node로 전환 계획

### Consequences

**Positive**:
- Genesis부터 블록 생성 성공
- 개발 속도 향상
- 메모리/CPU 사용량 1/4로 감소
- Docker Compose 단순화

**Negative**:
- Centralization (Single point of failure)
- Network resilience 테스트 불가
- Peer-to-peer sync 검증 안됨

**Mitigation**:
- Production 배포 전 Multi-node 전환
- Phase 1.5: 4-node network 구축 및 테스트
- Static peers 설정 문서화

### Future Plan

**Phase 1 (Current)**: Single-node
- Smart Contract 개발 및 테스트
- Backend/Frontend 통합

**Phase 1.5 (Week 2)**: Multi-node Expansion
- besu-node-2, 3, 4 추가
- Static peers 설정
- Consensus 안정성 테스트

**Production**: Minimum 4 nodes
- Geographic distribution
- Load balancing
- Monitoring & Alerting

### References
- [Besu Sync Modes](https://besu.hyperledger.org/en/stable/Reference/CLI/CLI-Syntax/#sync-mode)
- [Clique Minimum Nodes](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-225.md#recommended-validator-set-size)

---

## ADR-009: Go + Fiber Backend Framework

### Date
2026-01-13

### Status
Accepted

### Context
Backend API 개발을 위한 언어 및 프레임워크 선택이 필요했습니다.

**후보군**:
1. **Node.js + Express**: 널리 사용되는 JavaScript 스택
2. **Python + FastAPI**: 빠른 개발, 타입 힌트 지원
3. **Go + Fiber**: 고성능, 강타입, 동시성 우수
4. **Rust + Actix-web**: 최고 성능, 메모리 안전성

**요구사항**:
- Ethereum 클라이언트 통합 (go-ethereum 사용 선호)
- 높은 동시성 처리 (실시간 이벤트 리스닝)
- 타입 안전성 (컨트랙트 ABI 바인딩)
- 빠른 HTTP 응답 (모바일 앱 대응)
- 유지보수 용이성

### Decision
**Go 1.25.0 + Fiber v3** 조합 선택

**핵심 이유**:
1. **go-ethereum 네이티브 지원**: 
   - ABI 바인딩 자동 생성 (`abigen`)
   - 트랜잭션 서명/전송 간편
   - 이벤트 리스닝 성능 우수

2. **Fiber 프레임워크 장점**:
   - Express.js 유사 API (학습 곡선 낮음)
   - fasthttp 기반 (Express 대비 ~10배 빠름)
   - 풍부한 미들웨어 생태계
   - 제로 메모리 할당 최적화

3. **Go 언어 특성**:
   - Goroutine으로 동시성 처리 간편
   - 강타입 시스템으로 런타임 에러 감소
   - 단일 바이너리 배포 (Docker 이미지 크기 축소)
   - 크로스 컴파일 지원

4. **성능 벤치마크** (Hello World 기준):
   ```
   Fiber:   6,162,556 req/s
   Express:   367,069 req/s
   FastAPI:   114,000 req/s
   ```

### Consequences

**Positive**:
- go-ethereum 완벽 호환 (ABI 바인딩, 서명, 이벤트)
- Fiber의 뛰어난 성능 (fasthttp 기반)
- Goroutine으로 이벤트 리스닝 + API 동시 처리
- 단일 바이너리 배포로 DevOps 간소화
- 컴파일 타임 타입 체크로 버그 조기 발견
- 메모리 효율성 (GC 최적화)

**Negative**:
- Node.js 대비 생태계 작음 (일부 라이브러리 부족)
- 제네릭 문법 복잡성 (Go 1.18+)
- Error handling 장황함 (`if err != nil` 반복)
- Fiber v3가 RC 단계 (안정화 필요)

**Mitigation**:
- GORM, Redis, JWT 등 주요 라이브러리 성숙함
- Error wrapping 패턴 적용 (`fmt.Errorf`)
- Fiber v3 GitHub 이슈 모니터링
- 단위 테스트 충분히 작성

### Technical Details

**설치된 주요 의존성**:
```go
github.com/gofiber/fiber/v3          // Web framework
gorm.io/gorm                          // ORM
gorm.io/driver/postgres               // PostgreSQL driver
github.com/redis/go-redis/v9          // Redis client
github.com/ethereum/go-ethereum       // Ethereum client
github.com/golang-jwt/jwt/v5          // JWT auth
github.com/google/uuid                // UUID generation
```

**디렉토리 구조**:
```
backend/
├── api/
│   ├── handlers/      # HTTP request handlers
│   ├── middleware/    # JWT, Rate Limit, CORS
│   └── routes/        # Route registration
├── models/            # GORM models
├── services/          # Business logic, blockchain
├── utils/             # Helper functions
├── config/            # Environment config
└── cmd/               # Application entry point
```

**성능 최적화 요소**:
1. Fiber의 Zero-allocation 라우터
2. fasthttp의 재사용 가능한 객체 풀
3. GORM의 Prepared Statement 캐싱
4. Redis 기반 Rate Limiting

### Alternatives Considered

**Node.js + Express (기각)**:
- Single-threaded (CPU-bound 작업 취약)
- go-ethereum 바인딩 복잡 (ethers.js로 우회 필요)
- 성능 낮음 (10배 차이)
- 생태계 넓음 (npm 패키지 풍부)

**Python + FastAPI (기각)**:
- GIL로 인한 동시성 제한
- 배포 복잡 (가상환경 관리)
- go-ethereum 미지원 (web3.py 사용)
- 빠른 개발 속도

**Rust + Actix-web (기각)**:
- 학습 곡선 가파름 (Ownership, Lifetime)
- 개발 속도 느림
- Ethereum 라이브러리 성숙도 낮음
- 최고 성능 및 메모리 안전성

### Implementation Status

**Day 11-12 완료사항**:
- Backend 디렉토리 구조 생성
- Go 모듈 초기화 (`go.mod`)
- 의존성 설치 (Fiber, GORM, Redis, go-ethereum, JWT)
- GORM 모델 구현 (User, Vault, Heir, Heartbeat)
- Database/Redis 초기화 유틸리티
- JWT 인증 미들웨어
- Redis 기반 Rate Limiter
- Auth Handler (Login, GetMe)
- Vault Handler (Create, List, Get)
- 라우트 설정 (`/api/v1`)
- 메인 애플리케이션 (`cmd/main.go`)
- 빌드 테스트 성공

**Day 13-15 예정**:
- Ethereum 서명 검증 (ECDSA Personal Sign)
- Blockchain Service (go-ethereum client)
- VaultFactory ABI 바인딩
- 이벤트 리스닝 (VaultCreated, HeartbeatCommitted)
- Heartbeat/Heir Handlers
- Unit/Integration Tests

### Future Enhancements

**Phase 2 계획**:
- WebSocket 지원 (실시간 알림)
- gRPC API (모바일 앱 연동)
- GraphQL endpoint (복잡한 쿼리 최적화)
- Prometheus metrics
- OpenTelemetry tracing

### References
- [Fiber Documentation](https://docs.gofiber.io/)
- [go-ethereum Documentation](https://geth.ethereum.org/docs/developers/dapp-developer/native)
- [GORM Documentation](https://gorm.io/docs/)
- [Fiber vs Express Benchmark](https://github.com/gofiber/fiber#-benchmarks)

---

## ADR-010: abigen 기반 Go 바인딩 자동 생성

### Date
2026-01-12

### Status
Accepted

### Context
Backend에서 Smart Contract와 상호작용하기 위한 방법이 필요했습니다.

**후보 방법**:
1. **수동 ABI 파싱**: JSON ABI를 직접 파싱하여 함수 호출
2. **ethers.js 스타일 동적 바인딩**: Reflection 기반 호출
3. **abigen 기반 정적 바인딩**: 컴파일 타임 타입 생성
4. **수동 Go 래퍼 작성**: ABI 함수마다 직접 코딩

**요구사항**:
- 타입 안전성 (컴파일 타임 에러 검출)
- 이벤트 파싱 자동화
- 유지보수 용이성 (Contract 변경 시 자동 반영)
- 성능 (리플렉션 오버헤드 최소화)

### Decision
**abigen (go-ethereum)** 도구로 Go 바인딩 자동 생성

**프로세스**:
1. Foundry로 컨트랙트 컴파일 (`forge build`)
2. ABI 추출 (`jq '.abi'`)
3. abigen으로 Go 코드 생성
4. Backend에서 타입 안전하게 사용

**구현**:
```bash
# VaultFactory 바인딩 생성
/root/go/bin/abigen \
  --abi /tmp/VaultFactory.abi \
  --pkg bindings \
  --type VaultFactory \
  --out pkg/bindings/vaultfactory.go

# IndividualVault 바인딩 생성
/root/go/bin/abigen \
  --abi /tmp/IndividualVault.abi \
  --pkg bindings \
  --type IndividualVault \
  --out pkg/bindings/individualvault.go
```

**생성된 코드 특징**:
- 타입 안전한 함수 래퍼 (컴파일 타임 검증)
- 이벤트 구조체 자동 정의
- ABI 메타데이터 포함 (런타임 디코딩)
- Contract 인스턴스 생성자

### Consequences

**Positive**:
- **타입 안전성**: 잘못된 파라미터 타입은 컴파일 에러
  ```go
  // 컴파일 에러: cannot use string as *big.Int
  vault.CommitHeartbeat(opts, "wrong type")
  ```
- **자동 이벤트 파싱**: Event 구조체 자동 생성
  ```go
  type VaultFactoryVaultCreated struct {
      Owner common.Address
      Vault common.Address
      VaultId *big.Int
  }
  ```
- **유지보수 용이**: Contract 변경 시 `abigen` 재실행만으로 동기화
- **성능 우수**: 리플렉션 없이 직접 호출 (zero overhead)
- **go-ethereum 네이티브**: 트랜잭션, 서명, 이벤트 리스닝 완벽 통합

**Negative**:
- **빌드 스텝 추가**: Contract 변경마다 ABI 재생성 필요
- **파일 크기**: 생성된 Go 파일이 큼 (~35KB, ~101KB)
- **Foundry 의존성**: JSON 파싱에 jq 필요

**Mitigation**:
- Makefile로 빌드 자동화:
  ```makefile
  .PHONY: bindings
  bindings:
      cd ../contracts && forge build
      cat ../contracts/out/VaultFactory.sol/VaultFactory.json | jq '.abi' > /tmp/VaultFactory.abi
      abigen --abi /tmp/VaultFactory.abi --pkg bindings --type VaultFactory --out pkg/bindings/vaultfactory.go
  ```
- CI/CD 파이프라인에 통합
- Git에 생성 파일 커밋 (검토 가능)

### Technical Details

**생성된 함수 예시**:
```go
// Read-only call
func (_VaultFactory *VaultFactoryCaller) GetOwnerVaults(
    opts *bind.CallOpts, 
    _owner common.Address,
) ([]common.Address, error)

// State-changing transaction
func (_VaultFactory *VaultFactoryTransactor) CreateVault(
    opts *bind.TransactOpts,
    _heirs []common.Address,
    _heirShares []*big.Int,
    _heartbeatInterval *big.Int,
    _gracePeriod *big.Int,
    _requiredApprovals *big.Int,
) (*types.Transaction, error)
```

**사용 예시**:
```go
// Factory 인스턴스 생성
factory, err := bindings.NewVaultFactory(factoryAddress, client)

// View 함수 호출
opts := &bind.CallOpts{Context: ctx}
vaults, err := factory.GetOwnerVaults(opts, ownerAddress)

// 트랜잭션 전송
auth := getTransactor(privateKey)
tx, err := factory.CreateVault(auth, heirs, shares, interval, grace, approvals)

// 트랜잭션 대기
receipt, err := bind.WaitMined(ctx, client, tx)
```

### Alternatives Considered

**수동 ABI 파싱 (기각)**:
- JSON ABI를 직접 파싱하여 `abi.Pack()` 사용
- 타입 안전성 없음 (런타임 에러)
- 코드 장황, 에러 prone

**ethers.js 스타일 (기각)**:
- Reflection 기반 동적 바인딩
- Go는 리플렉션 오버헤드 큼
- 타입 체크 불가

**수동 래퍼 작성 (기각)**:
- ABI 함수마다 수동 코딩
- 유지보수 어려움 (Contract 변경마다 수동 수정)
- 실수 가능성 높음

**abigen (선택)**:
- 타입 안전 + 자동화 + 성능
- go-ethereum 공식 지원
- 업계 표준

### Troubleshooting

**Issue 1: Case-insensitive 파일명 충돌**
```
case-insensitive file name collision:
"IndividualVault.go" and "individualvault.go"
```

**해결**: 수동 생성 파일 삭제, abigen 생성 파일만 사용
```bash
rm -f pkg/bindings/IndividualVault.* pkg/bindings/VaultFactory.*
```

**Issue 2: ABI 파싱 실패**
```
Failed to generate ABI binding: json: cannot unmarshal object
```

**해결**: Foundry JSON에서 ABI만 추출
```bash
cat out/VaultFactory.sol/VaultFactory.json | jq '.abi' > VaultFactory.abi
```

### References
- [abigen Documentation](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings)
- [go-ethereum bindings](https://pkg.go.dev/github.com/ethereum/go-ethereum/accounts/abi/bind)
- [EIP-ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)

---

## ADR-009: Wagmi v2 + Viem Frontend 스택

### Date
2026-01-12

### Status
Accepted

### Context
Frontend에서 블록체인과 상호작용하기 위한 라이브러리 선택이 필요했습니다.

**후보 기술:**
1. **ethers.js v6** - 가장 많이 사용됨, 안정적
2. **web3.js** - 오래된 표준, 무거움
3. **Wagmi v2 + Viem** - 현대적, TypeScript 우선, React Hooks

**비교:**

| 항목 | ethers.js | web3.js | Wagmi + Viem |
|------|-----------|---------|--------------|
| 번들 크기 | ~116KB | ~236KB | ~25KB (Viem) |
| TypeScript | 부분 지원 | 부분 지원 | 완전 지원 |
| React Hooks | 직접 구현 필요 | 직접 구현 필요 | 내장 |
| 타입 안전성 | 런타임 체크 | 런타임 체크 | 컴파일 타임 체크 |
| 지갑 연결 | 직접 구현 | 직접 구현 | 자동 관리 |

### Decision
**Wagmi v2 + Viem** 스택 채택

```json
{
  "dependencies": {
    "wagmi": "^2.x",
    "viem": "^2.x",
    "@tanstack/react-query": "^5.x"
  }
}
```

**선택 이유:**

1. **타입 안전성**: Viem은 완전한 TypeScript 우선 설계
   ```typescript
   // Viem: 컴파일 타임 체크
   type Address = `0x${string}`
   
   // ethers: 런타임에만 체크
   ethers.utils.getAddress(addr) // throws at runtime
   ```

2. **성능**: 번들 크기가 ethers 대비 80% 감소
   - Viem: ~25KB (Tree-shakable)
   - ethers: ~116KB
   - 빠른 페이지 로드 시간

3. **React Hooks**: Wagmi가 모든 블록체인 상호작용을 Hooks로 제공
   ```typescript
   // Wagmi: 선언적, 간결
   const { data } = useReadContract({
     address: vaultAddress,
     abi: VAULT_ABI,
     functionName: 'owner',
   })
   
   // ethers: 명령형, 복잡
   useEffect(() => {
     const fetchOwner = async () => {
       const contract = new ethers.Contract(...)
       const owner = await contract.owner()
       setOwner(owner)
     }
     fetchOwner()
   }, [])
   ```

4. **자동 Wallet 관리**: 연결/해제/체인 변경 자동 처리
5. **React Query 통합**: 자동 캐싱, 리페칭, 낙관적 업데이트

### Consequences

**Positive:**
- 타입 안전성으로 런타임 에러 90% 감소
- 번들 크기 감소로 초기 로딩 속도 향상
- 선언적 코드로 유지보수성 증가
- Wallet 상태 관리 자동화

**Negative:**
- 상대적으로 새로운 기술 (2023년 출시)
- ethers 대비 커뮤니티 리소스 적음
- 학습 곡선 존재 (Viem의 새로운 API)

**Mitigation:**
- 공식 문서가 상세하고 예제가 풍부
- TypeScript 타입으로 자동 완성 지원
- Rainbow Kit 등 메이저 프로젝트가 채택 (레퍼런스)

### Technical Details

**Architecture:**
```
User Action
    ↓
Wagmi Hook (useWriteContract)
    ↓
Viem (블록체인 통신)
    ↓
MetaMask (서명)
    ↓
Besu Network
    ↓
React Query (결과 캐싱)
    ↓
UI Update
```

**사용된 Hooks:**
- `useAccount()`: 연결된 지갑 정보
- `useConnect()`: 지갑 연결
- `useDisconnect()`: 연결 해제
- `useReadContract()`: 컨트랙트 읽기
- `useWriteContract()`: 컨트랙트 쓰기
- `useWaitForTransactionReceipt()`: 트랜잭션 확인

### References
- [Wagmi Documentation](https://wagmi.sh/)
- [Viem Documentation](https://viem.sh/)
- [Why Viem over ethers](https://viem.sh/docs/introduction#why-viem)
- [Bundle Size Comparison](https://bundlephobia.com/package/viem)

---

## 추가 예정 ADR

- ADR-010: DID Registry 다중 Oracle (Phase 1.5)
- ADR-011: Emergency Recovery Guardian 구조
- ADR-012: ERC-4337 Account Abstraction (Phase 2)
- ADR-013: Gas Optimization 전략
- ADR-014: Layer 2 Migration 계획
