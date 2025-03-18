-- ## Suppliers Queries (CRUD)

-- Criar fornecedores (empresas)
INSERT INTO suppliers (name, cnpj, contact_info, created_at, updated_at)
VALUES
('Fornecedor A', '12.345.678/0001-11', 'Telefone: (11) 1234-5678, email: fornecedorA@example.com', NOW(), NOW()),
('Fornecedor B', '23.456.789/0001-22', 'Telefone: (21) 2345-6789, email: fornecedorB@example.com', NOW(), NOW()),
('Fornecedor C', '34.567.890/0001-33', 'Telefone: (31) 3456-7890, email: fornecedorC@example.com', NOW(), NOW()),
('Fornecedor D', '45.678.901/0001-44', 'Telefone: (41) 4567-8901, email: fornecedorD@example.com', NOW(), NOW()),
('Fornecedor E', '56.789.012/0001-55', 'Telefone: (51) 5678-9012, email: fornecedorE@example.com', NOW(), NOW());

-- Selecionar todos os fornecedores
SELECT * FROM suppliers;

-- Buscar fornecedor por ID
SELECT * FROM suppliers WHERE id = 1;

-- Buscar fornecedor por nome
SELECT * FROM suppliers WHERE name = 'Fornecedor A';

-- Buscar fornecedor por CNPJ
SELECT * FROM suppliers WHERE cnpj = '12.345.678/0001-11';

-- Buscar fornecedores com informações de contato
SELECT * FROM suppliers WHERE contact_info IS NOT NULL;

-- Atualizar o nome do fornecedor
UPDATE suppliers
SET name = 'Fornecedor Atualizado', updated_at = NOW()
WHERE id = 1;

-- Atualizar informações de contato do fornecedor
UPDATE suppliers
SET contact_info = 'Telefone: (11) 9876-5432, email: atualizado@example.com', updated_at = NOW()
WHERE id = 2;

-- Deletar fornecedor por ID
DELETE FROM suppliers WHERE id = 3;

-- Ordenar fornecedores por nome (alfabético)
SELECT * FROM suppliers ORDER BY name ASC;

-- Atualizar status de fornecedor
UPDATE suppliers
SET status = FALSE, updated_at = NOW()
WHERE id = 4;
