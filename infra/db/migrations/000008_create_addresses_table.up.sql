CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    client_id INTEGER REFERENCES clients(id) ON DELETE CASCADE,
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,

    street VARCHAR(255) NOT NULL,
    street_number VARCHAR(20),
    complement VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state CHAR(2) NOT NULL, 
    country VARCHAR(100) NOT NULL DEFAULT 'Brasil',
    postal_code VARCHAR(10) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE, 

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    
    CONSTRAINT chk_only_one_entity CHECK (
        (user_id IS NOT NULL)::int +
        (client_id IS NOT NULL)::int +
        (supplier_id IS NOT NULL)::int = 1
    )
);


CREATE INDEX idx_addresses_city_state ON addresses (city, state);
CREATE INDEX idx_addresses_postal_code ON addresses (postal_code);
CREATE INDEX idx_addresses_user_id ON addresses (user_id);
CREATE INDEX idx_addresses_client_id ON addresses (client_id);
CREATE INDEX idx_addresses_supplier_id ON addresses (supplier_id);


CREATE UNIQUE INDEX idx_addresses_client_unique ON addresses (client_id, street, street_number, postal_code);
