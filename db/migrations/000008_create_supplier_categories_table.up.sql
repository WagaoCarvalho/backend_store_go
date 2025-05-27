 CREATE TABLE supplier_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_supplier_categories_name ON supplier_categories (name);

INSERT INTO supplier_categories (name, description, created_at, updated_at)
VALUES
('Categoria A', 'Descrição da Categoria A', NOW(), NOW()),
('Categoria B', 'Descrição da Categoria B', NOW(), NOW()),
('Categoria C', 'Descrição da Categoria C', NOW(), NOW()),
('Categoria D', 'Descrição da Categoria D', NOW(), NOW());
