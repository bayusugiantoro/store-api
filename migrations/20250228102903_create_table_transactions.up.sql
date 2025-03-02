-- Membuat tabel utama transactions
CREATE TABLE IF NOT EXISTS transactions (
    -- Primary key dengan auto-increment
    id SERIAL PRIMARY KEY,

    -- ID customer yang melakukan transaksi
    -- NOT NULL: wajib diisi
    customer_id INTEGER NOT NULL,

    -- Total poin yang digunakan dalam transaksi
    -- NOT NULL: wajib diisi
    -- CHECK: memastikan total_points selalu positif
    total_points INTEGER NOT NULL CHECK (total_points > 0),

    -- Status transaksi
    -- NOT NULL: wajib diisi
    -- CHECK: membatasi nilai yang valid
    -- Hanya bisa: 'pending', 'completed', 'failed'
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),

    -- Timestamp pembuatan transaksi
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Timestamp update transaksi
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Membuat tabel detail transaksi
CREATE TABLE IF NOT EXISTS transaction_items (
    -- Primary key dengan auto-increment
    id SERIAL PRIMARY KEY,

    -- Foreign key ke tabel transactions
    -- NOT NULL: wajib diisi
    -- REFERENCES: memastikan transaction_id valid
    -- ON DELETE CASCADE: hapus item jika transaksi dihapus
    transaction_id INTEGER NOT NULL REFERENCES transactions(id),

    -- Foreign key ke tabel vouchers
    -- NOT NULL: wajib diisi
    -- REFERENCES: memastikan voucher_id valid
    voucher_id INTEGER NOT NULL REFERENCES vouchers(id),

    -- Jumlah poin yang digunakan untuk voucher ini
    -- NOT NULL: wajib diisi
    -- CHECK: memastikan points_used selalu positif
    points_used INTEGER NOT NULL CHECK (points_used > 0),

    -- Timestamp pembuatan item
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk optimasi query
-- Mempercepat pencarian transaksi berdasarkan customer
CREATE INDEX idx_transactions_customer_id ON transactions(customer_id);

-- Mempercepat pencarian item berdasarkan transaction
CREATE INDEX idx_transaction_items_transaction_id ON transaction_items(transaction_id);

-- Mempercepat pencarian item berdasarkan voucher
CREATE INDEX idx_transaction_items_voucher_id ON transaction_items(voucher_id);

-- Trigger untuk auto-update updated_at pada transactions
CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();