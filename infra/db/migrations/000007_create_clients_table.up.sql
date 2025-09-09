CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    cpf VARCHAR(14) UNIQUE,
    cnpj VARCHAR(18) UNIQUE,
    status BOOLEAN NOT NULL DEFAULT TRUE,
    allow_credit BOOLEAN NOT NULL DEFAULT FALSE,
    credit_limit DECIMAL(10,2) DEFAULT 0.00 CHECK (credit_limit >= 0),
    credit_balance DECIMAL(10,2) DEFAULT 0.00 CHECK (credit_balance >= 0),
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_clients_name ON clients (name);
CREATE INDEX idx_clients_cpf ON clients (cpf);
CREATE INDEX idx_clients_cnpj ON clients (cnpj);
