CREATE TABLE product_category_relations (
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES product_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

CREATE INDEX idx_product_category_relations_product_id ON product_category_relations (product_id);
CREATE INDEX idx_product_category_relations_category_id ON product_category_relations (category_id);
