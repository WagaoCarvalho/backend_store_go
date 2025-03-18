-- Relacionar produtos com categorias
INSERT INTO product_category_relations (product_id, category_id)
VALUES 
(1, 1), -- Smartphone pertence à categoria Eletrônicos
(2, 2), -- Camiseta pertence à categoria Roupas
(3, 3), -- Chocolate pertence à categoria Alimentos
(4, 4), -- Sofá pertence à categoria Móveis
(5, 5); -- Bola de Futebol pertence à categoria Esportes

-- Read: Buscar todos os produtos e suas categorias
SELECT p.id, p.name AS product_name, c.name AS category_name
FROM product_category_relations r
JOIN products p ON r.product_id = p.id
JOIN product_categories c ON r.category_id = c.id;

-- Buscar todas as categorias de um produto específico
SELECT c.name AS category_name 
FROM product_category_relations r
JOIN product_categories c ON r.category_id = c.id
WHERE r.product_id = 1;

-- Buscar todos os produtos de uma categoria específica
SELECT p.name AS product_name 
FROM product_category_relations r
JOIN products p ON r.product_id = p.id
WHERE r.category_id = 1;

-- Atualizar a categoria de um produto (removendo a antiga e adicionando uma nova)
DELETE FROM product_category_relations WHERE product_id = 1;
INSERT INTO product_category_relations (product_id, category_id) VALUES (1, 2);

-- Remover um relacionamento específico entre um produto e uma categoria
DELETE FROM product_category_relations WHERE product_id = 1 AND category_id = 1;