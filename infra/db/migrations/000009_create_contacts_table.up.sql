CREATE TABLE IF NOT EXISTS contacts (
    id BIGSERIAL PRIMARY KEY,
    contact_name VARCHAR(100) NOT NULL,
    contact_description TEXT,              
    email VARCHAR(255),                    
    phone VARCHAR(20),
    cell VARCHAR(20),
    contact_type VARCHAR(20) NOT NULL,     

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contacts_name ON contacts (contact_name);

CREATE UNIQUE INDEX idx_contacts_email ON contacts (email) WHERE email IS NOT NULL;

CREATE INDEX idx_contacts_type ON contacts (contact_type);
