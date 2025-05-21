CREATE TABLE contacts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    client_id INTEGER REFERENCES clients(id) ON DELETE SET NULL,
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE SET NULL,
    contact_name VARCHAR(255),
    contact_position VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    cell VARCHAR(20),
    contact_type VARCHAR(50),
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_contacts_user_id ON contacts (user_id);
CREATE INDEX idx_contacts_client_id ON contacts (client_id);
CREATE INDEX idx_contacts_supplier_id ON contacts (supplier_id);
CREATE INDEX idx_contacts_email ON contacts (email);
CREATE INDEX idx_contacts_phone ON contacts (phone);
CREATE INDEX idx_contacts_cell ON contacts (cell);
