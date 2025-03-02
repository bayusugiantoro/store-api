-- Membuat tabel vouchers
CREATE TABLE IF NOT EXISTS vouchers (
    -- Primary key dengan auto-increment
    id SERIAL PRIMARY KEY,

    -- Foreign key ke tabel brands
    -- NOT NULL: harus diisi
    -- REFERENCES: memastikan brand_id ada di tabel brands
    brand_id INTEGER NOT NULL REFERENCES brands(id),

    -- Kode voucher
    -- VARCHAR(50): maksimal 50 karakter
    -- NOT NULL: harus diisi
    -- UNIQUE: tidak boleh ada kode yang sama
    code VARCHAR(50) NOT NULL UNIQUE,

    -- Nama voucher
    -- VARCHAR(255): maksimal 255 karakter
    -- NOT NULL: harus diisi
    name VARCHAR(255) NOT NULL,

    -- Deskripsi voucher (opsional)
    -- TEXT: tidak ada batasan panjang
    description TEXT,

    -- Jumlah poin yang dibutuhkan
    -- INTEGER: bilangan bulat
    -- NOT NULL: harus diisi
    -- CHECK: memastikan points selalu lebih dari 0
    points INTEGER NOT NULL CHECK (points > 0),

    -- Tanggal kadaluarsa
    -- TIMESTAMP: menyimpan tanggal dan waktu
    -- NOT NULL: harus diisi
    valid_until TIMESTAMP NOT NULL,

    -- Timestamp pembuatan data
    -- DEFAULT CURRENT_TIMESTAMP: otomatis diisi waktu saat ini
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Timestamp update data
    -- Akan diupdate otomatis oleh trigger
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk optimasi query berdasarkan brand_id
-- Berguna untuk query: SELECT * FROM vouchers WHERE brand_id = X
CREATE INDEX idx_vouchers_brand_id ON vouchers(brand_id);

-- Index untuk optimasi pencarian berdasarkan kode
-- Berguna untuk query: SELECT * FROM vouchers WHERE code = 'X'
CREATE INDEX idx_vouchers_code ON vouchers(code);

-- Trigger untuk auto-update updated_at
CREATE TRIGGER update_vouchers_updated_at
    BEFORE UPDATE ON vouchers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();