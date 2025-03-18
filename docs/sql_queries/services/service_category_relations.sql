-- Relacionar serviços com categorias
INSERT INTO service_category_relations (service_id, category_id, created_at)
VALUES 
(1, 1, NOW()),
(2, 2, NOW()),
(3, 3, NOW()),
(4, 3, NOW()),
(5, 5, NOW());

-- Read: Buscar todos os serviços e suas categorias
SELECT s.id, s.name AS service_name, c.name AS category_name
FROM service_category_relations r
JOIN services s ON r.service_id = s.id
JOIN service_categories c ON r.category_id = c.id;

-- Buscar todas as categorias de um serviço específico
SELECT c.name AS category_name 
FROM service_category_relations r
JOIN service_categories c ON r.category_id = c.id
WHERE r.service_id = 1;

-- Buscar todos os serviços de uma categoria específica
SELECT s.name AS service_name 
FROM service_category_relations r
JOIN services s ON r.service_id = s.id
WHERE r.category_id = 1;

-- Atualizar a categoria de um serviço (removendo a antiga e adicionando uma nova)
DELETE FROM service_category_relations WHERE service_id = 1;
INSERT INTO service_category_relations (service_id, category_id, created_at) VALUES (1, 2, NOW());

-- Remover um relacionamento específico entre um serviço e uma categoria
DELETE FROM service_category_relations WHERE service_id = 1 AND category_id = 1;