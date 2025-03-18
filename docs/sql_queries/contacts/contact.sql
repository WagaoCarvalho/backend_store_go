-- ## Contacts Queries (CRUD)

-- Criar contato para um usuário, cliente ou fornecedor
INSERT INTO contacts (user_id, client_id, supplier_id, contact_name, contact_position, email, phone, cell, contact_type, created_at, updated_at)
VALUES
(1, NULL, NULL, 'João Silva', 'Gerente de Vendas', 'joao.silva@example.com', '(11) 1234-5678', '(11) 98765-4321', 'Venda', NOW(), NOW()),  -- Contato de um usuário
(NULL, 1, NULL, 'Maria Oliveira', 'Diretora Comercial', 'maria.oliveira@cliente.com', '(21) 2345-6789', '(21) 99876-5432', 'Suporte', NOW(), NOW()),  -- Contato de um cliente
(NULL, NULL, 1, 'Carlos Souza', 'Gestor de Compras', 'carlos.souza@fornecedor.com', '(31) 3456-7890', '(31) 98765-8765', 'Compras', NOW(), NOW());  -- Contato de um fornecedor

-- Selecionar todos os contatos
SELECT * FROM contacts;

-- Buscar contato por ID
SELECT * FROM contacts WHERE id = 1;

-- Buscar contatos de um usuário específico
SELECT * FROM contacts WHERE user_id = 1;

-- Buscar contatos de um cliente específico
SELECT * FROM contacts WHERE client_id = 1;

-- Buscar contatos de um fornecedor específico
SELECT * FROM contacts WHERE supplier_id = 1;

-- Buscar contatos de um tipo específico (ex: "Venda", "Suporte")
SELECT * FROM contacts WHERE contact_type = 'Venda';

-- Atualizar o nome de um contato
UPDATE contacts
SET contact_name = 'João Pereira', updated_at = NOW()
WHERE id = 1;

-- Atualizar telefone de um contato
UPDATE contacts
SET phone = '(11) 1111-1111', updated_at = NOW()
WHERE id = 2;

-- Deletar um contato por ID
DELETE FROM contacts WHERE id = 3;

-- Deletar todos os contatos de um usuário
DELETE FROM contacts WHERE user_id = 1;

-- Deletar todos os contatos de um cliente
DELETE FROM contacts WHERE client_id = 1;

-- Deletar todos os contatos de um fornecedor
DELETE FROM contacts WHERE supplier_id = 1;

-- Ordenar contatos por nome
SELECT * FROM contacts ORDER BY contact_name ASC;

-- Buscar contatos e suas respectivas entidades (usuários, clientes, fornecedores)
SELECT c.*, u.username AS user_name, cl.name AS client_name, s.name AS supplier_name
FROM contacts c
LEFT JOIN users u ON c.user_id = u.id
LEFT JOIN clients cl ON c.client_id = cl.id
LEFT JOIN suppliers s ON c.supplier_id = s.id;
