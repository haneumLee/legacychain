# LegacyChain

> Digital Asset Inheritance Platform - Smart Contract 기반 디지털 자산 상속 플랫폼

[![Solidity](https://img.shields.io/badge/Solidity-0.8.33-blue.svg)](https://soliditylang.org/)
[![Foundry](https://img.shields.io/badge/Foundry-1.5.1-red.svg)](https://getfoundry.sh/)
[![Go](https://img.shields.io/badge/Go-1.21.13-00ADD8.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 📖 개요

LegacyChain은 블록체인 기반 디지털 자산 상속 플랫폼입니다. 스마트 컨트랙트를 통해 안전하고 투명한 자산 관리 및 상속 프로세스를 제공합니다.

### 주요 기능

- 🔐 **Commit-Reveal Heartbeat**: Front-running 공격 방지
- 👨‍👩‍👧‍👦 **Multi-Heir Support**: 최대 10명의 상속인 지정
- ⏰ **Grace Period**: Owner 복귀 기회 제공
- 🔒 **Multi-Signature Approval**: 상속인 합의 기반 자산 청구
- ⚡ **Gas Optimization**: EIP-1167 Clone 패턴으로 94.4% 가스 절감
- 🛡️ **Emergency Pause**: 긴급 상황 대응

## 🏗️ 아키텍처

```
legacychain/
├── contracts/          # Smart Contracts (Solidity + Foundry)
│   ├── src/
│   │   ├── VaultFactory.sol
│   │   └── IndividualVault.sol
│   ├── test/
│   │   ├── unit/       # 30 unit tests
│   │   └── invariant/  # 5 invariant tests
│   └── script/         # Deployment scripts
├── backend/            # Backend API (Go + Fiber)
├── frontend/           # Frontend (Next.js 14)
├── docker/             # Infrastructure
│   ├── besu/           # Besu network configs
│   └── postgres/       # Database schema
└── docs/               # Documentation
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21.13+
- Node.js 18.x+
- Foundry 1.5.1+
- Docker 29.x+

### 1. Besu Network 시작

```bash
cd docker
docker compose up -d besu-node-1 postgres redis

# 블록 생성 확인
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545
```

### 2. Smart Contract 배포

```bash
cd contracts

# 테스트 실행
forge test -vv

# Besu 네트워크에 배포
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
forge script script/DeployVaultFactory.s.sol:DeployVaultFactory \
  --rpc-url http://localhost:8545 \
  --broadcast \
  --legacy
```

### 3. 배포된 컨트랙트 주소

```
VaultFactory:    0x5FbDB2315678afecb367f032d93F642f64180aa3
Implementation:  0xa16E02E87b7454126E5E10d957A927A7F5B5d2be
Deployer:        0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Network:         Besu Clique PoA (Chain ID 1337)
```

## 🧪 Testing

### Smart Contract Tests

```bash
cd contracts

# 전체 테스트 실행
forge test

# 커버리지 확인
forge coverage

# 가스 스냅샷
forge snapshot
```

**테스트 결과:**
- ✅ 30 unit tests (100% pass)
- ✅ 5 invariant tests (256 fuzz runs each)
- ✅ 90.15% line coverage
- ✅ 0 High/Medium security issues (Slither)

## 📊 Gas Optimization

**EIP-1167 Clone Pattern 효과:**
```
Before (Direct Deploy): ~800,000 gas
After (Clone):          ~45,000 gas
Reduction:              94.4%
```

## 🔧 Development

### Smart Contract 개발

```bash
cd contracts

# 빌드
forge build

# 린트
forge fmt

# 보안 분석
slither .
```

### Backend 개발 (준비 중)

```bash
cd backend
go mod download
go run main.go
```

### Frontend 개발 (준비 중)

```bash
cd frontend
npm install
npm run dev
```

## 📚 Documentation

- [PRD (Product Requirements Document)](docs/PRD.md)
- [Development Log](docs/DEV_LOG.md)
- [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
- [Security Report](docs/SECURITY_REPORT.md)
- [API Specification](docs/API_SPEC.md)

## 🛠️ Technology Stack

**Smart Contracts:**
- Solidity 0.8.33
- Foundry (Forge, Cast, Anvil)
- OpenZeppelin Contracts 5.5.0

**Blockchain:**
- Hyperledger Besu 24.12.0
- Clique PoA Consensus

**Backend:**
- Go 1.21.13
- Fiber Framework
- PostgreSQL 16
- Redis 7

**Frontend:**
- Next.js 14
- TypeScript
- Ethers.js

**Infrastructure:**
- Docker & Docker Compose
- GitHub Actions (CI/CD)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

## 📞 Contact

For questions or support, please open an issue on GitHub.

---

**Built with ❤️ by LegacyChain Team**
