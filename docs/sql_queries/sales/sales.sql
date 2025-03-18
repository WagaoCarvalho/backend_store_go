-- Inserir 5 vendas fictícias
INSERT INTO sales (customer_id, total_amount, sale_date, created_at, updated_at)
VALUES
(1, 500.00, NOW(), NOW(), NOW()),
(2, 150.00, NOW(), NOW(), NOW()),
(3, 1200.00, NOW(), NOW(), NOW()),
(4, 300.00, NOW(), NOW(), NOW()),
(5, 75.00, NOW(), NOW(), NOW());

-- Consultar todas as vendas
SELECT * FROM sales;

-- Consultar vendas de um cliente específico
SELECT * FROM sales
WHERE customer_id = 1;

-- Consultar vendas feitas por um vendedor específico
SELECT * FROM sales
WHERE seller_id = 2;

-- Consultar o total de vendas por vendedor
SELECT seller_id, SUM(total_amount) AS total_sales
FROM sales
GROUP BY seller_id;

-- Consultar vendas em um intervalo de datas
SELECT * FROM sales
WHERE sale_date BETWEEN '2025-03-01' AND '2025-03-17';

-- Atualizar o valor de uma venda
UPDATE sales
SET total_amount = 120.00
WHERE id = 1;

-- Deletar uma venda
DELETE FROM sales
WHERE id = 1;

-- Contar o número de vendas por vendedor
SELECT seller_id, COUNT(*) AS sales_count
FROM sales
GROUP BY seller_id;
