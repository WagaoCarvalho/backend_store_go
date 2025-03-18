-- ## Category Queries (CRUD)

-- Criar categorias
INSERT INTO users_categories (name, description, created_at, updated_at) 
VALUES 
('Admin', 'Administradores do sistema', NOW(), NOW()),
('Editor', 'Usuários que podem editar conteúdos', NOW(), NOW()),
('Viewer', 'Usuários com acesso somente leitura', NOW(), NOW()),
('Premium', 'Usuários com acesso premium', NOW(), NOW()),
('Guest', 'Usuários temporários', NOW(), NOW());

-- Todas as categorias
SELECT * FROM users_categories;

-- Buscar categoria por ID
SELECT * FROM users_categories WHERE id = 1;

-- Buscar categoria por nome
SELECT * FROM users_categories WHERE name = 'Admin';

-- Ordenar categorias por data de criação (mais recentes primeiro)
SELECT * FROM users_categories ORDER BY created_at DESC;

-- Atualizar nome e descrição da categoria
UPDATE users_categories 
SET name = 'Super Admin', description = 'Usuários com privilégios máximos', updated_at = NOW() 
WHERE id = 1;

-- Deletar categoria por ID
DELETE FROM users_categories WHERE id = 3;
