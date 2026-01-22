CREATE TABLE IF NOT EXISTS clients_cpf_credits (
    id SERIAL PRIMARY KEY,

    client_cpf_id INT NOT NULL,
    CONSTRAINT fk_clients_cpf_credits_client
        FOREIGN KEY (client_cpf_id)
        REFERENCES clients_cpf(id)
        ON DELETE CASCADE,

    allow_credit BOOLEAN NOT NULL DEFAULT FALSE,

    credit_limit NUMERIC(14,2) NOT NULL DEFAULT 0.00,
    credit_balance NUMERIC(14,2) NOT NULL DEFAULT 0.00,

    CONSTRAINT chk_credit_limit_non_negative
        CHECK (credit_limit >= 0),

    CONSTRAINT chk_credit_balance_valid
        CHECK (credit_balance >= 0 AND credit_balance <= credit_limit),

    CONSTRAINT chk_credit_allowed_consistency
        CHECK (
            (allow_credit = TRUE)
            OR
            (allow_credit = FALSE AND credit_limit = 0 AND credit_balance = 0)
        ),

    description TEXT,

    version INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT chk_clients_cpf_credits_version_positive CHECK (version > 0),

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_clients_cpf_credits_client UNIQUE (client_cpf_id)
);

-- Índices obrigatórios
CREATE INDEX idx_clients_cpf_credits_client_id
    ON clients_cpf_credits (client_cpf_id);

CREATE INDEX idx_clients_cpf_credits_allow_credit
    ON clients_cpf_credits (allow_credit);
