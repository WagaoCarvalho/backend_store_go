-- Criar usuário
INSERT INTO users (username, email, password_hash, created_at, updated_at) 
VALUES 
('user1', 'user1@example.com', 'hash_senha1', NOW(), NOW()),
('user2', 'user2@example.com', 'hash_senha2', NOW(), NOW()),
('user3', 'user3@example.com', 'hash_senha3', NOW(), NOW()),
('user4', 'user4@example.com', 'hash_senha4', NOW(), NOW()),
('user5', 'user5@example.com', 'hash_senha5', NOW(), NOW());

-- Todos os usuários
SELECT * FROM users;

-- Buscar usuário por ID
SELECT * FROM users WHERE id = 1;

-- Buscar usuário por username
SELECT * FROM users WHERE username = 'user1';

-- Buscar usuário por email
SELECT * FROM users WHERE email = 'user1@example.com';

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

-- Atualizar email de usuário
UPDATE users 
SET email = 'updated_email@example.com', updated_at = NOW() 
WHERE id = 1;

-- Atualizar status do usuário
UPDATE users 
SET status = FALSE, updated_at = NOW() 
WHERE id = 2;

-- Deletar usuário por ID
DELETE FROM users WHERE id = 3;

-- Deletar usuário por email
DELETE FROM users WHERE email = 'user3@example.com';