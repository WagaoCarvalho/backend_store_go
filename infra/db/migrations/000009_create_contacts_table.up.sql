CREATE TABLE IF NOT EXISTS contacts (
    id BIGSERIAL PRIMARY KEY,
    contact_name VARCHAR(100) NOT NULL,
    contact_description TEXT,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    cell VARCHAR(20),
    contact_type VARCHAR(30) NOT NULL,

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_contacts_name ON contacts (contact_name);
CREATE INDEX IF NOT EXISTS idx_contacts_type ON contacts (contact_type);
CREATE INDEX IF NOT EXISTS idx_contacts_type_name ON contacts (contact_type, contact_name);

-- Trigger para manter updated_at atualizado
CREATE OR REPLACE FUNCTION set_contacts_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_contacts_updated_at
BEFORE UPDATE ON contacts
FOR EACH ROW
EXECUTE FUNCTION set_contacts_updated_at();
