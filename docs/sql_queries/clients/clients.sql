-- ## Client Queries (CRUD)

-- Criar clientes
INSERT INTO clients (name, cpf, cnpj, created_at, updated_at)
VALUES
('Cliente 1', '123.456.789-01', '12.345.678/0001-11', NOW(), NOW()),
('Cliente 2', '987.654.321-00', '23.456.789/0001-22', NOW(), NOW()),
('Cliente 3', '111.222.333-44', '34.567.890/0001-33', NOW(), NOW()),
('Cliente 4', '555.666.777-88', '45.678.901/0001-44', NOW(), NOW()),
('Cliente 5', '999.888.777-66', '56.789.012/0001-55', NOW(), NOW());

-- Todos os clientes
SELECT * FROM clients;

-- Buscar cliente por ID
SELECT * FROM clients WHERE id = 1;

-- Buscar cliente por nome
SELECT * FROM clients WHERE name = 'Cliente 1';

-- Buscar clientes com CPF ou CNPJ
SELECT * FROM clients WHERE cpf IS NOT NULL OR cnpj IS NOT NULL;

-- Buscar clientes ativos
SELECT * FROM clients WHERE status = TRUE;

-- Buscar clientes e suas categorias
SELECT c.*, cat.name AS category_name
FROM clients c
LEFT JOIN client_category_relations r ON c.id = r.client_id
LEFT JOIN clients_categories cat ON r.category_id = cat.id;

-- Ordenar clientes por data de criação (mais recentes primeiro)
SELECT * FROM clients ORDER BY created_at DESC;

-- Atualizar nome de cliente
UPDATE clients
SET name = 'Cliente Atualizado 1', updated_at = NOW()
WHERE id = 1;

-- Atualizar status do cliente
UPDATE clients
SET status = FALSE, updated_at = NOW()
WHERE id = 2;

-- Deletar cliente por ID
DELETE FROM clients WHERE id = 3;
