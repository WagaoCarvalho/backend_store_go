CREATE TABLE sale_category_relations (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER REFERENCES sales(id) ON DELETE CASCADE,          -- Relacionamento com a tabela de vendas
    category_id INTEGER REFERENCES sale_categories(id) ON DELETE CASCADE,  -- Relacionamento com a tabela de categorias de venda
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índice para otimizar buscas nas relações
CREATE INDEX idx_sale_category_relations_sale_id ON sale_category_relations(sale_id);
CREATE INDEX idx_sale_category_relations_category_id ON sale_category_relations(category_id);