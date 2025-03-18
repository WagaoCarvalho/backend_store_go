-- Criar a tabela de compras (purchases)
CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE SET NULL,  -- Relacionamento com a tabela de fornecedores
    buyer_id INTEGER REFERENCES users(id) ON DELETE SET NULL,       -- Relacionamento com a tabela de usuários (compradores)
    total_amount DECIMAL(10, 2) NOT NULL,                             -- Valor total da compra
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                -- Data da compra
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices para otimizar buscas nas compras
CREATE INDEX idx_purchases_purchase_date ON purchases(purchase_date);
CREATE INDEX idx_purchases_supplier ON purchases(supplier_id);
CREATE INDEX idx_purchases_buyer ON purchases(buyer_id);
