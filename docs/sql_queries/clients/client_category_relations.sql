-- ## Client Category Relations Queries (CRUD)

-- Criar relações entre clientes e categorias
INSERT INTO client_category_relations (client_id, category_id, created_at)
VALUES
(1, 2, NOW()),  -- Cliente 1 relacionado à categoria "Pessoa Jurídica"
(2, 1, NOW()),  -- Cliente 2 relacionado à categoria "Pessoa Física"
(3, 3, NOW()),  -- Cliente 3 relacionado à categoria "VIP"
(4, 4, NOW());  -- Cliente 4 relacionado à categoria "Revendedor"

-- Todos os relacionamentos entre clientes e categorias
SELECT * FROM client_category_relations;

-- Buscar relação por ID
SELECT * FROM client_category_relations WHERE id = 1;

-- Buscar relações de um cliente específico
SELECT * FROM client_category_relations WHERE client_id = 1;

-- Buscar relações de uma categoria específica
SELECT * FROM client_category_relations WHERE category_id = 2;

-- Buscar clientes com suas categorias
SELECT c.id AS client_id, c.name AS client_name, cat.name AS category_name
FROM client_category_relations r
JOIN clients c ON r.client_id = c.id
JOIN clients_categories cat ON r.category_id = cat.id
ORDER BY c.id;

-- Deletar relação de cliente com categoria
DELETE FROM client_category_relations WHERE client_id = 1 AND category_id = 2;

-- Deletar todas as relações de um cliente
DELETE FROM client_category_relations WHERE client_id = 1;

-- Deletar todas as relações de uma categoria
DELETE FROM client_category_relations WHERE category_id = 1;
