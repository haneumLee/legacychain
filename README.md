# LegacyChain

> Digital Asset Inheritance Platform - Smart Contract 기반 디지털 자산 상속 플랫폼

[![Solidity](https://img.shields.io/badge/Solidity-0.8.33-blue.svg)](https://soliditylang.org/)
[![Foundry](https://img.shields.io/badge/Foundry-1.5.1-red.svg)](https://getfoundry.sh/)
[![Go](https://img.shields.io/badge/Go-1.25.0-00ADD8.svg)](https://go.dev/)
[![Fiber](https://img.shields.io/badge/Fiber-v3.0.0--rc.3-00ACD7.svg)](https://gofiber.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 개요

LegacyChain은 블록체인 기반 디지털 자산 상속 플랫폼입니다. 스마트 컨트랙트를 통해 안전하고 투명한 자산 관리 및 상속 프로세스를 제공합니다.

### 주요 기능

- **Commit-Reveal Heartbeat**: Front-running 공격 방지
- **Multi-Heir Support**: 최대 10명의 상속인 지정
- **Grace Period**: Owner 복귀 기회 제공
- **Multi-Signature Approval**: 상속인 합의 기반 자산 청구
- **Gas Optimization**: EIP-1167 Clone 패턴으로 94.4% 가스 절감
- **Emergency Pause**: 긴급 상황 대응

## 아키텍처

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

## Quick Start

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
# WARNING: 테스트 환경 전용 - 절대 실제 자산을 보내지 마세요!
# .env.besu 파일에서 BESU_TEST_PRIVATE_KEY 사용 권장
source ../.env.besu
forge script script/DeployVaultFactory.s.sol:DeployVaultFactory \
  --rpc-url $BESU_RPC_URL \
  --private-key $BESU_TEST_PRIVATE_KEY \
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

## Testing

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
- 30 unit tests (100% pass)
- 5 invariant tests (256 fuzz runs each)
- 90.15% line coverage
- 0 High/Medium security issues (Slither)

## Gas Optimization

**EIP-1167 Clone Pattern 효과:**
```
Before (Direct Deploy): ~800,000 gas
After (Clone):          ~45,000 gas
Reduction:              94.4%
```

## Development

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

## Security: Private Key Management

**IMPORTANT: Never commit private keys to git!**

### Key Files Structure

```
docker/
├── .env                   # Environment variables with private keys (gitignored)
└── besu/
    ├── entrypoint.sh      # Runtime key file generator
    ├── genesis.json
    └── static-nodes.json
```

### How It Works

1. **Private keys are stored in `docker/.env`** as environment variables:
   ```bash
   NODE1_PRIVATE_KEY=ac0974...
   NODE2_PRIVATE_KEY=38e49b...
   NODE3_PRIVATE_KEY=77dd7b...
   NODE4_PRIVATE_KEY=8897d0...
   ```

2. **Entrypoint script creates key files at runtime**:
   - Reads `$NODE_PRIVATE_KEY` from environment
   - Creates `/config/node-key` inside container
   - Sets 600 permissions
   - Besu uses the generated file

3. **Benefits**:
   - Single source of truth (`docker/.env`)
   - No sensitive files in repository
   - Easy to rotate keys
   - Production-ready (use Docker secrets in prod)

### Regenerating Validator Keys

If you need to change validator keys:

```bash
# 1. Generate new private key
cast wallet new

# 2. Get validator address
cast wallet address --private-key <YOUR_NEW_KEY>

# 3. Update docker/.env
NODE1_PRIVATE_KEY=<new_key>

# 4. Update genesis.json (extraData + alloc) and static-nodes.json (enode URLs)

# 5. Reset network
cd docker
docker compose down -v
docker compose up -d
```

**Production Deployment**: Use Docker secrets or vault instead of .env files.


# 3. Update docker/.env
NODE1_PRIVATE_KEY=<new_key>

# 4. Update genesis.json (extraData + alloc) and static-nodes.json (enode URLs)

# 5. Reset network
cd docker
docker compose down -v
docker compose up -d
```

**Production Deployment**: Use Docker secrets or vault instead of .env files.
## Environment Variables

Environment variables are organized by service. Each folder contains its own `.env` file.

### Docker Infrastructure (docker/.env)

**Besu Validator Keys:**
- `NODE1_PRIVATE_KEY` = Private key for Besu validator node 1
- `NODE2_PRIVATE_KEY` = Private key for Besu validator node 2
- `NODE3_PRIVATE_KEY` = Private key for Besu validator node 3
- `NODE4_PRIVATE_KEY` = Private key for Besu validator node 4

**Network Configuration:**
- `SUBNET` = Docker network subnet (default: 172.20.0.0/16)
- `NODE1_IP` ~ `NODE4_IP` = Fixed IP addresses for validator nodes
- `POSTGRES_IP` = PostgreSQL server IP address
- `REDIS_IP` = Redis server IP address

**Besu Network:**
- `BESU_RPC_URL` = HTTP RPC endpoint URL
- `BESU_WS_URL` = WebSocket endpoint URL
- `BESU_CHAIN_ID` = Network chain ID
- `BESU_NETWORK_ID` = Network ID

**Deployed Smart Contracts:**
- `VAULT_FACTORY_ADDRESS` = VaultFactory contract address
- `IMPLEMENTATION_ADDRESS` = Vault implementation contract address

**Database:**
- `POSTGRES_USER` = PostgreSQL username
- `POSTGRES_PASSWORD` = PostgreSQL password
- `POSTGRES_DB` = Database name
- `POSTGRES_PORT` = PostgreSQL port (default: 5432)

**Redis:**
- `REDIS_PORT` = Redis port (default: 6379)
- `REDIS_PASSWORD` = Redis authentication password

### Backend API (backend/.env)

**Server Configuration:**
- `PORT` = API server port (default: 8080)
- `ENV` = Environment mode (development/production/test)
- `LOG_LEVEL` = Logging level (debug/info/warn/error)

**Database:**
- `DB_HOST` = Database host address
- `DB_PORT` = Database port
- `DB_USER` = Database username
- `DB_PASSWORD` = Database password
- `DB_NAME` = Database name
- `DB_SSLMODE` = SSL mode (disable/require/verify-full)

**Redis:**
- `REDIS_HOST` = Redis host address
- `REDIS_PORT` = Redis port
- `REDIS_PASSWORD` = Redis password
- `REDIS_DB` = Redis database number (0-15)

**Blockchain:**
- `RPC_URL` = Ethereum RPC endpoint
- `WS_URL` = Ethereum WebSocket endpoint
- `CHAIN_ID` = Network chain ID
- `VAULT_FACTORY_ADDRESS` = VaultFactory contract address
- `PRIVATE_KEY` = Backend wallet private key for transactions

**JWT Authentication:**
- `JWT_SECRET` = Secret key for JWT token signing
- `JWT_EXPIRY` = Token expiration time (e.g., 24h)
- `JWT_REFRESH_SECRET` = Secret key for refresh token
- `JWT_REFRESH_EXPIRY` = Refresh token expiration time

**API Rate Limiting:**
- `RATE_LIMIT_MAX` = Maximum requests per window
- `RATE_LIMIT_WINDOW` = Time window for rate limiting (e.g., 1m)

**Email (SMTP):**
- `SMTP_HOST` = SMTP server address
- `SMTP_PORT` = SMTP server port
- `SMTP_USER` = SMTP username
- `SMTP_PASSWORD` = SMTP password
- `EMAIL_FROM` = Default sender email address

**DID/Identity Verification:**
- `NICE_CLIENT_ID` = NICE API client ID
- `NICE_SECRET` = NICE API secret key
- `PASS_CLIENT_ID` = PASS API client ID
- `PASS_SECRET` = PASS API secret key

**External APIs:**
- `INFURA_API_KEY` = Infura project API key
- `ALCHEMY_API_KEY` = Alchemy API key
- `ETHERSCAN_API_KEY` = Etherscan API key

**Monitoring:**
- `PROMETHEUS_PORT` = Prometheus metrics port
- `GRAFANA_PORT` = Grafana dashboard port

**Security:**
- `CORS_ORIGINS` = Allowed CORS origins (comma-separated)
- `TRUSTED_PROXIES` = Trusted proxy IP addresses
- `ENCRYPTION_KEY` = Data encryption key (32 bytes)

**Testing:**
- `TEST_DB_NAME` = Test database name
- `TEST_TIMEOUT` = Test execution timeout

**Feature Flags:**
- `FEATURE_EMAIL_VERIFICATION` = Enable email verification (true/false)
- `FEATURE_KYC` = Enable KYC verification (true/false)

### Smart Contracts (contracts/.env)

**Deployment:**
- `PRIVATE_KEY` = Deployer wallet private key
- `RPC_URL` = Target network RPC endpoint
- `ETHERSCAN_API_KEY` = Etherscan API key for verification
- `CHAIN_ID` = Target chain ID

**Verification:**
- `VERIFY_CONTRACT` = Auto-verify on Etherscan (true/false)

### Security Best Practices

1. **Never commit `.env` files** - All `.env` files are gitignored
2. **Use `.env.example` templates** - Each folder has example templates
3. **Rotate keys regularly** - Change private keys and secrets periodically
4. **Use strong secrets** - Generate random 32+ character secrets
5. **Production environments** - Use Docker secrets, AWS Secrets Manager, or HashiCorp Vault

### Quick Setup

```bash
# Copy example files
cp docker/.env.example docker/.env
cp backend/.env.example backend/.env
cp contracts/.env.example contracts/.env

# Edit each file with your values
nano docker/.env
nano backend/.env
nano contracts/.env
```
## �Documentation

- [PRD (Product Requirements Document)](docs/PRD.md)
- [Development Log](docs/DEV_LOG.md)
- [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
- [Security Report](docs/SECURITY_REPORT.md)
- [API Specification](docs/API_SPEC.md)

## Technology Stack

**Smart Contracts:**
- Solidity 0.8.33
- Foundry (Forge, Cast, Anvil)
- OpenZeppelin Contracts 5.5.0

**Blockchain:**
- Hyperledger Besu 24.12.0
- Clique PoA Consensus

**Backend:**
- Go 1.25.0
- Fiber v3 Framework
- GORM (PostgreSQL Driver)
- PostgreSQL 16
- Redis 7
- go-ethereum

**Frontend:**
- Next.js 14
- TypeScript
- Ethers.js

**Infrastructure:**
- Docker & Docker Compose
- GitHub Actions (CI/CD)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

## Contact

For questions or support, please open an issue on GitHub.

---

**Built with by LegacyChain Team**
