-- Inserir 5 categorias de serviços
INSERT INTO service_categories (name, description, created_at, updated_at)
VALUES 
('Consultoria', 'Serviços de consultoria para empresas e profissionais', NOW(), NOW()),
('TI e Tecnologia', 'Serviços relacionados a tecnologia da informação', NOW(), NOW()),
('Educação e Treinamento', 'Cursos e treinamentos para desenvolvimento profissional', NOW(), NOW()),
('Reparos e Manutenção', 'Serviços de conserto e manutenção de equipamentos', NOW(), NOW()),
('Atendimento e Suporte', 'Suporte técnico e atendimento ao cliente', NOW(), NOW());

-- Read: Buscar todas as categorias de serviços
SELECT * FROM service_categories;

-- Buscar uma categoria específica pelo ID
SELECT * FROM service_categories WHERE id = 1;

-- Buscar uma categoria pelo nome
SELECT * FROM service_categories WHERE name LIKE '%Consultoria%';

-- Atualizar uma categoria de serviço
UPDATE service_categories
SET name = 'Consultoria Empresarial', description = 'Consultoria para empresas e autônomos', updated_at = NOW()
WHERE id = 1;

-- Deletar uma categoria de serviço pelo ID
DELETE FROM service_categories WHERE id = 1;