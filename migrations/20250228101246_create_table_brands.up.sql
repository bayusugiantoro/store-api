-- Membuat tabel brands
CREATE TABLE IF NOT EXISTS brands (
    -- Primary key dengan auto-increment
    -- SERIAL: tipe data untuk auto-increment integer
    -- PRIMARY KEY: menjadi unique identifier dan tidak boleh NULL
    id SERIAL PRIMARY KEY,

    -- Nama brand
    -- VARCHAR(255): string dengan maksimal 255 karakter
    -- NOT NULL: field ini wajib diisi
    name VARCHAR(255) NOT NULL,

    -- Deskripsi brand
    -- TEXT: tipe data untuk text panjang tanpa batasan
    -- Boleh NULL karena tidak ada constraint NOT NULL
    description TEXT,

    -- Timestamp pembuatan data
    -- TIMESTAMP: tipe data untuk tanggal dan waktu
    -- DEFAULT CURRENT_TIMESTAMP: otomatis diisi waktu saat ini
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Timestamp terakhir update data
    -- Akan diupdate otomatis oleh trigger saat data diubah
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Membuat fungsi untuk auto-update timestamp
-- Fungsi ini akan dipanggil oleh trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    -- Set kolom updated_at ke waktu saat ini
    NEW.updated_at = CURRENT_TIMESTAMP;
    -- Kembalikan baris yang sudah dimodifikasi
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Membuat trigger untuk auto-update timestamp
-- BEFORE UPDATE: trigger dijalankan sebelum update
-- FOR EACH ROW: trigger dijalankan untuk setiap baris yang diupdate
CREATE TRIGGER update_brands_updated_at
    BEFORE UPDATE ON brands
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();