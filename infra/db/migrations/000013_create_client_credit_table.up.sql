CREATE TABLE IF NOT EXISTS client_credit (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    allow_credit BOOLEAN NOT NULL DEFAULT FALSE,
    credit_limit DECIMAL(10,2) NOT NULL DEFAULT 0.00 CHECK (credit_limit >= 0),
    credit_balance DECIMAL(10,2) NOT NULL DEFAULT 0.00 CHECK (credit_balance >= 0 AND credit_balance <= credit_limit),
    description TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT uq_client_credit UNIQUE (client_id),
    CONSTRAINT chk_credit_balance_limit CHECK (credit_balance <= credit_limit)
);

CREATE INDEX idx_client_credit_allow_balance 
    ON client_credit (allow_credit, credit_balance);
