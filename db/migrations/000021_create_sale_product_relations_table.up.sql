-- Criar a tabela de relações de venda e produtos (sale_product_relations)
CREATE TABLE sale_product_relations (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER REFERENCES sales(id) ON DELETE CASCADE,          -- Relacionamento com a tabela de vendas
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,    -- Relacionamento com a tabela de produtos
    quantity INTEGER NOT NULL,                                        -- Quantidade de produtos vendidos
    price DECIMAL(10,2) NOT NULL,                                      -- Preço do produto na venda
    total_amount DECIMAL(10,2) GENERATED ALWAYS AS (quantity * price) STORED, -- Total da venda do produto
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices para otimizar buscas nas relações
CREATE INDEX idx_sale_product_relations_sale_id ON sale_product_relations(sale_id);
CREATE INDEX idx_sale_product_relations_product_id ON sale_product_relations(product_id);
