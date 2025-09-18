CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    cpf VARCHAR(14) UNIQUE,
    cnpj VARCHAR(18) UNIQUE,
    client_type VARCHAR(20) NOT NULL CHECK (client_type IN ('PF','PJ')),
    status BOOLEAN NOT NULL DEFAULT TRUE,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);


-- Busca rápida por nome
CREATE INDEX idx_clients_name ON clients (name);
CREATE INDEX idx_clients_status ON clients (status);
CREATE INDEX idx_clients_type ON clients (client_type);

-- Documentos são UNIQUE (logo já viram índices implícitos):
-- cpf, cnpj, email
