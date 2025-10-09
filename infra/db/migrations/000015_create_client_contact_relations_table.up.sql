CREATE TABLE IF NOT EXISTS client_contact_relations (
    contact_id BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    client_id BIGINT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (contact_id, client_id)
);