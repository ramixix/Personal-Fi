-- Drop triggers first
DROP TRIGGER IF EXISTS update_goals_timestamp;
DROP TRIGGER IF EXISTS update_accounts_timestamp;
DROP TRIGGER IF EXISTS update_transactions_timestamp;

-- Drop indexes
DROP INDEX IF EXISTS idx_sync_table;
DROP INDEX IF EXISTS idx_sync_ready;
DROP INDEX IF EXISTS idx_goal_contributions_date;
DROP INDEX IF EXISTS idx_goal_contributions_goal;
DROP INDEX IF EXISTS idx_goals_currency;
DROP INDEX IF EXISTS idx_goals_category;
DROP INDEX IF EXISTS idx_goals_status;
DROP INDEX IF EXISTS idx_account_transactions_date;
DROP INDEX IF EXISTS idx_account_transactions_account;
DROP INDEX IF EXISTS idx_accounts_currency;
DROP INDEX IF EXISTS idx_accounts_name;
DROP INDEX IF EXISTS idx_transactions_currency;
DROP INDEX IF EXISTS idx_transactions_category;
DROP INDEX IF EXISTS idx_transactions_type;
DROP INDEX IF EXISTS idx_transactions_date;

-- Drop tables in reverse order (foreign key dependencies)
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS currency_rates;
DROP TABLE IF EXISTS sync_metadata;
DROP TABLE IF EXISTS goal_contributions;
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS account_transactions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS transactions;