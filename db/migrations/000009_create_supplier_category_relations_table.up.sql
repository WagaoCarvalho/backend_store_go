CREATE TABLE supplier_category_relations (
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES supplier_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (supplier_id, category_id)
);

-- Índices adicionais, se necessário
CREATE INDEX idx_supplier_category_relations_supplier_id ON supplier_category_relations (supplier_id);
CREATE INDEX idx_supplier_category_relations_category_id ON supplier_category_relations (category_id);
