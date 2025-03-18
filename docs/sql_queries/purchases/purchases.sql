-- Inserir múltiplas compras
INSERT INTO purchases (supplier_id, buyer_id, total_amount, purchase_date)
VALUES
(1, 2, 500.00, '2025-03-17 10:00:00'),
(2, 3, 150.00, '2025-03-17 11:00:00'),
(3, 4, 250.00, '2025-03-17 12:00:00'),
(4, 5, 350.00, '2025-03-17 13:00:00'),
(1, 6, 450.00, '2025-03-17 14:00:00');

-- Consultar todas as compras
SELECT * FROM purchases;

-- Consultar as compras de um fornecedor específico
SELECT p.id, p.total_amount, p.purchase_date
FROM purchases p
JOIN suppliers s ON p.supplier_id = s.id
WHERE s.id = 1;

-- Consultar as compras feitas por um comprador específico
SELECT p.id, p.total_amount, p.purchase_date
FROM purchases p
JOIN users u ON p.buyer_id = u.id
WHERE u.id = 2;

-- Deletar uma compra específica
DELETE FROM purchases WHERE id = 1;

-- Contar o número de compras feitas por um fornecedor específico
SELECT COUNT(*) AS total_purchases
FROM purchases
WHERE supplier_id = 1;

-- Contar o total de compras feitas por um comprador específico
SELECT SUM(total_amount) AS total_spent
FROM purchases
WHERE buyer_id = 2;
