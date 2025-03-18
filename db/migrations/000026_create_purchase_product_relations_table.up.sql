-- Criar a tabela de relações entre compras e produtos (purchase_product_relations)
CREATE TABLE purchase_product_relations (
    id SERIAL PRIMARY KEY,
    purchase_id INTEGER REFERENCES purchases(id) ON DELETE CASCADE, -- Relacionamento com a tabela de compras
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,   -- Relacionamento com a tabela de produtos
    quantity INTEGER NOT NULL, -- Quantidade do produto comprado
    price DECIMAL(10,2) NOT NULL, -- Preço do produto no momento da compra
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices para otimizar buscas
CREATE INDEX idx_purchase_product_relations_purchase_id ON purchase_product_relations(purchase_id);
CREATE INDEX idx_purchase_product_relations_product_id ON purchase_product_relations(product_id);
