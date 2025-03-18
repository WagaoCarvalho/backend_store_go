-- Inserir múltiplas categorias de compras
INSERT INTO purchase_categories (name, description)
VALUES
('Tecnologia', 'Categoria para compras de produtos eletrônicos e tecnológicos'),
('Eletrodomésticos', 'Categoria para compras de eletrodomésticos como geladeiras, fogões, etc.'),
('Móveis', 'Categoria para compras de móveis para casa e escritório'),
('Roupas', 'Categoria para compras de roupas e acessórios de moda'),
('Alimentos', 'Categoria para compras de alimentos e bebidas');

-- Consultar todas as categorias de compras
SELECT * FROM purchase_categories;

-- Consultar uma categoria específica de compra
SELECT * FROM purchase_categories WHERE name = 'Tecnologia';

-- Deletar uma categoria específica de compra
DELETE FROM purchase_categories WHERE id = 1;

-- Atualizar o nome e a descrição de uma categoria de compra
UPDATE purchase_categories
SET name = 'Eletrodomésticos', description = 'Categoria para compras de eletrodomésticos'
WHERE id = 1;
