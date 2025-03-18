-- ## Addresses Queries (CRUD)

-- Criar endereço para um usuário, cliente ou fornecedor
INSERT INTO addresses (user_id, client_id, supplier_id, street, city, state, country, postal_code, created_at, updated_at)
VALUES
(1, NULL, NULL, 'Rua Exemplo, 123', 'São Paulo', 'SP', 'Brasil', '01000-000', NOW(), NOW()),  -- Endereço de um usuário
(NULL, 1, NULL, 'Avenida Cliente, 456', 'Rio de Janeiro', 'RJ', 'Brasil', '20000-000', NOW(), NOW()),  -- Endereço de um cliente
(NULL, NULL, 1, 'Rua Fornecedor, 789', 'Belo Horizonte', 'MG', 'Brasil', '30000-000', NOW(), NOW());  -- Endereço de um fornecedor

-- Selecionar todos os endereços
SELECT * FROM addresses;

-- Buscar endereço por ID
SELECT * FROM addresses WHERE id = 1;

-- Buscar endereços de um usuário específico
SELECT * FROM addresses WHERE user_id = 1;

-- Buscar endereços de um cliente específico
SELECT * FROM addresses WHERE client_id = 1;

-- Buscar endereços de um fornecedor específico
SELECT * FROM addresses WHERE supplier_id = 1;

-- Atualizar o endereço de um usuário, cliente ou fornecedor
UPDATE addresses
SET street = 'Nova Rua, 123', city = 'Curitiba', state = 'PR', postal_code = '40000-000', updated_at = NOW()
WHERE id = 1;

-- Deletar um endereço por ID
DELETE FROM addresses WHERE id = 2;

-- Deletar todos os endereços de um usuário
DELETE FROM addresses WHERE user_id = 1;

-- Deletar todos os endereços de um cliente
DELETE FROM addresses WHERE client_id = 1;

-- Deletar todos os endereços de um fornecedor
DELETE FROM addresses WHERE supplier_id = 1;

-- Ordenar endereços por cidade e estado
SELECT * FROM addresses ORDER BY city, state;

-- Buscar endereços e suas respectivas entidades (usuários, clientes, fornecedores)
SELECT a.*, u.username AS user_name, c.name AS client_name, s.name AS supplier_name
FROM addresses a
LEFT JOIN users u ON a.user_id = u.id
LEFT JOIN clients c ON a.client_id = c.id
LEFT JOIN suppliers s ON a.supplier_id = s.id;
