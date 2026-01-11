# LegacyChain - Development Log

> **목적**: 개발 과정, 의사결정, 이슈 해결 기록  
> **시작일**: 2026년 1월 12일

---

## 📋 목차

1. [Phase 0: 개발 환경 구축](#phase-0-개발-환경-구축)
2. [Phase 1: Smart Contract 개발](#phase-1-smart-contract-개발)
3. [Backend 개발](#backend-개발)
4. [Frontend 개발](#frontend-개발)
5. [통합 및 배포](#통합-및-배포)

---

## Phase 0: 개발 환경 구축

### [2026-01-12] Day 0: 개발 도구 설치 및 프로젝트 초기화

#### 작업 내용
프로젝트 개발에 필요한 모든 도구 설치 및 디렉토리 구조 생성

#### 1. 개발 도구 설치

##### 1.1 초기 상태 확인
```bash
# 설치된 도구
✅ Node.js: v18.19.1
✅ npm: 9.2.0
✅ Docker: 29.1.4

# 미설치 도구
❌ Go (필수)
❌ Foundry (필수)
```

##### 1.2 Go 1.21.13 설치
```bash
# 설치 과정
cd /tmp
wget -q https://go.dev/dl/go1.21.13.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.13.linux-amd64.tar.gz
rm go1.21.13.linux-amd64.tar.gz

# PATH 설정
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin

# 확인
go version
# Output: go version go1.21.13 linux/amd64
```

**설치 이유**: Backend API 서버 개발에 Go 1.21+ 필수 (PRD 명세)

##### 1.3 Foundry 설치
```bash
# Foundryup 설치
curl -L https://foundry.paradigm.xyz | bash

# PATH 설정
export PATH="$HOME/.foundry/bin:$PATH"
echo 'export PATH="$HOME/.foundry/bin:$PATH"' >> ~/.bashrc

# Foundry 설치
source ~/.bashrc
foundryup

# 확인
forge --version
# Output: forge Version: 1.5.1-stable (b0a9dd9ced 2025-12-22)
cast --version
# Output: cast Version: 1.5.1-stable
anvil --version
# Output: anvil Version: 1.5.1-stable
```

**설치 이유**: Smart Contract 개발, 테스트, 배포에 Foundry 필수 (Hardhat 대비 빠른 테스트 속도)

##### 1.4 최종 설치 도구 버전
```
✅ Node.js: v18.19.1        (Frontend - Next.js 14)
✅ npm: 9.2.0               (패키지 관리)
✅ Go: go1.21.13            (Backend API)
✅ Foundry - forge: 1.5.1   (Smart Contract 개발)
✅ Foundry - cast: 1.5.1    (Blockchain 상호작용)
✅ Foundry - anvil: 1.5.1   (로컬 테스트 노드)
✅ Docker: 29.1.4           (Besu, PostgreSQL, Redis)
```

#### 2. 프로젝트 디렉토리 구조 생성

##### 2.1 디렉토리 생성 스크립트
```bash
cd /root/legacychain

# Smart Contract
mkdir -p contracts/{src,test/{unit,invariant,integration},script,lib}

# Backend
mkdir -p backend/{cmd/server,internal/{handler,service,repository,model,blockchain,middleware},pkg/{logger,validator,crypto},config,migrations}

# Frontend
mkdir -p frontend/{app/{vault,dashboard,did},components/{ui,vault,dashboard,layout},lib/{hooks,utils,contracts},public/{images,icons}}

# Infrastructure
mkdir -p docker/{besu,postgres,redis}
mkdir -p infrastructure/{aws,k8s,terraform}
mkdir -p scripts
```

##### 2.2 최종 프로젝트 구조
```
legacychain/
├── contracts/              # Smart Contract (Solidity)
│   ├── src/               # 컨트랙트 소스
│   ├── test/              # 테스트
│   │   ├── unit/          # 단위 테스트
│   │   ├── invariant/     # 속성 기반 테스트
│   │   └── integration/   # 통합 테스트
│   ├── script/            # 배포 스크립트
│   └── lib/               # 라이브러리 (forge-std 등)
│
├── backend/               # Go API Server
│   ├── cmd/server/        # 메인 엔트리포인트
│   ├── internal/          # 내부 패키지
│   │   ├── handler/       # HTTP 핸들러
│   │   ├── service/       # 비즈니스 로직
│   │   ├── repository/    # DB 액세스
│   │   ├── model/         # 데이터 모델
│   │   ├── blockchain/    # 블록체인 연동
│   │   └── middleware/    # 미들웨어 (Auth, CORS 등)
│   ├── pkg/               # 공개 패키지
│   │   ├── logger/        # 로깅
│   │   ├── validator/     # 검증
│   │   └── crypto/        # 암호화
│   ├── config/            # 설정 파일
│   └── migrations/        # DB 마이그레이션
│
├── frontend/              # Next.js 14 App
│   ├── app/               # App Router
│   │   ├── vault/         # Vault 관리
│   │   ├── dashboard/     # 대시보드
│   │   └── did/           # DID 관리
│   ├── components/        # React 컴포넌트
│   │   ├── ui/            # shadcn/ui
│   │   ├── vault/         # Vault 관련
│   │   ├── dashboard/     # Dashboard 관련
│   │   └── layout/        # 레이아웃
│   └── lib/               # 유틸리티
│       ├── hooks/         # Custom Hooks
│       ├── utils/         # 헬퍼 함수
│       └── contracts/     # Contract ABI/주소
│
├── docker/                # Docker 설정
│   ├── besu/              # Besu 노드
│   ├── postgres/          # PostgreSQL
│   └── redis/             # Redis
│
├── infrastructure/        # IaC
│   ├── aws/               # AWS 리소스
│   ├── k8s/               # Kubernetes
│   └── terraform/         # Terraform
│
├── scripts/               # 자동화 스크립트
└── docs/                  # 문서
```

#### 3. 프로젝트 초기화

##### 3.1 Foundry 프로젝트 초기화
```bash
cd /root/legacychain/contracts
forge init --force --no-git

# 설치된 라이브러리
✅ forge-std (Foundry 표준 라이브러리)
```

**결과**: 
- `foundry.toml` 생성 (Foundry 설정)
- `lib/forge-std` 설치 (테스트 유틸리티)
- 샘플 컨트랙트 생성 (나중에 제거 예정)

##### 3.2 Go 모듈 초기화
```bash
cd /root/legacychain/backend
go mod init github.com/haneumLee/legacychain/backend

# 결과
✅ go.mod 생성
```

**설정**:
- 모듈 경로: `github.com/haneumLee/legacychain/backend`
- Go 버전: 1.21

##### 3.3 Next.js 14 프로젝트 생성
```bash
cd /root/legacychain
rm -rf frontend  # 빈 디렉토리 제거
mkdir frontend
cd frontend

npx -y create-next-app@14 . \
  --typescript \
  --tailwind \
  --app \
  --no-src-dir \
  --import-alias "@/*" \
  --skip-install
```

**설정**:
- ✅ TypeScript
- ✅ Tailwind CSS
- ✅ App Router (Next.js 14)
- ✅ ESLint
- ✅ Import alias: `@/*`

**설치된 패키지**:
```json
{
  "dependencies": {
    "react": "^18",
    "react-dom": "^18",
    "next": "14.2.35"
  },
  "devDependencies": {
    "@types/node": "^20",
    "@types/react": "^18",
    "@types/react-dom": "^18",
    "eslint": "^8",
    "eslint-config-next": "14.2.35",
    "postcss": "^8",
    "tailwindcss": "^3.4.1",
    "typescript": "^5"
  }
}
```

#### 4. 이슈 및 해결

##### Issue 1: Next.js 프로젝트 충돌
**문제**: 수동으로 생성한 `frontend/{app,components,lib}` 디렉토리와 create-next-app이 충돌

**에러 메시지**:
```
The directory frontend contains files that could conflict:
  app/
  components/
  lib/
  public/
```

**해결**:
```bash
# 디렉토리 완전 재생성
rm -rf frontend
mkdir frontend
cd frontend
npx -y create-next-app@14 . --typescript --tailwind --app
```

**교훈**: create-next-app은 빈 디렉토리에서 실행해야 함

##### Issue 2: npm 보안 취약점 경고
**경고**:
```
3 high severity vulnerabilities
To address all issues (including breaking changes), run:
  npm audit fix --force
```

**대응**: 
- 현재는 무시 (개발 초기)
- 배포 전 `npm audit fix` 실행 예정
- TROUBLESHOOTING.md에 기록

#### 5. 다음 단계 (Day 3-6)

##### 준비 완료 체크리스트
```yaml
✅ 개발 도구 설치
  ├─ Go 1.21.13
  ├─ Foundry 1.5.1-stable
  └─ Node.js 18.19.1

✅ 프로젝트 구조 생성
  ├─ contracts/
  ├─ backend/
  ├─ frontend/
  ├─ docker/
  └─ infrastructure/

✅ 프로젝트 초기화
  ├─ Foundry (forge-std 설치)
  ├─ Go modules
  └─ Next.js 14

⏳ Smart Contract 개발 준비
  ├─ OpenZeppelin Contracts 설치 (예정)
  ├─ VaultFactory.sol 작성 (예정)
  └─ IndividualVault.sol 작성 (예정)
```

##### 즉시 진행 가능한 작업
```bash
# 1. OpenZeppelin 설치
cd /root/legacychain/contracts
forge install OpenZeppelin/openzeppelin-contracts
forge install OpenZeppelin/openzeppelin-contracts-upgradeable

# 2. VaultFactory.sol 작성 시작
# 3. IndividualVault.sol 작성 시작
# 4. 단위 테스트 작성
```

#### 시간 기록
- 개발 도구 설치: ~5분
- 프로젝트 구조 생성: ~2분
- 프로젝트 초기화: ~3분
- **총 소요 시간**: ~10분

#### 참고 자료
- [Foundry Book](https://book.getfoundry.sh/)
- [Go Modules](https://go.dev/doc/modules/)
- [Next.js 14 Documentation](https://nextjs.org/docs)
- [PRD 문서](/root/legacychain/docs/PRD.md)

---

## Phase 1: Smart Contract 개발

_작성 예정 (Day 3-6)_

---

## Backend 개발

_작성 예정 (Week 2-3)_

---

## Frontend 개발

_작성 예정 (Week 4)_

---

## 통합 및 배포

_작성 예정_

---

**Last Updated**: 2026-01-12  
**Status**: Phase 0 완료, Phase 1 준비 중
