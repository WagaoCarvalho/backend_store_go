-- Inserir 5 serviços
INSERT INTO services (name, description, created_at, updated_at)
VALUES 
('Consultoria Empresarial', 'Análise e estratégia para empresas', NOW(), NOW()),
('Manutenção de Computadores', 'Reparo e otimização de computadores', NOW(), NOW()),
('Desenvolvimento Web', 'Criação de sites e sistemas online', NOW(), NOW()),
('Treinamento Corporativo', 'Capacitação profissional para equipes', NOW(), NOW()),
('Suporte Técnico', 'Atendimento e resolução de problemas técnicos', NOW(), NOW());

-- Read: Buscar todos os serviços
SELECT * FROM services;

-- Buscar um serviço pelo ID
SELECT * FROM services WHERE id = 1;

-- Buscar um serviço pelo nome
SELECT * FROM services WHERE name LIKE '%Desenvolvimento%';

-- Atualizar um serviço
UPDATE services
SET name = 'Desenvolvimento de Software', description = 'Criação de sistemas personalizados', updated_at = NOW()
WHERE id = 3;

-- Deletar um serviço pelo ID
DELETE FROM services WHERE id = 1;