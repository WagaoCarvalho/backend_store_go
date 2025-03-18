-- ## Supplier Categories Queries (CRUD)

-- Criar categorias de fornecedores
INSERT INTO supplier_categories (name, description, created_at, updated_at)
VALUES
('Categoria A', 'Descrição da Categoria A', NOW(), NOW()),
('Categoria B', 'Descrição da Categoria B', NOW(), NOW()),
('Categoria C', 'Descrição da Categoria C', NOW(), NOW()),
('Categoria D', 'Descrição da Categoria D', NOW(), NOW());

-- Selecionar todas as categorias de fornecedores
SELECT * FROM supplier_categories;

-- Buscar categoria por ID
SELECT * FROM supplier_categories WHERE id = 1;

-- Buscar categoria por nome
SELECT * FROM supplier_categories WHERE name = 'Categoria A';

-- Atualizar nome da categoria
UPDATE supplier_categories
SET name = 'Categoria Atualizada', updated_at = NOW()
WHERE id = 1;

-- Atualizar descrição da categoria
UPDATE supplier_categories
SET description = 'Nova descrição para Categoria A', updated_at = NOW()
WHERE id = 2;

-- Deletar categoria por ID
DELETE FROM supplier_categories WHERE id = 3;

-- Ordenar categorias por nome (alfabético)
SELECT * FROM supplier_categories ORDER BY name ASC;

-- Ordenar categorias por data de criação (mais recentes primeiro)
SELECT * FROM supplier_categories ORDER BY created_at DESC;
