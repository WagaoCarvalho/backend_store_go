-- Criar a tabela de categorias de compras (purchase_categories)
CREATE TABLE purchase_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,             -- Nome da categoria
    description TEXT,                       -- Descrição da categoria
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índice para otimizar buscas por nome de categoria
CREATE INDEX idx_purchase_categories_name ON purchase_categories(name);
