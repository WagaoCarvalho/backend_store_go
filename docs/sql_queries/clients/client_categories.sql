-- ## Client Categories Queries (CRUD)

-- Criar categorias de clientes
INSERT INTO client_categories (name, description, created_at, updated_at)
VALUES
('Pessoa Física', 'Clientes individuais que utilizam CPF', NOW(), NOW()),
('Pessoa Jurídica', 'Empresas que utilizam CNPJ', NOW(), NOW()),
('VIP', 'Clientes premium com benefícios exclusivos', NOW(), NOW()),
('Revendedor', 'Clientes que compram para revenda', NOW(), NOW());

-- Todos os clientes e suas categorias
SELECT * FROM client_categories;

-- Buscar categoria por ID
SELECT * FROM client_categories WHERE id = 1;

-- Buscar categoria por nome
SELECT * FROM client_categories WHERE name = 'Pessoa Física';

-- Buscar categorias ativas
SELECT * FROM client_categories WHERE status = TRUE;

-- Ordenar categorias por data de criação (mais recentes primeiro)
SELECT * FROM client_categories ORDER BY created_at DESC;

-- Atualizar nome da categoria
UPDATE client_categories
SET name = 'Categoria Atualizada', updated_at = NOW()
WHERE id = 1;

-- Atualizar status da categoria
UPDATE client_categories
SET status = FALSE, updated_at = NOW()
WHERE id = 2;

-- Deletar categoria por ID
DELETE FROM client_categories WHERE id = 3;
