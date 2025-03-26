-- Inserindo 5 produtos (com todos os campos)
INSERT INTO products (
    product_name, 
    manufacturer, 
    product_description, 
    cost_price, 
    sale_price, 
    stock_quantity, 
    barcode, 
    created_at, 
    updated_at
)
VALUES
('Smartphone X', 'Samsung', 'Smartphone top de linha', 800.00, 999.99, 50, '1234567890123', NOW(), NOW()),
('Notebook Pro', 'Dell', 'Notebook para profissionais', 2500.00, 2999.99, 30, '2345678901234', NOW(), NOW()),
('Tablet Lite', 'Apple', 'Tablet compacto e potente', 400.00, 499.99, 20, '3456789012345', NOW(), NOW()),
('Fone Bluetooth', 'Sony', 'Fone sem fio premium', 150.00, 199.99, 10, '4567890123456', NOW(), NOW()),
('Smartwatch 4', 'Xiaomi', 'Relógio inteligente', 120.00, 149.99, 15, '5678901234567', NOW(), NOW());

-- Read: Buscar todos os produtos
SELECT * FROM products;

-- Buscar um produto específico pelo id
SELECT * FROM products WHERE id = 1;

-- Buscar um produto pelo nome (usando product_name em vez de name)
SELECT * FROM products WHERE product_name LIKE '%Smartphone%';

-- Buscar produtos com preço de venda superior a 200 (usando sale_price em vez de price)
SELECT * FROM products WHERE sale_price > 200;

-- Buscar produtos com estoque abaixo de 20
SELECT * FROM products WHERE stock_quantity < 20;

-- Update: Atualizar o preço de venda de um produto (usando sale_price)
UPDATE products
SET sale_price = 1099.99, updated_at = NOW()
WHERE id = 1;

-- Atualizar a quantidade em estoque de um produto
UPDATE products
SET stock_quantity = stock_quantity + 10, updated_at = NOW()
WHERE id = 1;

-- Delete: Excluir um produto
DELETE FROM products WHERE id = 1;

-- Buscar produtos entre dois preços (faixa de preço de venda)
SELECT * FROM products WHERE sale_price BETWEEN 100 AND 300;

-- Buscar produtos com código de barras específico
SELECT * FROM products WHERE barcode = '1234567890123';

-- Buscar produtos por fabricante
SELECT * FROM products WHERE manufacturer = 'Samsung';

-- Buscar produtos com margem de lucro alta (sale_price > 1.5 * cost_price)
SELECT * FROM products WHERE sale_price > (cost_price * 1.5);

-- Atualizar preço de custo e venda simultaneamente
UPDATE products
SET cost_price = 850.00, sale_price = 1049.99, updated_at = NOW()
WHERE id = 2;