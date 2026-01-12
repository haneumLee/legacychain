# LegacyChain - API Specification

> **Version**: v1.0  
> **Base URL**: `http://localhost:8080/api/v1`  
> **Last Updated**: 2026-01-12

---

## 목차 (Table of Contents)

1. [개요 (Overview)](#개요-overview)
2. [인증 (Authentication)](#인증-authentication)
3. [Vault 관리 (Vault Management)](#vault-관리-vault-management)
4. [Heartbeat](#heartbeat)
5. [상속인 (Heir Management)](#상속인-heir-management)
6. [에러 코드 (Error Codes)](#에러-코드-error-codes)
7. [Examples](#examples)

---

## 개요 (Overview)

### Base URL
```
Production: https://api.legacychain.io/api/v1
Development: http://localhost:8080/api/v1
```

### 인증 방식 (Authentication Method)
- **JWT Bearer Token** (Protected 엔드포인트)
- Header: `Authorization: Bearer <token>`

### 응답 형식 (Response Format)

**성공 응답:**
```json
{
  "data": { ... },
  "message": "Success"
}
```

**에러 응답:**
```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

---

## 인증 (Authentication)

### 1. Get Nonce

로그인을 위한 Nonce를 생성합니다.

**Endpoint:** `GET /auth/nonce`

**Request:**
```http
GET /api/v1/auth/nonce
```

**Response:**
```json
{
  "nonce": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Login to LegacyChain\nNonce: 550e8400-e29b-41d4-a716-446655440000\nTimestamp: 1673456789",
  "timestamp": 1673456789
}
```

**응답 필드:**
- `nonce`: 일회성 UUID (5분 유효)
- `message`: MetaMask 서명용 메시지
- `timestamp`: Unix timestamp (nonce 생성 시각)

---

### 2. Login (Signature Verification)

Ethereum 서명을 검증하여 JWT 토큰을 발급합니다.

**Endpoint:** `POST /auth/login`

**Request:**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "signature": "0x1234...abcd",
  "message": "Login to LegacyChain\nNonce: 550e8400...\nTimestamp: 1673456789",
  "nonce": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": 1673456789
}
```

**Request Fields:**
- `address` (required): Ethereum 지갑 주소
- `signature` (required): MetaMask 서명 (EIP-191)
- `message` (required): 서명된 메시지 (Get Nonce에서 받은 값)
- `nonce` (required): Get Nonce에서 받은 nonce
- `timestamp` (required): Get Nonce에서 받은 timestamp

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "created_at": "2026-01-12T10:30:00Z"
  }
}
```

**Response Fields:**
- `token`: JWT access token (24시간 유효)
- `user`: 사용자 정보 객체

**Errors:**
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid nonce, signature, or expired timestamp

---

### 3. Get Current User

현재 인증된 사용자 정보를 조회합니다.

**Endpoint:** `GET /auth/me`

**Request:**
```http
GET /api/v1/auth/me
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "created_at": "2026-01-12T10:30:00Z"
}
```

**Errors:**
- `401 Unauthorized`: Missing or invalid token

---

## Vault 관리 (Vault Management)

### 1. Create Vault

새로운 Vault를 데이터베이스에 등록합니다 (스마트 컨트랙트 배포 후 호출).

**Endpoint:** `POST /vaults`

**Request:**
```http
POST /api/v1/vaults
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": 1,
  "contract_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "heartbeat_interval": 2592000,
  "grace_period": 2592000,
  "required_approvals": 2,
  "heir_addresses": [
    "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
    "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
  ],
  "heir_shares": [5000, 5000]
}
```

**Request Fields:**
- `vault_id` (required): 온체인 Vault ID (BIGINT)
- `contract_address` (required): Vault 컨트랙트 주소
- `heartbeat_interval` (required): Heartbeat 주기 (초)
- `grace_period` (required): Grace Period (초)
- `required_approvals` (required): 필요한 승인 수
- `heir_addresses` (required): 상속인 주소 배열
- `heir_shares` (required): 상속인 지분 배열 (BPS: 5000 = 50%)

**Validation:**
- `heir_addresses`와 `heir_shares` 배열 길이 동일
- `heir_shares` 합계 = 10000 (100%)
- `heartbeat_interval` >= 259200 (3일)
- `required_approvals` <= 상속인 수

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "vault_id": 1,
  "contract_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "owner_id": "550e8400-e29b-41d4-a716-446655440001",
  "heartbeat_interval": 2592000,
  "grace_period": 2592000,
  "required_approvals": 2,
  "status": "locked",
  "created_at": "2026-01-12T11:00:00Z",
  "owner": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  },
  "heirs": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "vault_id": "550e8400-e29b-41d4-a716-446655440002",
      "address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
      "share_bps": 5000
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440004",
      "vault_id": "550e8400-e29b-41d4-a716-446655440002",
      "address": "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
      "share_bps": 5000
    }
  ]
}
```

**Errors:**
- `400 Bad Request`: Invalid request body or validation error
- `404 Not Found`: User not found
- `500 Internal Server Error`: Database error

---

### 2. Get Vault

Vault 상세 정보를 조회합니다.

**Endpoint:** `GET /vaults/:id`

**Request:**
```http
GET /api/v1/vaults/550e8400-e29b-41d4-a716-446655440002
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Vault UUID

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "vault_id": 1,
  "contract_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "owner_id": "550e8400-e29b-41d4-a716-446655440001",
  "heartbeat_interval": 2592000,
  "grace_period": 2592000,
  "required_approvals": 2,
  "status": "locked",
  "created_at": "2026-01-12T11:00:00Z",
  "owner": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  },
  "heirs": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
      "share_bps": 5000
    }
  ],
  "heartbeats": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440005",
      "vault_id": "550e8400-e29b-41d4-a716-446655440002",
      "commit_hash": "0xabc123...",
      "commit_tx_hash": "0xdef456...",
      "reveal_tx_hash": "0xghi789...",
      "status": "revealed",
      "committed_at": "2026-01-12T12:00:00Z",
      "revealed_at": "2026-01-12T12:05:00Z"
    }
  ]
}
```

**Errors:**
- `400 Bad Request`: Invalid vault ID format
- `404 Not Found`: Vault not found

---

### 3. List Vaults

인증된 사용자의 모든 Vault를 조회합니다.

**Endpoint:** `GET /vaults`

**Request:**
```http
GET /api/v1/vaults
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "vault_id": 1,
    "contract_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
    "owner_id": "550e8400-e29b-41d4-a716-446655440001",
    "heartbeat_interval": 2592000,
    "grace_period": 2592000,
    "required_approvals": 2,
    "status": "locked",
    "created_at": "2026-01-12T11:00:00Z",
    "heirs": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440003",
        "address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
        "share_bps": 5000
      }
    ]
  }
]
```

**Errors:**
- `404 Not Found`: User not found
- `500 Internal Server Error`: Database error

---

## Heartbeat

### 1. Commit Heartbeat

Heartbeat을 커밋합니다 (Commit-Reveal 패턴의 Phase 1).

**Endpoint:** `POST /heartbeat/commit`

**Request:**
```http
POST /api/v1/heartbeat/commit
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": "550e8400-e29b-41d4-a716-446655440002",
  "nonce": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
}
```

**Request Fields:**
- `vault_id` (required): Vault UUID
- `nonce` (required): 랜덤 값 (32 bytes hex string)

**Response:**
```json
{
  "tx_hash": "0xabc123def456...",
  "commit_hash": "0x789ghi012jkl...",
  "message": "Heartbeat committed successfully. Remember to reveal within the timeout period."
}
```

**Response Fields:**
- `tx_hash`: 블록체인 트랜잭션 해시
- `commit_hash`: keccak256(address + nonce)
- `message`: 안내 메시지

**Errors:**
- `400 Bad Request`: Invalid request body or vault ID
- `404 Not Found`: Vault not found or permission denied
- `500 Internal Server Error`: Blockchain or database error

---

### 2. Reveal Heartbeat

커밋된 Heartbeat을 공개합니다 (Commit-Reveal 패턴의 Phase 2).

**Endpoint:** `POST /heartbeat/reveal`

**Request:**
```http
POST /api/v1/heartbeat/reveal
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": "550e8400-e29b-41d4-a716-446655440002",
  "nonce": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
}
```

**Request Fields:**
- `vault_id` (required): Vault UUID
- `nonce` (required): Commit 시 사용한 nonce

**Response:**
```json
{
  "tx_hash": "0xdef456ghi789...",
  "message": "Heartbeat revealed successfully. Vault is now refreshed."
}
```

**Errors:**
- `400 Bad Request`: Invalid request body, vault ID, or nonce format
- `404 Not Found`: Vault not found or no pending commit
- `500 Internal Server Error`: Blockchain or database error

---

### 3. Get Heartbeat Status

Vault의 Heartbeat 상태를 조회합니다.

**Endpoint:** `GET /heartbeat/status/:vault_id`

**Request:**
```http
GET /api/v1/heartbeat/status/550e8400-e29b-41d4-a716-446655440002
Authorization: Bearer <token>
```

**Path Parameters:**
- `vault_id`: Vault UUID

**Response:**
```json
{
  "vault_id": "550e8400-e29b-41d4-a716-446655440002",
  "latest_commit": {
    "id": "550e8400-e29b-41d4-a716-446655440005",
    "commit_hash": "0x789ghi012jkl...",
    "status": "revealed",
    "committed_at": "2026-01-12T12:00:00Z",
    "revealed_at": "2026-01-12T12:05:00Z"
  },
  "last_heartbeat_timestamp": "1673520000",
  "onchain_status": "locked"
}
```

**Errors:**
- `400 Bad Request`: Invalid vault ID format
- `404 Not Found`: Vault not found
- `500 Internal Server Error`: Blockchain or database error

---

### 4. List Heartbeats

Vault의 모든 Heartbeat 기록을 조회합니다.

**Endpoint:** `GET /heartbeat/list/:vault_id`

**Request:**
```http
GET /api/v1/heartbeat/list/550e8400-e29b-41d4-a716-446655440002
Authorization: Bearer <token>
```

**Path Parameters:**
- `vault_id`: Vault UUID

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440006",
    "vault_id": "550e8400-e29b-41d4-a716-446655440002",
    "commit_hash": "0x789ghi012jkl...",
    "commit_tx_hash": "0xabc123def456...",
    "reveal_tx_hash": "0xdef456ghi789...",
    "status": "revealed",
    "committed_at": "2026-01-12T12:00:00Z",
    "revealed_at": "2026-01-12T12:05:00Z"
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440005",
    "vault_id": "550e8400-e29b-41d4-a716-446655440002",
    "commit_hash": "0x123abc456def...",
    "commit_tx_hash": "0x789ghi012jkl...",
    "reveal_tx_hash": "0xmno345pqr678...",
    "status": "revealed",
    "committed_at": "2026-01-11T12:00:00Z",
    "revealed_at": "2026-01-11T12:03:00Z"
  }
]
```

**Errors:**
- `400 Bad Request`: Invalid vault ID format
- `404 Not Found`: Vault not found
- `500 Internal Server Error`: Database error

---

## 상속인 (Heir Management)

### 1. Approve Inheritance

상속인으로서 상속 승인을 등록합니다 (Multi-sig 승인).

**Endpoint:** `POST /heir/approve`

**Request:**
```http
POST /api/v1/heir/approve
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

**Request Fields:**
- `vault_id` (required): Vault UUID

**Response:**
```json
{
  "tx_hash": "0xabc123def456...",
  "message": "Inheritance approved successfully"
}
```

**Errors:**
- `400 Bad Request`: Invalid request body or vault ID
- `403 Forbidden`: You are not an heir of this vault
- `404 Not Found`: Vault not found
- `500 Internal Server Error`: Blockchain or database error

---

### 2. Claim Inheritance

상속받을 자산을 청구합니다 (충분한 승인 필요).

**Endpoint:** `POST /heir/claim`

**Request:**
```http
POST /api/v1/heir/claim
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

**Request Fields:**
- `vault_id` (required): Vault UUID

**Response:**
```json
{
  "tx_hash": "0xdef456ghi789...",
  "message": "Inheritance claimed successfully"
}
```

**Errors:**
- `400 Bad Request`: Invalid request body or vault ID
- `403 Forbidden`: You are not an heir of this vault
- `404 Not Found`: Vault not found
- `500 Internal Server Error`: Blockchain or database error

---

### 3. Get Approval Status

Vault의 승인 상태를 조회합니다.

**Endpoint:** `GET /heir/status/:vault_id`

**Request:**
```http
GET /api/v1/heir/status/550e8400-e29b-41d4-a716-446655440002
Authorization: Bearer <token>
```

**Path Parameters:**
- `vault_id`: Vault UUID

**Response:**
```json
{
  "vault_id": "550e8400-e29b-41d4-a716-446655440002",
  "heir_address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
  "approval_count": "2",
  "required_count": 2,
  "can_claim": true
}
```

**Response Fields:**
- `vault_id`: Vault UUID
- `heir_address`: 현재 사용자의 주소
- `approval_count`: 현재 승인 수
- `required_count`: 필요한 승인 수
- `can_claim`: 청구 가능 여부

**Errors:**
- `400 Bad Request`: Invalid vault ID format
- `403 Forbidden`: You are not an heir of this vault
- `404 Not Found`: Vault not found
- `500 Internal Server Error`: Blockchain or database error

---

### 4. List Heirs

Vault의 모든 상속인을 조회합니다.

**Endpoint:** `GET /heir/list/:vault_id`

**Request:**
```http
GET /api/v1/heir/list/550e8400-e29b-41d4-a716-446655440002
Authorization: Bearer <token>
```

**Path Parameters:**
- `vault_id`: Vault UUID

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "vault_id": "550e8400-e29b-41d4-a716-446655440002",
    "address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
    "share_bps": 5000,
    "created_at": "2026-01-12T11:00:00Z"
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440004",
    "vault_id": "550e8400-e29b-41d4-a716-446655440002",
    "address": "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
    "share_bps": 5000,
    "created_at": "2026-01-12T11:00:00Z"
  }
]
```

**Errors:**
- `400 Bad Request`: Invalid vault ID format
- `404 Not Found`: Vault not found
- `500 Internal Server Error`: Database error

---

## 에러 코드 (Error Codes)

| HTTP Status | Error Code | Description |
|-------------|-----------|-------------|
| 400 | `BAD_REQUEST` | Invalid request body or parameters |
| 401 | `UNAUTHORIZED` | Missing or invalid authentication token |
| 403 | `FORBIDDEN` | Insufficient permissions |
| 404 | `NOT_FOUND` | Resource not found |
| 429 | `RATE_LIMIT_EXCEEDED` | Too many requests |
| 500 | `INTERNAL_SERVER_ERROR` | Server error |

---

## Examples

### Complete Workflow Example

#### 1. Login Flow
```bash
# Step 1: Get nonce
curl -X GET http://localhost:8080/api/v1/auth/nonce

# Response:
{
  "nonce": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Login to LegacyChain\nNonce: 550e8400...\nTimestamp: 1673456789",
  "timestamp": 1673456789
}

# Step 2: Sign message with MetaMask (client-side)
# const signature = await ethereum.request({
#   method: 'personal_sign',
#   params: [message, address]
# });

# Step 3: Login with signature
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "signature": "0x1234...abcd",
    "message": "Login to LegacyChain\nNonce: 550e8400...\nTimestamp: 1673456789",
    "nonce": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": 1673456789
  }'

# Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { ... }
}
```

#### 2. Create Vault
```bash
curl -X POST http://localhost:8080/api/v1/vaults \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "vault_id": 1,
    "contract_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
    "heartbeat_interval": 2592000,
    "grace_period": 2592000,
    "required_approvals": 2,
    "heir_addresses": [
      "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
      "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
    ],
    "heir_shares": [5000, 5000]
  }'
```

#### 3. Heartbeat Flow (Commit-Reveal)
```bash
# Step 1: Generate random nonce (client-side)
# const nonce = ethers.hexlify(ethers.randomBytes(32));

# Step 2: Commit heartbeat
curl -X POST http://localhost:8080/api/v1/heartbeat/commit \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "vault_id": "550e8400-e29b-41d4-a716-446655440002",
    "nonce": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
  }'

# Response:
{
  "tx_hash": "0xabc123...",
  "commit_hash": "0x789ghi...",
  "message": "Heartbeat committed successfully..."
}

# Step 3: Reveal heartbeat (after commit transaction is mined)
curl -X POST http://localhost:8080/api/v1/heartbeat/reveal \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "vault_id": "550e8400-e29b-41d4-a716-446655440002",
    "nonce": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
  }'

# Response:
{
  "tx_hash": "0xdef456...",
  "message": "Heartbeat revealed successfully..."
}
```

#### 4. Inheritance Approval Flow
```bash
# Step 1: Heir A approves
curl -X POST http://localhost:8080/api/v1/heir/approve \
  -H "Authorization: Bearer <heir_a_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "vault_id": "550e8400-e29b-41d4-a716-446655440002"
  }'

# Step 2: Check approval status
curl -X GET http://localhost:8080/api/v1/heir/status/550e8400-e29b-41d4-a716-446655440002 \
  -H "Authorization: Bearer <heir_a_token>"

# Response:
{
  "vault_id": "550e8400-e29b-41d4-a716-446655440002",
  "heir_address": "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
  "approval_count": "1",
  "required_count": 2,
  "can_claim": false
}

# Step 3: Heir B approves (reaching required threshold)
curl -X POST http://localhost:8080/api/v1/heir/approve \
  -H "Authorization: Bearer <heir_b_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "vault_id": "550e8400-e29b-41d4-a716-446655440002"
  }'

# Step 4: Claim inheritance
curl -X POST http://localhost:8080/api/v1/heir/claim \
  -H "Authorization: Bearer <heir_a_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "vault_id": "550e8400-e29b-41d4-a716-446655440002"
  }'

# Response:
{
  "tx_hash": "0xghi789...",
  "message": "Inheritance claimed successfully"
}
```

---

**Last Updated**: 2026-01-12
