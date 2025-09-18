CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    cpf VARCHAR(14) UNIQUE,
    cnpj VARCHAR(18) UNIQUE,
    status BOOLEAN NOT NULL DEFAULT TRUE,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT chk_cpf_or_cnpj CHECK (
        (cpf IS NOT NULL AND cnpj IS NULL)
        OR (cpf IS NULL AND cnpj IS NOT NULL)
    )
);

-- Índices auxiliares
CREATE INDEX idx_clients_name ON clients (name);
CREATE INDEX idx_clients_status ON clients (status);
