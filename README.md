## Security: Private Key Management

**IMPORTANT: Never commit private keys to git!**

### Key Files Structure

```
docker/
 .env                   # Environment variables with private keys (gitignored)
 besu/
     entrypoint.sh      # Runtime key file generator
     genesis.json
     static-nodes.json
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

## Documentation

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

##  Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.

##  Contact

For questions or support, please open an issue on GitHub.

---

**Built with  by LegacyChain Team**
