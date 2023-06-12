CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS wallet CASCADE;

CREATE TABLE wallets
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name         VARCHAR(255) NOT NULL,
    mobile       VARCHAR(255) NOT NULL,
    avatar       VARCHAR(255),
    balance      BIGINT       NOT NULL,
    description  TEXT         NOT NULL,
    updated_at   TIMESTAMP,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wallet_mobile ON wallets (mobile);

CREATE TABLE transactions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_id   UUID,
    amount      bigint NOT NULL,
    type        VARCHAR(64) NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ
);

CREATE INDEX ON transactions (wallet_id);
CREATE INDEX idx_transaction_wallet_id ON transactions (wallet_id);
CREATE INDEX idx_transaction_type ON transactions (type);

