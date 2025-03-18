-- Criar a tabela de categorias de vendas (sale_categories)
CREATE TABLE sale_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL, -- Nome da categoria
    description TEXT,           -- Descrição da categoria
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índice para otimizar buscas pelo nome da categoria
CREATE INDEX idx_sale_categories_name ON sale_categories(name);
