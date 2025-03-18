-- Criar a tabela de relações entre compras e categorias (purchase_category_relations)
CREATE TABLE purchase_category_relations (
    id SERIAL PRIMARY KEY,
    purchase_id INTEGER REFERENCES purchases(id) ON DELETE CASCADE,    -- Relacionamento com a tabela de compras
    category_id INTEGER REFERENCES purchase_categories(id) ON DELETE CASCADE,  -- Relacionamento com a tabela de categorias de compras
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índice para otimizar buscas nas relações entre compras e categorias
CREATE INDEX idx_purchase_category_relations_purchase_id ON purchase_category_relations(purchase_id);
CREATE INDEX idx_purchase_category_relations_category_id ON purchase_category_relations(category_id);
