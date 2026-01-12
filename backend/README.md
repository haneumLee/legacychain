# LegacyChain Backend API

> Go + Fiber 기반 고성능 REST API 서버

## 📖 개요

LegacyChain의 Backend API는 Go 언어와 Fiber 프레임워크를 사용하여 구축된 고성능 RESTful API 서버입니다. PostgreSQL과 Redis를 사용하며, go-ethereum을 통해 스마트 컨트랙트와 통신합니다.

## 🏗️ 아키텍처

```
backend/
├── api/
│   ├── handlers/       # HTTP 요청 핸들러
│   │   ├── auth.go    # 인증 (Login, GetMe)
│   │   └── vault.go   # Vault CRUD
│   ├── middleware/     # 미들웨어
│   │   ├── auth.go    # JWT 인증
│   │   └── ratelimit.go # Rate Limiting
│   └── routes/         # 라우트 설정
├── models/             # GORM 모델
│   ├── user.go
│   ├── vault.go
│   ├── heir.go
│   └── heartbeat.go
├── services/           # 비즈니스 로직 (예정)
├── utils/              # 유틸리티 함수
│   ├── database.go    # DB 초기화
│   └── redis.go       # Redis 초기화
├── config/             # 설정 관리
│   └── config.go
├── cmd/                # 애플리케이션 진입점
│   └── main.go
├── .env.example        # 환경 변수 템플릿
└── go.mod              # Go 모듈 정의
```

## 🚀 Quick Start

### 1. 환경 변수 설정

```bash
cp .env.example .env
# .env 파일을 편집하여 설정 변경
```

### 2. 의존성 설치

```bash
go mod download
```

### 3. 빌드 및 실행

```bash
# 빌드
go build -o bin/server ./cmd/main.go

# 실행
./bin/server
```

또는 직접 실행:

```bash
go run ./cmd/main.go
```

서버는 기본적으로 `http://localhost:8080`에서 실행됩니다.

## 📡 API Endpoints

### Health Check

```
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "legacychain-backend"
}
```

### Authentication

#### Login
```
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "signature": "0x...",
  "message": "Login to LegacyChain"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "created_at": "2026-01-13T10:00:00Z"
  }
}
```

#### Get Current User
```
GET /api/v1/auth/me
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "email": "user@example.com",
  "nickname": "Alice",
  "created_at": "2026-01-13T10:00:00Z"
}
```

### Vaults

#### Create Vault
```
POST /api/v1/vaults
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "vault_id": 1,
  "contract_address": "0x1234...",
  "heartbeat_interval": 2592000,
  "grace_period": 604800,
  "required_approvals": 2,
  "heir_addresses": [
    "0xHeir1...",
    "0xHeir2...",
    "0xHeir3..."
  ],
  "heir_shares": [5000, 3000, 2000]
}
```

**Response:**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "vault_id": 1,
  "contract_address": "0x1234...",
  "owner_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "locked",
  "heartbeat_interval": 2592000,
  "grace_period": 604800,
  "required_approvals": 2,
  "heirs": [...],
  "created_at": "2026-01-13T10:30:00Z"
}
```

#### List Vaults
```
GET /api/v1/vaults
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "vault_id": 1,
    "contract_address": "0x1234...",
    "status": "locked",
    "heirs": [...]
  }
]
```

#### Get Vault
```
GET /api/v1/vaults/:id
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "vault_id": 1,
  "contract_address": "0x1234...",
  "owner": {...},
  "heirs": [...],
  "heartbeats": [...],
  "status": "locked"
}
```

## 🔐 Authentication

이 API는 JWT (JSON Web Token) 기반 인증을 사용합니다.

1. `/api/v1/auth/login`으로 Ethereum 서명 검증 후 JWT 발급
2. 이후 모든 요청의 `Authorization` 헤더에 `Bearer <token>` 포함
3. 토큰은 기본 24시간 유효 (`.env`에서 변경 가능)

## 📊 Rate Limiting

Redis 기반 Rate Limiting 적용:

- 기본: IP당 100 requests/minute
- 응답 헤더:
  - `X-RateLimit-Limit`: 최대 요청 수
  - `X-RateLimit-Remaining`: 남은 요청 수
  - `X-RateLimit-Reset`: 리셋 시간 (Unix timestamp)

초과 시 `429 Too Many Requests` 응답

## 🗄️ Database Models

### User
- `id` (UUID, PK)
- `address` (Ethereum address, unique)
- `email`, `nickname` (optional)
- Soft Delete 지원

### Vault
- `id` (UUID, PK)
- `vault_id` (int, unique, on-chain ID)
- `contract_address` (Ethereum address, unique)
- `owner_id` (FK → User)
- `status` (locked, unlocked, claimed)
- `heartbeat_interval`, `grace_period`
- `required_approvals`

### Heir
- `id` (UUID, PK)
- `vault_id` (FK → Vault)
- `address` (Ethereum address)
- `share_bps` (Basis Points: 0-10000)
- `has_approved`, `has_claimed` (boolean)

### Heartbeat
- `id` (UUID, PK)
- `vault_id` (FK → Vault)
- `tx_hash` (unique, on-chain transaction)
- `timestamp`

## 🛠️ Technology Stack

- **Language**: Go 1.25.0
- **Framework**: Fiber v3.0.0-rc.3
- **ORM**: GORM v1.31.1
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Blockchain**: go-ethereum v1.16.7
- **Auth**: JWT v5.3.0

## 📦 Dependencies

주요 의존성:

```go
github.com/gofiber/fiber/v3          // Web framework
gorm.io/gorm                          // ORM
gorm.io/driver/postgres               // PostgreSQL driver
github.com/redis/go-redis/v9          // Redis client
github.com/ethereum/go-ethereum       // Ethereum client
github.com/golang-jwt/jwt/v5          // JWT auth
github.com/google/uuid                // UUID generation
github.com/joho/godotenv              // .env support
```

## 🔧 Development

### 코드 포맷팅
```bash
go fmt ./...
```

### Linting
```bash
go vet ./...
```

### 테스트 (예정)
```bash
go test ./...
```

## 🚧 TODO (Day 13-15)

- [ ] Ethereum 서명 검증 구현 (ECDSA Personal Sign)
- [ ] Blockchain Service 구현
  - [ ] go-ethereum 클라이언트 설정
  - [ ] VaultFactory ABI 바인딩
  - [ ] 이벤트 리스닝 (VaultCreated, HeartbeatCommitted)
- [ ] Heartbeat Handlers (Commit, Reveal, Status)
- [ ] Heir Handlers (Approve, Claim)
- [ ] Unit Tests
- [ ] Integration Tests
- [ ] Swagger/OpenAPI 문서

## 📄 Environment Variables

`.env.example` 참고:

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=legacychain
DB_PASSWORD=legacychain_password
DB_NAME=legacychain

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=legacychain_redis_password

# Blockchain
BESU_RPC_URL=http://localhost:8545
BESU_WS_URL=ws://localhost:8546
CHAIN_ID=1337
VAULT_FACTORY_ADDRESS=0x5FbDB2315678afecb367f032d93F642f64180aa3

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRES_IN=24h

# Rate Limiting
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=1m
```

## 📚 References

- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [go-ethereum Documentation](https://geth.ethereum.org/docs/developers/dapp-developer/native)

---

**Last Updated**: 2026-01-13
