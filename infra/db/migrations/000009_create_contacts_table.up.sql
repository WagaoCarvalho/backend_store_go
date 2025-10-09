CREATE TABLE IF NOT EXISTS contacts (
    id BIGSERIAL PRIMARY KEY,
    contact_name VARCHAR(100) NOT NULL,
    contact_description VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    cell VARCHAR(20),
    contact_type VARCHAR(20),

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);