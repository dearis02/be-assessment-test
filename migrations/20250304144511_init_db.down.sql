DROP TABLE IF EXISTS ledger_entries;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS bank_accounts;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS ledger_entry_type;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS transaction_type;

DROP INDEX IF EXISTS idx_ledger_entries_type;