CREATE TABLE user_category_relations (
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, category_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES user_categories (id) ON DELETE CASCADE
);

-- Índices adicionais, se necessário
CREATE INDEX idx_user_category_relations_user_id ON user_category_relations (user_id);
CREATE INDEX idx_user_category_relations_category_id ON user_category_relations (category_id);
