CREATE TABLE IF NOT EXISTS supplier_contact_relations (
    contact_id BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    supplier_id BIGINT NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (contact_id, supplier_id)
);