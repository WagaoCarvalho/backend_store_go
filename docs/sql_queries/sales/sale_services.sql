-- Inserir múltiplas relações entre venda e serviço
INSERT INTO sale_services (sale_id, service_id, quantity, price)
VALUES
(1, 1, 1, 300.00),
(2, 2, 2, 150.00),
(3, 3, 1, 450.00),
(4, 4, 3, 100.00),
(5, 1, 1, 500.00);

-- Consultar todas as relações de venda e serviços
SELECT * FROM sale_services;

-- Consultar os serviços associados a uma venda específica
SELECT s.name, ss.quantity, ss.price, ss.total_amount
FROM sale_services ss
JOIN services s ON ss.service_id = s.id
WHERE ss.sale_id = 1;

-- Consultar as vendas associadas a um serviço específico
SELECT sa.id, sa.total_amount, sa.sale_date
FROM sale_services ss
JOIN sales sa ON ss.sale_id = sa.id
WHERE ss.service_id = 1;

-- Deletar uma relação de venda e serviço
DELETE FROM sale_services
WHERE sale_id = 1 AND service_id = 1;

-- Contar o número de serviços vendidos em uma venda
SELECT SUM(quantity) AS total_services
FROM sale_services
WHERE sale_id = 1;

-- Contar o número de vendas que envolvem um serviço específico
SELECT COUNT(DISTINCT sale_id) AS total_sales
FROM sale_services
WHERE service_id = 1;
