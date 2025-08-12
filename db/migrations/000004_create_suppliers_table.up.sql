CREATE TABLE IF NOT EXISTS suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cnpj VARCHAR(18) UNIQUE,
    cpf VARCHAR(14) UNIQUE,
    status BOOLEAN NOT NULL DEFAULT TRUE,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    -- Garante que exatamente um dos dois campos seja preenchido
    CONSTRAINT chk_cpf_xor_cnpj CHECK (
        (cpf IS NOT NULL AND cnpj IS NULL) OR
        (cpf IS NULL AND cnpj IS NOT NULL)
    )
);

-- Índices para busca mais eficiente
CREATE INDEX idx_suppliers_name ON suppliers (name);
CREATE INDEX idx_suppliers_cnpj ON suppliers (cnpj);
CREATE INDEX idx_suppliers_cpf ON suppliers (cpf);
