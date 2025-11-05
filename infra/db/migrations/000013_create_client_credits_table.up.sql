CREATE TABLE IF NOT EXISTS client_credits (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    allow_credit BOOLEAN NOT NULL DEFAULT FALSE,
    credit_limit NUMERIC(10,2) NOT NULL DEFAULT 0.00 CHECK (credit_limit >= 0),
    credit_balance NUMERIC(10,2) NOT NULL DEFAULT 0.00 CHECK (credit_balance >= 0 AND credit_balance <= credit_limit),
    description TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_client_credits UNIQUE (client_id)
);

CREATE INDEX idx_client_credits_allow_balance 
    ON client_credits (allow_credit, credit_balance);
