-- Inserir múltiplas relações entre compras e serviços
INSERT INTO purchase_services (purchase_id, service_id, quantity, price)
VALUES
(6, 1, 1, 500.00),  -- Compra 1, Serviço 1, 1 unidade, preço 500.00
(7, 3, 2, 750.00),  -- Compra 2, Serviço 3, 2 unidades, preço 750.00
(8, 5, 1, 200.00),  -- Compra 3, Serviço 5, 1 unidade, preço 200.00
(9, 2, 3, 450.00),  -- Compra 4, Serviço 2, 3 unidades, preço 450.00
(10, 4, 1, 900.00);  -- Compra 5, Serviço 4, 1 unidade, preço 900.00


-- Consultar todas as relações entre compras e serviços
SELECT * FROM purchase_services;

-- Consultar os serviços adquiridos em uma compra específica
SELECT ps.quantity, ps.price, s.name
FROM purchase_services ps
JOIN services s ON ps.service_id = s.id
WHERE ps.purchase_id = 1;

-- Consultar todas as compras de um serviço específico
SELECT ps.purchase_id, p.total_amount, p.purchase_date
FROM purchase_services ps
JOIN purchases p ON ps.purchase_id = p.id
WHERE ps.service_id = 2;

-- Deletar uma relação específica entre compra e serviço
DELETE FROM purchase_services WHERE id = 1;
