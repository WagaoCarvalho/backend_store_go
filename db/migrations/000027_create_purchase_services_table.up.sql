-- Criar a tabela de compras de serviços
CREATE TABLE purchase_services (
    id SERIAL PRIMARY KEY,
    purchase_id INTEGER REFERENCES purchases(id) ON DELETE CASCADE, -- Relacionamento com a tabela de compras
    service_id INTEGER REFERENCES services(id) ON DELETE CASCADE,   -- Relacionamento com a tabela de serviços
    quantity INTEGER NOT NULL, -- Quantidade de serviços adquiridos
    price DECIMAL(10,2) NOT NULL, -- Preço do serviço no momento da compra
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices para otimizar buscas
CREATE INDEX idx_purchase_services_purchase_id ON purchase_services(purchase_id);
CREATE INDEX idx_purchase_services_service_id ON purchase_services(service_id);
