CREATE TABLE service_category_relations (
    service_id INTEGER REFERENCES services(id),
    category_id INTEGER REFERENCES service_categories(id),
    PRIMARY KEY (service_id, category_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_service_category_relations_service ON service_category_relations (service_id);
CREATE INDEX idx_service_category_relations_category ON service_category_relations (category_id);