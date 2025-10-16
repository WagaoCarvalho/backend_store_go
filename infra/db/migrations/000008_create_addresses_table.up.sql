CREATE TABLE IF NOT EXISTS addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    client_id INTEGER REFERENCES clients(id) ON DELETE CASCADE,
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
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    -- Garante que apenas um tipo de entidade (user, client, supplier) possua endereço
    CONSTRAINT chk_only_one_entity CHECK (
        (user_id IS NOT NULL)::int +
        (client_id IS NOT NULL)::int +
        (supplier_id IS NOT NULL)::int = 1
    )
);

-- Índices úteis para buscas por localização
CREATE INDEX idx_addresses_city_state ON addresses (city, state);
CREATE INDEX idx_addresses_postal_code ON addresses (postal_code);

-- Restrições únicas por entidade (evita duplicidade de endereço)
CREATE UNIQUE INDEX idx_addresses_client_unique 
    ON addresses (client_id, street, street_number, postal_code)
    WHERE client_id IS NOT NULL;

CREATE UNIQUE INDEX idx_addresses_user_unique 
    ON addresses (user_id, street, street_number, postal_code)
    WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX idx_addresses_supplier_unique 
    ON addresses (supplier_id, street, street_number, postal_code)
    WHERE supplier_id IS NOT NULL;
