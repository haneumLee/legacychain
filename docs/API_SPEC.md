# LegacyChain - API Specification

> **Version**: v1.0  
> **Base URL**: `http://localhost:8080/api/v1`  
> ****: 2026 1 12

---

## 

1. [](#)
2. [ (Authentication)](#-authentication)
3. [Vault ](#vault-)
4. [Heartbeat](#heartbeat)
5. [Heir ()](#heir-)
6. [ ](#-)
7. [Examples](#examples)

---

## 

### Base URL
```
Production: https://api.legacychain.io/api/v1
Development: http://localhost:8080/api/v1
```

###  
- **JWT Bearer Token** ( Protected )
- Header: `Authorization: Bearer <token>`

###   

**:**
```json
{
  "data": { ... },
  "message": "Success"
}
```

**:**
```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

---

##  (Authentication)

### 1. Get Nonce

  Nonce .

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

**:**
- `nonce`: 1 UUID (5 )
- `message`: MetaMask  
- `timestamp`: Unix timestamp ( )

---

### 2. Login (Signature Verification)

Ethereum   JWT  .

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

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "a1b2c3d4-...",
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "created_at": "2026-01-12T10:30:00Z"
  }
}
```

**Validation:**
- Nonce Redis  1 
- Timestamp 5 
- Message Nonce + Timestamp 
- Signature EIP-191  

**:**
```json
{
  "error": "Invalid or expired nonce"
}
```

---

### 3. Get Current User

    .

**Endpoint:** `GET /auth/me`

**Request:**
```http
GET /api/v1/auth/me
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": "a1b2c3d4-...",
  "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "created_at": "2026-01-12T10:30:00Z",
  "updated_at": "2026-01-12T10:30:00Z"
}
```

---

## Vault 

### 1. Create Vault

 Vault .

**Endpoint:** `POST /vaults`

**Request:**
```http
POST /api/v1/vaults
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": 1,
  "contract_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "status": "active"
}
```

**Response:**
```json
{
  "id": "vault-uuid-...",
  "vault_id": 1,
  "owner_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "contract_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "status": "active",
  "created_at": "2026-01-12T10:30:00Z"
}
```

---

### 2. List Vaults

  Vault  .

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
    "id": "vault-uuid-1",
    "vault_id": 1,
    "owner_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "contract_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
    "status": "active",
    "created_at": "2026-01-12T10:30:00Z"
  }
]
```

---

### 3. Get Vault Details

 Vault   .

**Endpoint:** `GET /vaults/:id`

**Request:**
```http
GET /api/v1/vaults/vault-uuid-1
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": "vault-uuid-1",
  "vault_id": 1,
  "owner_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "contract_address": "0x5FbDB2315678afecb367f032d93F642f64180aa3",
  "status": "active",
  "created_at": "2026-01-12T10:30:00Z",
  "heirs": [
    {
      "id": "heir-uuid-1",
      "address": "0xHeir1Address...",
      "share_bps": 5000
    },
    {
      "id": "heir-uuid-2",
      "address": "0xHeir2Address...",
      "share_bps": 5000
    }
  ]
}
```

---

## Heartbeat

### 1. Commit Heartbeat

Heartbeat Commit Hash   (Commit-Reveal Step 1).

**Endpoint:** `POST /heartbeat/commit`

**Request:**
```http
POST /api/v1/heartbeat/commit
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": "vault-uuid-1",
  "nonce": "0x1234567890abcdef..."
}
```

**Parameters:**
- `vault_id` (string, required): Vault UUID
- `nonce` (string, required): Random hex string (client generates)

**Response:**
```json
{
  "tx_hash": "0xTransactionHash...",
  "commit_hash": "0xCommitHash...",
  "message": "Heartbeat committed successfully. Remember to reveal within the timeout period."
}
```

**Commit Hash  :**
```javascript
// Frontend
const nonce = ethers.randomBytes(32);
const commitHash = ethers.keccak256(
  ethers.concat([ownerAddress, nonce])
);
```

---

### 2. Reveal Heartbeat

 Commit Nonce  Heartbeat  (Commit-Reveal Step 2).

**Endpoint:** `POST /heartbeat/reveal`

**Request:**
```http
POST /api/v1/heartbeat/reveal
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": "vault-uuid-1",
  "nonce": "0x1234567890abcdef..."
}
```

**Response:**
```json
{
  "tx_hash": "0xRevealTxHash...",
  "message": "Heartbeat revealed successfully"
}
```

**Validation:**
- Nonce Database Committed  
- Blockchain CommitHash 
- Grace Period    (Owner )

---

### 3. Get Heartbeat Status

Vault  Heartbeat  .

**Endpoint:** `GET /heartbeat/status/:vault_id`

**Request:**
```http
GET /api/v1/heartbeat/status/vault-uuid-1
Authorization: Bearer <token>
```

**Response:**
```json
{
  "vault_id": "vault-uuid-1",
  "latest_commit": {
    "id": "heartbeat-uuid-1",
    "vault_id": "vault-uuid-1",
    "commit_hash": "0xCommitHash...",
    "commit_tx_hash": "0xCommitTx...",
    "reveal_tx_hash": "0xRevealTx...",
    "nonce": "0x1234...",
    "status": "revealed",
    "committed_at": "2026-01-12T10:30:00Z",
    "revealed_at": "2026-01-12T10:35:00Z"
  },
  "last_heartbeat_timestamp": "1673456789",
  "onchain_status": "Last heartbeat: 2026-01-12T10:35:00Z"
}
```

**Status Values:**
- `committed`: Commit , Reveal 
- `revealed`: Reveal 
- `failed`: Reveal 

---

### 4. List Heartbeats

Vault  Heartbeat  .

**Endpoint:** `GET /heartbeat/list/:vault_id`

**Request:**
```http
GET /api/v1/heartbeat/list/vault-uuid-1
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": "heartbeat-uuid-2",
    "vault_id": "vault-uuid-1",
    "commit_hash": "0x...",
    "commit_tx_hash": "0x...",
    "reveal_tx_hash": "0x...",
    "status": "revealed",
    "committed_at": "2026-01-12T10:30:00Z",
    "revealed_at": "2026-01-12T10:35:00Z"
  },
  {
    "id": "heartbeat-uuid-1",
    "vault_id": "vault-uuid-1",
    "commit_hash": "0x...",
    "commit_tx_hash": "0x...",
    "status": "committed",
    "committed_at": "2026-01-11T10:30:00Z",
    "revealed_at": null
  }
]
```

---

## Heir ()

### 1. Approve Inheritance

   (Multi-sig ).

**Endpoint:** `POST /heir/approve`

**Request:**
```http
POST /api/v1/heir/approve
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": "vault-uuid-1"
}
```

**Response:**
```json
{
  "tx_hash": "0xApproveTxHash...",
  "message": "Inheritance approved successfully"
}
```

**Requirements:**
- Caller Vault Heir 
- Vault Unlocked 
-    (On-chain)

---

### 2. Claim Inheritance

     .

**Endpoint:** `POST /heir/claim`

**Request:**
```http
POST /api/v1/heir/claim
Authorization: Bearer <token>
Content-Type: application/json

{
  "vault_id": "vault-uuid-1"
}
```

**Response:**
```json
{
  "tx_hash": "0xClaimTxHash...",
  "message": "Inheritance claimed successfully"
}
```

**Requirements:**
- Vault Unlocked 
- Required Approvals  ()
- Grace Period 
- Heir  Claim 

**:**
```json
{
  "error": "Insufficient approvals: 2/3 (need 2). You have approved."
}
```

---

### 3. Get Approval Status

   .

**Endpoint:** `GET /heir/status/:vault_id`

**Request:**
```http
GET /api/v1/heir/status/vault-uuid-1
Authorization: Bearer <token>
```

**Response:**
```json
{
  "vault_id": "vault-uuid-1",
  "heir_address": "0xYourAddress...",
  "approval_count": "2 (You have approved)",
  "required_count": 2,
  "can_claim": true
}
```

**Fields:**
- `approval_count`:    (   )
- `required_count`:    ()
- `can_claim`:   

---

### 4. List Heirs

Vault    .

**Endpoint:** `GET /heir/list/:vault_id`

**Request:**
```http
GET /api/v1/heir/list/vault-uuid-1
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": "heir-uuid-1",
    "vault_id": "vault-uuid-1",
    "address": "0xHeir1Address...",
    "share_bps": 5000,
    "created_at": "2026-01-12T10:30:00Z"
  },
  {
    "id": "heir-uuid-2",
    "vault_id": "vault-uuid-1",
    "address": "0xHeir2Address...",
    "share_bps": 5000,
    "created_at": "2026-01-12T10:30:00Z"
  }
]
```

**Note:** `share_bps` Basis Points (10000 = 100%)

---

##  

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Bad Request (  ) |
| 401 | Unauthorized (JWT  /) |
| 403 | Forbidden ( ) |
| 404 | Not Found ( ) |
| 500 | Internal Server Error |

###  

```json
{
  "error": "Invalid request body"
}
```

```json
{
  "error": "Invalid vault ID format"
}
```

```json
{
  "error": "Vault not found or you don't have permission"
}
```

```json
{
  "error": "You are not an heir of this vault"
}
```

---

## Examples

### Complete Authentication Flow

```javascript
// 1. Get Nonce
const nonceRes = await fetch('http://localhost:8080/api/v1/auth/nonce');
const { nonce, message, timestamp } = await nonceRes.json();

// 2. Sign with MetaMask
const signature = await ethereum.request({
  method: 'personal_sign',
  params: [message, accounts[0]]
});

// 3. Login
const loginRes = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    address: accounts[0],
    signature,
    message,
    nonce,
    timestamp
  })
});

const { token, user } = await loginRes.json();

// 4. Use Token
const vaultsRes = await fetch('http://localhost:8080/api/v1/vaults', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const vaults = await vaultsRes.json();
```

---

### Complete Heartbeat Flow

```javascript
// 1. Generate Nonce
const nonce = ethers.randomBytes(32);
const nonceHex = ethers.hexlify(nonce);

// 2. Commit
const commitRes = await fetch('http://localhost:8080/api/v1/heartbeat/commit', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    vault_id: 'vault-uuid-1',
    nonce: nonceHex
  })
});

const { tx_hash: commitTx, commit_hash } = await commitRes.json();

// 3. Wait for Commit Tx
await provider.waitForTransaction(commitTx);

// 4. Reveal (within timeout)
const revealRes = await fetch('http://localhost:8080/api/v1/heartbeat/reveal', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    vault_id: 'vault-uuid-1',
    nonce: nonceHex
  })
});

const { tx_hash: revealTx } = await revealRes.json();
await provider.waitForTransaction(revealTx);

console.log('Heartbeat completed!');
```

---

### Heir Approval and Claim Flow

```javascript
// 1. Check Status (as Heir)
const statusRes = await fetch('http://localhost:8080/api/v1/heir/status/vault-uuid-1', {
  headers: { 'Authorization': `Bearer ${heirToken}` }
});
const status = await statusRes.json();
// { approval_count: "1 (You have not approved)", required_count: 2, can_claim: false }

// 2. Approve Inheritance
const approveRes = await fetch('http://localhost:8080/api/v1/heir/approve', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${heirToken}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ vault_id: 'vault-uuid-1' })
});

const { tx_hash: approveTx } = await approveRes.json();
await provider.waitForTransaction(approveTx);

// 3. Re-check Status
const newStatus = await fetch('http://localhost:8080/api/v1/heir/status/vault-uuid-1', {
  headers: { 'Authorization': `Bearer ${heirToken}` }
}).then(r => r.json());
// { approval_count: "2 (You have approved)", required_count: 2, can_claim: true }

// 4. Claim Inheritance (if Grace Period ended)
const claimRes = await fetch('http://localhost:8080/api/v1/heir/claim', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${heirToken}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({ vault_id: 'vault-uuid-1' })
});

const { tx_hash: claimTx } = await claimRes.json();
await provider.waitForTransaction(claimTx);

console.log('Inheritance claimed!');
```

---

## Rate Limiting

- **Rate Limit**: 100 requests per minute per IP
- **Header**: `X-RateLimit-Remaining: 99`
- **429 Error**: Too Many Requests

```json
{
  "error": "Rate limit exceeded. Try again in 60 seconds."
}
```

---

## CORS

**Allowed Origins (Development):**
- `*` ( Origin )

**Allowed Headers:**
- `Origin`, `Content-Type`, `Accept`, `Authorization`

**Allowed Methods:**
- `GET`, `POST`, `PUT`, `DELETE`, `OPTIONS`

**Production:**    

---

## WebSocket Events (Planned)

 WebSocket  :

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.on('vault_created', (event) => {
  console.log('New vault:', event.vault_address);
});

ws.on('heartbeat_committed', (event) => {
  console.log('Heartbeat committed:', event.tx_hash);
});

ws.on('inheritance_claimed', (event) => {
  console.log('Claimed by:', event.heir_address);
});
```

---

**Last Updated**: 2026-01-12  
**Version**: v1.0  
**Maintained by**: LegacyChain Dev Team
