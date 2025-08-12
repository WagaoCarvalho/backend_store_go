CREATE TABLE IF NOT EXISTS supplier_category_relations (
    supplier_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (supplier_id, category_id),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES supplier_categories(id) ON DELETE CASCADE
);
