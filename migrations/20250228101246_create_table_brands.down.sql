-- Menghapus trigger terlebih dahulu
DROP TRIGGER IF EXISTS update_brands_updated_at ON brands;

-- Menghapus fungsi
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Menghapus tabel
DROP TABLE IF EXISTS brands;