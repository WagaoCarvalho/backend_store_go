-- Inserir múltiplas relações entre venda e produto
INSERT INTO sale_product_relations (sale_id, product_id, quantity, price)
VALUES
(1, 1, 2, 50.00),
(2, 3, 1, 100.00),
(3, 2, 5, 20.00),
(4, 4, 3, 75.00),
(5, 1, 1, 150.00);

-- Consultar todas as relações de venda e produtos
SELECT * FROM sale_product_relations;

-- Consultar os produtos associados a uma venda específica
SELECT p.name, spr.quantity, spr.price, spr.total_amount
FROM sale_product_relations spr
JOIN products p ON spr.product_id = p.id
WHERE spr.sale_id = 1;

-- Consultar as vendas associadas a um produto específico
SELECT s.id, s.total_amount, s.sale_date
FROM sale_product_relations spr
JOIN sales s ON spr.sale_id = s.id
WHERE spr.product_id = 1;

-- Deletar uma relação de venda e produto
DELETE FROM sale_product_relations
WHERE sale_id = 1 AND product_id = 1;

-- Contar o número de produtos vendidos em uma venda
SELECT SUM(quantity) AS total_products
FROM sale_product_relations
WHERE sale_id = 1;

-- Contar o número de vendas que envolvem um produto específico
SELECT COUNT(DISTINCT sale_id) AS total_sales
FROM sale_product_relations
WHERE product_id = 1;
