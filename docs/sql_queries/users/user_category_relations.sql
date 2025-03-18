-- ## User-Category Relations Queries (CRUD)

-- Associar usuários a categorias
INSERT INTO user_category_relations (user_id, category_id) 
VALUES 
(1, 1),  -- Usuário 1 -> Categoria Admin
(2, 2),  -- Usuário 2 -> Categoria Editor
(3, 3),  -- Usuário 3 -> Categoria Viewer
(4, 4),  -- Usuário 4 -> Categoria Premium
(5, 5),  -- Usuário 5 -> Categoria Guest
(1, 2),  -- Usuário 1 também pertence à categoria Editor
(2, 3);  -- Usuário 2 também pertence à categoria Viewer

-- Todas as associações entre usuários e categorias
SELECT * FROM user_category_relations;

-- Buscar todas as categorias de um usuário específico (por ID)
SELECT uc.*, c.name AS category_name 
FROM user_category_relations uc
JOIN users_categories c ON uc.category_id = c.id
WHERE uc.user_id = 1;

-- Buscar todos os usuários de uma categoria específica (por ID)
SELECT uc.*, u.username 
FROM user_category_relations uc
JOIN users u ON uc.user_id = u.id
WHERE uc.category_id = 2;

-- Remover uma associação específica entre usuário e categoria
DELETE FROM user_category_relations 
WHERE user_id = 2 AND category_id = 3;

-- Remover todas as associações de um usuário (exemplo: se o usuário for deletado ou alterado)
DELETE FROM user_category_relations WHERE user_id = 3;

-- Remover todas as associações de uma categoria (exemplo: se a categoria for deletada ou alterada)
DELETE FROM user_category_relations WHERE category_id = 5;
