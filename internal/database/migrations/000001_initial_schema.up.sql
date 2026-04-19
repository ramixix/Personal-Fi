-- 1. Transactions
CREATE TABLE IF NOT EXISTS transactions (
    id TEXT PRIMARY KEY,
    date DATETIME NOT NULL,
    amount REAL NOT NULL,
    category TEXT NOT NULL,
    description TEXT NOT NULL,
    type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
    currency_code TEXT NOT NULL DEFAULT 'USD',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 2. Accounts
CREATE TABLE IF NOT EXISTS accounts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    balance REAL NOT NULL DEFAULT 0,
    currency_code TEXT NOT NULL DEFAULT 'USD',
    created DATETIME NOT NULL,  -- Original creation date (user-facing)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Record metadata
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 3. Account Transactions (Depends on Accounts)
CREATE TABLE IF NOT EXISTS account_transactions (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL,
    amount REAL NOT NULL,
    date DATETIME NOT NULL,
    note TEXT DEFAULT '',
    automatic BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

-- 4. Goals
CREATE TABLE IF NOT EXISTS goals (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    target_amount REAL NOT NULL,
    current_amount REAL DEFAULT 0,
    currency_code TEXT NOT NULL DEFAULT 'USD',
    deadline DATETIME,
    has_deadline BOOLEAN DEFAULT 0,
    category TEXT DEFAULT '',
    priority TEXT DEFAULT 'medium' CHECK(priority IN ('high', 'medium', 'low')),
    status TEXT DEFAULT 'active' CHECK(status IN ('active', 'completed', 'paused', 'cancelled')),
    linked_account_id TEXT DEFAULT '',
    created DATETIME NOT NULL,  -- Original creation date
    completed_date DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 5. Goal Contributions
CREATE TABLE IF NOT EXISTS goal_contributions (
    id TEXT PRIMARY KEY,
    goal_id TEXT NOT NULL,
    amount REAL NOT NULL,
    date DATETIME NOT NULL,
    note TEXT DEFAULT '',
    automatic BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (goal_id) REFERENCES goals(id) ON DELETE CASCADE
);

-- 6. Sync metadata table (for future cloud sync)
CREATE TABLE IF NOT EXISTS sync_metadata (
    id TEXT PRIMARY KEY,
    table_name TEXT NOT NULL,
    record_id TEXT NOT NULL,
    action TEXT NOT NULL CHECK(action IN ('insert', 'update', 'delete')),
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    synced BOOLEAN DEFAULT 0,
    device_id TEXT DEFAULT ''
);

-- Currency rates table for future multi-currency support
CREATE TABLE IF NOT EXISTS currency_rates (
    id TEXT PRIMARY KEY,
    from_currency TEXT NOT NULL,
    to_currency TEXT NOT NULL,
    rate REAL NOT NULL,
    date DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(from_currency, to_currency, date)
);

-- User/Device table for multi-device support:
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    device_name TEXT NOT NULL,
    device_type TEXT CHECK(device_type IN ('phone', 'desktop', 'tablet', 'web')),
    last_sync DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category);
CREATE INDEX IF NOT EXISTS idx_transactions_currency ON transactions(currency_code);

CREATE INDEX IF NOT EXISTS idx_accounts_name ON accounts(name);
CREATE INDEX IF NOT EXISTS idx_accounts_currency ON accounts(currency_code);

CREATE INDEX IF NOT EXISTS idx_account_transactions_account ON account_transactions(account_id);
CREATE INDEX IF NOT EXISTS idx_account_transactions_date ON account_transactions(date);

CREATE INDEX IF NOT EXISTS idx_goals_status ON goals(status);
CREATE INDEX IF NOT EXISTS idx_goals_category ON goals(category);
CREATE INDEX IF NOT EXISTS idx_goals_currency ON goals(currency_code);

CREATE INDEX IF NOT EXISTS idx_goal_contributions_goal ON goal_contributions(goal_id);
CREATE INDEX IF NOT EXISTS idx_goal_contributions_date ON goal_contributions(date);

CREATE INDEX IF NOT EXISTS idx_sync_ready ON sync_metadata(synced) WHERE synced = 0;
CREATE INDEX IF NOT EXISTS idx_sync_table ON sync_metadata(table_name, record_id);

-- Triggers to automatically update updated_at timestamp
CREATE TRIGGER IF NOT EXISTS update_transactions_timestamp 
    AFTER UPDATE ON transactions
BEGIN
    UPDATE transactions SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_accounts_timestamp 
    AFTER UPDATE ON accounts
BEGIN
    UPDATE accounts SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_goals_timestamp 
    AFTER UPDATE ON goals
BEGIN
    UPDATE goals SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;