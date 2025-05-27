CREATE TABLE supplier_category_relations (
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES supplier_categories(id) ON DELETE CASCADE,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (supplier_id, category_id),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES supplier_categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_supplier_category_relations_supplier_id ON supplier_category_relations (supplier_id);
CREATE INDEX idx_supplier_category_relations_category_id ON supplier_category_relations (category_id);


