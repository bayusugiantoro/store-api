-- Menghapus trigger
DROP TRIGGER IF EXISTS update_vouchers_updated_at ON vouchers;

-- Menghapus index
DROP INDEX IF EXISTS idx_vouchers_code;
DROP INDEX IF EXISTS idx_vouchers_brand_id;

-- Menghapus tabel
DROP TABLE IF EXISTS vouchers;