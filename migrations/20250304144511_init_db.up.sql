DO $$
BEGIN
    CREATE TYPE transaction_type AS ENUM (
        'withdraw', 
        'deposit'
    );
    EXCEPTION WHEN duplicate_object THEN 
        RAISE NOTICE 'transaction_type type already exists';
END $$;

DO $$
BEGIN
    CREATE TYPE transaction_status AS ENUM (
        'pending', 
        'completed',
        'failed'
    );
    EXCEPTION WHEN duplicate_object THEN 
        RAISE NOTICE 'transaction_status type already exists';
END $$;

DO $$
BEGIN
    CREATE TYPE ledger_entry_type AS ENUM (
        'debit', 
        'credit'
    );
    EXCEPTION WHEN duplicate_object THEN 
        RAISE NOTICE 'ledger_entry_type type already exists';
END $$;


CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(300) NOT NULL,
    national_identity_number VARCHAR(16) UNIQUE NOT NULL,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS bank_accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    balance DECIMAL(20,0) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    bank_account_id UUID NOT NULL,
    type transaction_type NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id)
);

CREATE TABLE IF NOT EXISTS ledger_entries (
    id UUID PRIMARY KEY,
    bank_account_id UUID NOT NULL,
    transaction_id UUID NOT NULL,
    type ledger_entry_type NOT NULL,
    amount DECIMAL(20,0) NOT NULL,
    balance_after DECIMAL(20,0) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id),
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
);

CREATE INDEX IF NOT EXISTS idx_ledger_entries_type ON ledger_entries USING HASH(type)