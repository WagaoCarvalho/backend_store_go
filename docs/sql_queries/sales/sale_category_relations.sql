-- Inserir múltiplas relações entre venda e categoria
INSERT INTO sale_category_relations (sale_id, category_id)
VALUES
(1, 2),
(2, 3),
(3, 4),
(4, 5),
(5, 1);

-- Consultar todas as relações de venda e categorias
SELECT * FROM sale_category_relations;

-- Consultar as categorias associadas a uma venda específica
SELECT sc.name, sc.description
FROM sale_category_relations scr
JOIN sale_categories sc ON scr.category_id = sc.id
WHERE scr.sale_id = 1;

-- Consultar as vendas associadas a uma categoria específica
SELECT s.id, s.total_amount, s.sale_date
FROM sale_category_relations scr
JOIN sales s ON scr.sale_id = s.id
WHERE scr.category_id = 2;

-- Deletar uma relação de venda e categoria
DELETE FROM sale_category_relations
WHERE sale_id = 1 AND category_id = 2;

-- Contar o número de categorias associadas a uma venda
SELECT COUNT(*) AS total_categories
FROM sale_category_relations
WHERE sale_id = 1;

-- Contar o número de vendas associadas a uma categoria
SELECT COUNT(*) AS total_sales
FROM sale_category_relations
WHERE category_id = 2;
