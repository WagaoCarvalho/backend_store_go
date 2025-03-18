-- Inserir múltiplas relações entre compras e produtos
INSERT INTO purchase_product_relations (purchase_id, product_id, quantity, price)
VALUES
(6, 1, 5, 120.00),  -- Compra 1, Produto 1, 5 unidades, preço 120.00
(7, 3, 2, 250.00),  -- Compra 2, Produto 3, 2 unidades, preço 250.00
(8, 5, 10, 50.00),  -- Compra 3, Produto 5, 10 unidades, preço 50.00
(9, 2, 3, 75.00),  -- Compra 4, Produto 2, 3 unidades, preço 75.00
(10, 4, 1, 500.00);  -- Compra 5, Produto 4, 1 unidade, preço 500.00

-- Consultar todas as relações entre compras e produtos
SELECT * FROM purchase_product_relations;

-- Consultar os produtos comprados em uma compra específica
SELECT ppr.quantity, ppr.price, pr.name
FROM purchase_product_relations ppr
JOIN products pr ON ppr.product_id = pr.id
WHERE ppr.purchase_id = 1;

-- Consultar todas as compras de um produto específico
SELECT ppr.purchase_id, p.total_amount, p.purchase_date
FROM purchase_product_relations ppr
JOIN purchases p ON ppr.purchase_id = p.id
WHERE ppr.product_id = 2;

-- Deletar uma relação específica entre compra e produto
DELETE FROM purchase_product_relations WHERE id = 1;
