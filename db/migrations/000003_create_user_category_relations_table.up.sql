CREATE TABLE IF NOT EXISTS user_category_relations (
    user_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY (user_id, category_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES user_categories(id) ON DELETE CASCADE
);
