CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    client_id INTEGER REFERENCES clients(id) ON DELETE CASCADE,
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,
    street VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(50) NOT NULL,
    country VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    -- Garante que só uma FK esteja preenchida
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