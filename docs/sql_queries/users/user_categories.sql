-- ## Category Queries (CRUD)

-- Criar categorias
INSERT INTO user_categories (id, name, description, created_at, updated_at) VALUES
(1, 'Admin', 'Usuários com acesso total ao sistema', NOW(), NOW()),
(2, 'Cliente', 'Usuários que compram produtos', NOW(), NOW()),
(3, 'Fornecedor', 'Usuários que fornecem produtos', NOW(), NOW()),
(4, 'Gerente', 'Usuários que gerenciam vendas e estoque', NOW(), NOW()),
(5, 'Suporte', 'Usuários responsáveis pelo atendimento ao cliente', NOW(), NOW());


-- Todas as categorias
SELECT * FROM user_categories;

-- Buscar categoria por ID
SELECT * FROM user_categories WHERE id = 1;

-- Buscar categoria por nome
SELECT * FROM user_categories WHERE name = 'Admin';

-- Ordenar categorias por data de criação (mais recentes primeiro)
SELECT * FROM user_categories ORDER BY created_at DESC;

-- Atualizar nome e descrição da categoria
UPDATE user_categories 
SET name = 'Super Admin', description = 'Usuários com privilégios máximos', updated_at = NOW() 
WHERE id = 1;

-- Deletar categoria por ID
DELETE FROM user_categories WHERE id = 3;
