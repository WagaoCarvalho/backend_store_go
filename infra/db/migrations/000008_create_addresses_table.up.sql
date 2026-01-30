CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    client_cpf_id INTEGER REFERENCES clients_cpf(id) ON DELETE CASCADE,
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,

    street VARCHAR(255) NOT NULL,
    street_number VARCHAR(20) NOT NULL,
    complement VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state CHAR(2) NOT NULL, 
    country VARCHAR(100) NOT NULL DEFAULT 'Brasil',
    postal_code VARCHAR(10) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE, 

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_addresses_city_state ON addresses (city, state);
CREATE INDEX idx_addresses_postal_code ON addresses (postal_code);
