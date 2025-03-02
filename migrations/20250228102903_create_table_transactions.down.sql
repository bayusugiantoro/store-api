-- Menghapus dalam urutan yang benar (karena ada foreign key)
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP INDEX IF EXISTS idx_transaction_items_voucher_id;
DROP INDEX IF EXISTS idx_transaction_items_transaction_id;
DROP INDEX IF EXISTS idx_transactions_customer_id;
DROP TABLE IF EXISTS transaction_items;
DROP TABLE IF EXISTS transactions;