# LegacyChain Backend API

> Go + Fiber   REST API 

## 

LegacyChain Backend API Go  Fiber     RESTful API . PostgreSQL Redis , go-ethereum    .

## 

```
backend/
 api/
    handlers/       # HTTP  
       auth.go    #  (Login, GetMe)
       vault.go   # Vault CRUD
    middleware/     # 
       auth.go    # JWT 
       ratelimit.go # Rate Limiting
    routes/         #  
 models/             # GORM 
    user.go
    vault.go
    heir.go
    heartbeat.go
 services/           #   ()
 utils/              #  
    database.go    # DB 
    redis.go       # Redis 
 config/             #  
    config.go
 cmd/                #  
    main.go
 .env.example        #   
 go.mod              # Go  
```

## Quick Start

### 1.   

```bash
cp .env.example .env
# .env    
```

### 2.  

```bash
go mod download
```

### 3.   

```bash
# 
go build -o bin/server ./cmd/main.go

# 
./bin/server
```

  :

```bash
go run ./cmd/main.go
```

  `http://localhost:8080` .

##  API Endpoints

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

## Authentication

 API JWT (JSON Web Token)   .

1. `/api/v1/auth/login` Ethereum    JWT 
2.    `Authorization`  `Bearer <token>` 
3.   24  (`.env`  )

## Rate Limiting

Redis  Rate Limiting :

- : IP 100 requests/minute
-  :
  - `X-RateLimit-Limit`:   
  - `X-RateLimit-Remaining`:   
  - `X-RateLimit-Reset`:   (Unix timestamp)

  `429 Too Many Requests` 

##  Database Models

### User
- `id` (UUID, PK)
- `address` (Ethereum address, unique)
- `email`, `nickname` (optional)
- Soft Delete 

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

## Technology Stack

- **Language**: Go 1.25.0
- **Framework**: Fiber v3.0.0-rc.3
- **ORM**: GORM v1.31.1
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Blockchain**: go-ethereum v1.16.7
- **Auth**: JWT v5.3.0

##  Dependencies

 :

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

## Development

###  
```bash
go fmt ./...
```

### Linting
```bash
go vet ./...
```

###  ()
```bash
go test ./...
```

##  TODO (Day 13-15)

- [ ] Ethereum    (ECDSA Personal Sign)
- [ ] Blockchain Service 
  - [ ] go-ethereum  
  - [ ] VaultFactory ABI 
  - [ ]   (VaultCreated, HeartbeatCommitted)
- [ ] Heartbeat Handlers (Commit, Reveal, Status)
- [ ] Heir Handlers (Approve, Claim)
- [ ] Unit Tests
- [ ] Integration Tests
- [ ] Swagger/OpenAPI 

## Environment Variables

`.env.example` :

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

## References

- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [go-ethereum Documentation](https://geth.ethereum.org/docs/developers/dapp-developer/native)

---

**Last Updated**: 2026-01-13
