-- Inserir múltiplas relações entre compras e categorias
INSERT INTO purchase_category_relations (purchase_id, category_id)
VALUES
(6, 1),  -- Compra 1 na categoria Tecnologia
(7, 2),  -- Compra 2 na categoria Eletrodomésticos
(8, 3),  -- Compra 3 na categoria Móveis
(9, 4),  -- Compra 4 na categoria Roupas
(10, 5);  -- Compra 5 na categoria Alimentos


-- Consultar todas as relações entre compras e categorias
SELECT * FROM purchase_category_relations;

-- Consultar as categorias de uma compra específica
SELECT pc.name, pc.description
FROM purchase_category_relations pcr
JOIN purchase_categories pc ON pcr.category_id = pc.id
WHERE pcr.purchase_id = 1;

-- Consultar as compras de uma categoria específica
SELECT p.id, p.total_amount, p.purchase_date
FROM purchase_category_relations pcr
JOIN purchases p ON pcr.purchase_id = p.id
WHERE pcr.category_id = 2;

-- Deletar uma relação específica entre compra e categoria
DELETE FROM purchase_category_relations WHERE id = 1;
