CREATE TABLE IF NOT EXISTS sales (
    id SERIAL PRIMARY KEY,
    client_id INTEGER REFERENCES clients(id) ON DELETE SET NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    sale_date TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    total_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00 CHECK (total_amount >= 0),
    total_discount DECIMAL(12,2) DEFAULT 0.00 CHECK (total_discount >= 0),
    total_tax DECIMAL(12,2) DEFAULT 0.00 CHECK (total_tax >= 0),
    payment_type VARCHAR(50) NOT NULL, -- Ex: 'cash', 'card', 'credit'
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- Ex: 'active', 'canceled', 'returned'
    notes TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índices para pesquisas comuns
CREATE INDEX idx_sales_client_id ON sales (client_id);
CREATE INDEX idx_sales_user_id ON sales (user_id);
CREATE INDEX idx_sales_sale_date ON sales (sale_date);
CREATE INDEX idx_sales_status ON sales (status);
