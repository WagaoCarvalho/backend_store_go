CREATE TABLE IF NOT EXISTS clients_cpf (
    id SERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,

    email VARCHAR(255) NOT NULL,
    CONSTRAINT uq_clients_cpf_email UNIQUE (email),
    CONSTRAINT chk_clients_cpf_email_lower CHECK (email = LOWER(email)),

    cpf CHAR(11) NOT NULL,
    CONSTRAINT uq_clients_cpf_cpf UNIQUE (cpf),
    CONSTRAINT chk_clients_cpf_cpf_format CHECK (cpf ~ '^[0-9]{11}$'),

    description TEXT,

    status BOOLEAN NOT NULL DEFAULT TRUE,

    version INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT chk_clients_cpf_version_positive CHECK (version > 0),

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índices coerentes com uso real
CREATE INDEX idx_clients_cpf_name ON clients_cpf (name);
CREATE INDEX idx_clients_cpf_status_true ON clients_cpf (id) WHERE status = TRUE;
CREATE INDEX idx_clients_cpf_email_lower ON clients_cpf (LOWER(email));
