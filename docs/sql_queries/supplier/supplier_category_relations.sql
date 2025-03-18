-- ## Suppliers Categories Relation Queries (CRUD)

-- Criar relações entre fornecedores e categorias
INSERT INTO supplier_category_relations (supplier_id, category_id)
VALUES
(1, 2),  -- Fornecedor 1 relacionado à Categoria 2
(2, 1),  -- Fornecedor 2 relacionado à Categoria 1
(3, 3),  -- Fornecedor 3 relacionado à Categoria 3
(4, 4);  -- Fornecedor 4 relacionado à Categoria 4

-- Selecionar todas as relações entre fornecedores e categorias
SELECT * FROM supplier_category_relations;

-- Buscar relações por ID de fornecedor
SELECT * FROM supplier_category_relations WHERE supplier_id = 1;

-- Buscar relações por ID de categoria
SELECT * FROM supplier_category_relations WHERE category_id = 2;

-- Buscar fornecedores e suas categorias
SELECT s.id AS supplier_id, s.name AS supplier_name, c.name AS category_name
FROM supplier_category_relations r
JOIN suppliers s ON r.supplier_id = s.id
JOIN supplier_categories c ON r.category_id = c.id
ORDER BY s.id;

-- Deletar relação entre fornecedor e categoria por IDs
DELETE FROM supplier_category_relations WHERE supplier_id = 1 AND category_id = 2;

-- Deletar todas as relações de um fornecedor
DELETE FROM supplier_category_relations WHERE supplier_id = 3;

-- Deletar todas as relações de uma categoria
DELETE FROM supplier_category_relations WHERE category_id = 1;
