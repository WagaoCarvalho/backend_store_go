CREATE TABLE user_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_categories_name ON user_categories (name);

INSERT INTO user_categories (id, name, description, created_at, updated_at) VALUES
(1, 'Admin', 'Usuários com acesso total ao sistema', NOW(), NOW()),
(2, 'Cliente', 'Usuários que compram produtos', NOW(), NOW()),
(3, 'Fornecedor', 'Usuários que fornecem produtos', NOW(), NOW()),
(4, 'Gerente', 'Usuários que gerenciam vendas e estoque', NOW(), NOW()),
(5, 'Suporte', 'Usuários responsáveis pelo atendimento ao cliente', NOW(), NOW());