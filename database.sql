-- Create custom types
CREATE TYPE user_status AS ENUM ('ACTIVE', 'INACTIVE');
CREATE TYPE transaction_type AS ENUM ('DEPOSIT', 'WITHDRAW', 'TRANSFER');
CREATE TYPE transaction_status AS ENUM ('PENDING', 'COMPLETED', 'FAILED');
CREATE TYPE payment_method_type AS ENUM ('CREDIT_CARD', 'DEBIT_CARD', 'E_WALLET', 'BANK_ACCOUNT');
CREATE TYPE account_type AS ENUM ('WALLET', 'BANK_ACCOUNT');

-- Extensions
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
EXTENSION IF NOT EXISTS "pgcrypto";

-- Core User Management
CREATE TABLE users
(
    user_id       BIGSERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    email         VARCHAR(100) NOT NULL UNIQUE,
    phone_number  VARCHAR(20) UNIQUE,
    password_hash TEXT         NOT NULL,
    status        user_status              DEFAULT 'ACTIVE',
    preferences   JSONB                    DEFAULT '{}',

    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email_phone ON users (email, phone_number);

-- Wallet Management
CREATE TABLE wallets
(
    wallet_id     BIGSERIAL PRIMARY KEY,
    user_id       BIGINT         NOT NULL REFERENCES users (user_id),
    wallet_number UUID                     DEFAULT uuid_generate_v4() UNIQUE,
    balance       NUMERIC(15, 2) NOT NULL  DEFAULT 0.00 CHECK (balance >= 0),
    status        user_status              DEFAULT 'ACTIVE',
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wallets_user ON wallets (user_id);
CREATE INDEX idx_wallets_number ON wallets (wallet_number);

-- Transaction Records
CREATE TABLE transactions
(
    transaction_id        BIGSERIAL PRIMARY KEY,
    source_wallet_id      BIGINT           NOT NULL REFERENCES wallets (wallet_id),
    destination_wallet_id BIGINT REFERENCES wallets (wallet_id),
    transaction_type      transaction_type NOT NULL,
    amount                NUMERIC(15, 2)   NOT NULL CHECK (amount > 0),
    reference_id          UUID                     DEFAULT uuid_generate_v4() UNIQUE,
    status                transaction_status       DEFAULT 'PENDING',
    description           TEXT,
    created_at            TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_wallets ON transactions (source_wallet_id, destination_wallet_id);
CREATE INDEX idx_transactions_reference ON transactions (reference_id);
CREATE INDEX idx_transactions_created_at ON transactions (created_at);

-- Payment Methods
CREATE TABLE payment_methods
(
    payment_method_id BIGSERIAL PRIMARY KEY,
    user_id           BIGINT              NOT NULL REFERENCES users (user_id),
    method_type       payment_method_type NOT NULL,
    account_number    VARCHAR(50)         NOT NULL,
    bank_name         VARCHAR(100),
    is_default        BOOLEAN                  DEFAULT false,
    status            user_status              DEFAULT 'ACTIVE',
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, account_number)
);

CREATE INDEX idx_payment_methods_user ON payment_methods (user_id);

-- Beneficiaries
CREATE TABLE beneficiaries
(
    beneficiary_id     BIGSERIAL PRIMARY KEY,
    user_id            BIGINT       NOT NULL REFERENCES users (user_id),
    beneficiary_name   VARCHAR(100) NOT NULL,
    account_identifier VARCHAR(50)  NOT NULL,
    account_type       account_type NOT NULL,
    bank_name          VARCHAR(100),
    created_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, account_identifier, account_type)
);

CREATE INDEX idx_beneficiaries_user ON beneficiaries (user_id);

-- Transaction Limits
CREATE TABLE transaction_limits
(
    limit_id         BIGSERIAL PRIMARY KEY,
    user_id          BIGINT           NOT NULL REFERENCES users (user_id),
    transaction_type transaction_type NOT NULL,
    daily_limit      NUMERIC(15, 2)   NOT NULL CHECK (daily_limit >= 0),
    monthly_limit    NUMERIC(15, 2)   NOT NULL CHECK (monthly_limit >= 0),
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, transaction_type),
    CONSTRAINT check_monthly_greater_daily CHECK (monthly_limit >= daily_limit)
);

CREATE INDEX idx_transaction_limits_user ON transaction_limits (user_id);

-- Notifications
CREATE TABLE notifications
(
    notification_id   BIGSERIAL PRIMARY KEY,
    user_id           BIGINT       NOT NULL REFERENCES users (user_id),
    title             VARCHAR(100) NOT NULL,
    content           TEXT         NOT NULL,
    notification_type VARCHAR(50)  NOT NULL,
    is_read           BOOLEAN                  DEFAULT false,
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications (user_id);
CREATE INDEX idx_notifications_user_read ON notifications (user_id, is_read) WHERE NOT is_read;

-- Audit Logs
CREATE TABLE audit_logs
(
    log_id      BIGSERIAL PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users (user_id),
    action      VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id   BIGINT      NOT NULL,
    old_value   JSONB,
    new_value   JSONB,
    ip_address  INET,
    user_agent  TEXT,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user ON audit_logs (user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at);


-- Triggers for updated_at
CREATE
OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at
= CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$
language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE
    ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function for handling transactions
CREATE
OR REPLACE FUNCTION process_transaction(
    p_source_wallet_id BIGINT,
    p_destination_wallet_id BIGINT,
    p_amount NUMERIC,
    p_transaction_type transaction_type,
    p_description TEXT DEFAULT NULL
) RETURNS BOOLEAN AS $$
BEGIN
    -- Check sufficient balance
    IF
NOT EXISTS (
        SELECT 1 FROM wallets
        WHERE wallet_id = p_source_wallet_id
        AND balance >= p_amount
        AND
    -- Create transaction record
INSERT INTO transactions (
    source_wallet_id,
    destination_wallet_id, status = 'ACTIVE'
) THEN
    RETURN FALSE;
END IF;

    -- Update source wallet
UPDATE wallets
SET balance = balance - p_amount
WHERE wallet_id = p_source_wallet_id;

-- Update destination wallet if transfer
IF p_destination_wallet_id IS NOT NULL THEN
UPDATE wallets
SET balance = balance + p_amount
WHERE wallet_id = p_destination_wallet_id;
END IF;

    amount,
    transaction_type,
    status,
    description
) VALUES (
             p_source_wallet_id,
             p_destination_wallet_id,
             p_amount,
             p_transaction_type,
             'COMPLETED',
             p_description
         );

RETURN TRUE;
EXCEPTION
    WHEN OTHERS THEN
        RETURN FALSE;
END;
$$
LANGUAGE plpgsql;



-- Create role enum type
CREATE TYPE user_role AS ENUM ('USER', 'ADMIN');

-- Add role column to users table with default value
ALTER TABLE users
    ADD COLUMN role user_role NOT NULL DEFAULT 'USER';

-- Create index for role column
CREATE INDEX idx_users_role ON users(role);

-- Update existing admin users (optional)
-- Replace user_id values with actual admin user IDs
UPDATE users
SET role = 'ADMIN'
WHERE user_id IN (1, 2,11); -- Example admin user IDs

-- Add role to audit_logs table to track role changes
ALTER TABLE audit_logs
    ADD COLUMN user_role user_role;

-- Create trigger to track role changes
CREATE OR REPLACE FUNCTION audit_role_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.role <> NEW.role THEN
        INSERT INTO audit_logs (
            user_id,
            action,
            entity_type,
            entity_id,
            old_value,
            new_value,
            user_role
        ) VALUES (
            NEW.user_id,
            'ROLE_CHANGE',
            'USER',
            NEW.user_id,
            jsonb_build_object('role', OLD.role),
            jsonb_build_object('role', NEW.role),
            NEW.role
        );
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_role_audit
    AFTER UPDATE ON users
    FOR EACH ROW
    WHEN (OLD.role IS DISTINCT FROM NEW.role)
    EXECUTE FUNCTION audit_role_changes();
