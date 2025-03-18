-- ## User Queries (CRUD)

-- Criar usuário
INSERT INTO users (username, password_hash, created_at, updated_at) 
VALUES 
('user1', 'hash_senha1', NOW(), NOW()),
('user2', 'hash_senha2', NOW(), NOW()),
('user3', 'hash_senha3', NOW(), NOW()),
('user4', 'hash_senha4', NOW(), NOW()),
('user5', 'hash_senha5', NOW(), NOW());

-- Todos os usuários
SELECT * FROM users;

-- Buscar usuário por ID
SELECT * FROM users WHERE id = 1;

-- Buscar usuário por username
SELECT * FROM users WHERE username = 'user1';

-- Buscar usuários ativos
SELECT * FROM users WHERE status = TRUE;

-- Buscar usuários e suas categorias
SELECT u.*, c.name AS category_name 
FROM users u
LEFT JOIN user_category_relations uc ON u.id = uc.user_id
LEFT JOIN users_categories c ON uc.category_id = c.id;

-- Ordenar usuários por data de criação (mais recentes primeiro)
SELECT * FROM users ORDER BY created_at DESC;

-- Atualizar nome de usuário
UPDATE users 
SET username = 'updated_user1', updated_at = NOW() 
WHERE id = 1;

-- Atualizar status do usuário
UPDATE users 
SET status = FALSE, updated_at = NOW() 
WHERE id = 2;

-- Deletar usuário por ID
DELETE FROM users WHERE id = 3;
