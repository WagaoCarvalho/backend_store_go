CREATE TABLE client_category_relations (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (client_id) REFERENCES clients (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES client_categories (id) ON DELETE CASCADE,
    UNIQUE (client_id, category_id) -- Evita duplicatas de associação
);

CREATE INDEX idx_client_category_client ON client_category_relations (client_id);
CREATE INDEX idx_client_category_category ON client_category_relations (category_id);
