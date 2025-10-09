CREATE TABLE IF NOT EXISTS user_contact_relations (
    contact_id BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (contact_id, user_id)
);