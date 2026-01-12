-- LegacyChain PostgreSQL Initialization Script

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create schemas
CREATE SCHEMA IF NOT EXISTS legacychain;

-- Set search path
SET search_path TO legacychain, public;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_address VARCHAR(42) UNIQUE NOT NULL,
    email VARCHAR(255),
    did VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Vaults table
CREATE TABLE IF NOT EXISTS vaults (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    contract_address VARCHAR(42) UNIQUE NOT NULL,
    owner_address VARCHAR(42) NOT NULL,
    heartbeat_interval BIGINT NOT NULL,
    grace_period BIGINT NOT NULL,
    required_approvals INTEGER NOT NULL,
    is_locked BOOLEAN DEFAULT TRUE,
    last_heartbeat TIMESTAMP,
    unlock_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (owner_address) REFERENCES users(wallet_address)
);

-- Heirs table
CREATE TABLE IF NOT EXISTS heirs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vault_id UUID NOT NULL,
    wallet_address VARCHAR(42) NOT NULL,
    share NUMERIC(10, 2) NOT NULL CHECK (share >= 0 AND share <= 100),
    has_approved BOOLEAN DEFAULT FALSE,
    has_claimed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE,
    FOREIGN KEY (wallet_address) REFERENCES users(wallet_address)
);

-- Heartbeats table (audit log)
CREATE TABLE IF NOT EXISTS heartbeats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vault_id UUID NOT NULL,
    commitment VARCHAR(66) NOT NULL,
    revealed BOOLEAN DEFAULT FALSE,
    tx_hash VARCHAR(66),
    block_number BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE
);

-- Transactions table (audit log)
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vault_id UUID NOT NULL,
    tx_type VARCHAR(50) NOT NULL, -- 'DEPOSIT', 'WITHDRAW', 'CLAIM', 'APPROVE'
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    amount NUMERIC(78, 0), -- Wei amount (max 2^256)
    tx_hash VARCHAR(66) UNIQUE NOT NULL,
    block_number BIGINT,
    status VARCHAR(20) DEFAULT 'PENDING', -- 'PENDING', 'CONFIRMED', 'FAILED'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_users_wallet_address ON users(wallet_address);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_vaults_contract_address ON vaults(contract_address);
CREATE INDEX idx_vaults_owner_address ON vaults(owner_address);
CREATE INDEX idx_vaults_is_locked ON vaults(is_locked);
CREATE INDEX idx_heirs_vault_id ON heirs(vault_id);
CREATE INDEX idx_heirs_wallet_address ON heirs(wallet_address);
CREATE INDEX idx_heartbeats_vault_id ON heartbeats(vault_id);
CREATE INDEX idx_heartbeats_created_at ON heartbeats(created_at);
CREATE INDEX idx_transactions_vault_id ON transactions(vault_id);
CREATE INDEX idx_transactions_tx_hash ON transactions(tx_hash);
CREATE INDEX idx_transactions_status ON transactions(status);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vaults_updated_at BEFORE UPDATE ON vaults
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_heirs_updated_at BEFORE UPDATE ON heirs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert test data (for development)
INSERT INTO users (wallet_address, email) VALUES
    ('0xfe3b557e8fb62b89f4916b721be55ceb828dbd73', 'owner@legacychain.dev'),
    ('0x70997970c51812dc3a010c7d01b50e0d17dc79c8', 'heir1@legacychain.dev'),
    ('0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc', 'heir2@legacychain.dev')
ON CONFLICT (wallet_address) DO NOTHING;

-- Grant permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA legacychain TO legacychain;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA legacychain TO legacychain;
GRANT USAGE ON SCHEMA legacychain TO legacychain;
