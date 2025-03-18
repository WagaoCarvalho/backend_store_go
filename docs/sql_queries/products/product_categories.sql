-- Inserir 5 categorias de produtos
INSERT INTO product_categories (name, description, created_at, updated_at)
VALUES 
('Eletrônicos', 'Produtos eletrônicos e gadgets', NOW(), NOW()),
('Roupas', 'Vestuário masculino e feminino', NOW(), NOW()),
('Alimentos', 'Comida e bebidas em geral', NOW(), NOW()),
('Móveis', 'Móveis para casa e escritório', NOW(), NOW()),
('Esportes', 'Artigos esportivos e fitness', NOW(), NOW());

-- Read: Buscar todas as categorias
SELECT * FROM product_categories;

-- Buscar uma categoria pelo ID
SELECT * FROM product_categories WHERE id = 1;

-- Buscar uma categoria pelo nome
SELECT * FROM product_categories WHERE name LIKE '%Eletrônicos%';

-- Update: Atualizar o nome e descrição de uma categoria
UPDATE product_categories
SET name = 'Eletrônicos e Acessórios', description = 'Todos os eletrônicos e acessórios', updated_at = NOW()
WHERE id = 1;

-- Delete: Remover uma categoria pelo ID
DELETE FROM product_categories WHERE id = 1;