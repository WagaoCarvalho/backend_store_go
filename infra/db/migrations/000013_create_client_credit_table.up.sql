CREATE TABLE IF NOT EXISTS client_credit (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    allow_credit BOOLEAN NOT NULL DEFAULT FALSE,
    credit_limit DECIMAL(10,2) DEFAULT 0.00 CHECK (credit_limit >= 0),
    credit_balance DECIMAL(10,2) DEFAULT 0.00 CHECK (credit_balance >= 0)
);


CREATE UNIQUE INDEX idx_client_credit_client_id ON client_credit (client_id);
CREATE INDEX idx_client_credit_allow ON client_credit (allow_credit);
CREATE INDEX idx_client_credit_limit ON client_credit (credit_limit);
