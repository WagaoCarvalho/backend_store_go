-- Inserindo 5 produtos
INSERT INTO products (name, description, price, stock_quantity, barcode, created_at, updated_at)
VALUES
('Produto 1', 'Descrição do produto 1', 99.99, 50, '1234567890123', NOW(), NOW()),
('Produto 2', 'Descrição do produto 2', 199.99, 30, '2345678901234', NOW(), NOW()),
('Produto 3', 'Descrição do produto 3', 299.99, 20, '3456789012345', NOW(), NOW()),
('Produto 4', 'Descrição do produto 4', 399.99, 10, '4567890123456', NOW(), NOW()),
('Produto 5', 'Descrição do produto 5', 499.99, 15, '5678901234567', NOW(), NOW());

-- Read: Buscar todos os produtos
SELECT * FROM products;

-- Buscar um produto específico pelo id
SELECT * FROM products WHERE id = 1;

-- Buscar um produto pelo nome
SELECT * FROM products WHERE name LIKE '%Produto 1%';

-- Buscar produtos com preço superior a 200
SELECT * FROM products WHERE price > 200;

-- Buscar produtos com estoque abaixo de 20
SELECT * FROM products WHERE stock_quantity < 20;

-- Update: Atualizar o preço de um produto
UPDATE products
SET price = 120.00, updated_at = NOW()
WHERE id = 1;

-- Atualizar a quantidade em estoque de um produto
UPDATE products
SET stock_quantity = stock_quantity + 10, updated_at = NOW()
WHERE id = 1;

-- Delete: Excluir um produto
DELETE FROM products WHERE id = 1;

-- Buscar produtos entre dois preços (faixa de preço)
SELECT * FROM products WHERE price BETWEEN 100 AND 300;

-- Buscar produtos com código de barras específico
SELECT * FROM products WHERE barcode = '1234567890123';

