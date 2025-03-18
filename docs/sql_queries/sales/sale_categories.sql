-- Inserir múltiplas categorias de venda
INSERT INTO sale_categories (name, description)
VALUES
('Promoção', 'Categoria para vendas promocionais com descontos especiais'),
('Liquidação', 'Categoria para produtos em liquidação com grandes descontos'),
('Novidades', 'Categoria para lançamentos e novos produtos no mercado'),
('Ofertas Especiais', 'Categoria para ofertas limitadas e descontos exclusivos'),
('VIP', 'Categoria exclusiva para clientes VIP com benefícios personalizados');

-- Consultar todas as categorias de venda
SELECT * FROM sale_categories;

-- Consultar uma categoria de venda específica pelo nome
SELECT * FROM sale_categories
WHERE name = 'Promoção';

-- Atualizar a descrição de uma categoria de venda
UPDATE sale_categories
SET description = 'Descontos sazonais para produtos em promoção'
WHERE id = 1;

-- Deletar uma categoria de venda
DELETE FROM sale_categories
WHERE id = 1;

-- Contar o número de categorias de venda
SELECT COUNT(*) AS total_categories
FROM sale_categories;
