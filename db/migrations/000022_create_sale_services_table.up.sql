-- Criar a tabela de vendas de serviços (sale_services)
CREATE TABLE sale_services (
    id SERIAL PRIMARY KEY,
    sale_id INTEGER REFERENCES sales(id) ON DELETE CASCADE,         -- Relacionamento com a tabela de vendas
    service_id INTEGER REFERENCES services(id) ON DELETE CASCADE,   -- Relacionamento com a tabela de serviços
    quantity INTEGER NOT NULL,                                       -- Quantidade de serviços vendidos
    price DECIMAL(10, 2) NOT NULL,                                    -- Preço do serviço na venda
    total_amount DECIMAL(10, 2) GENERATED ALWAYS AS (quantity * price) STORED, -- Total da venda do serviço
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Criar índices para otimizar buscas nas relações
CREATE INDEX idx_sale_services_sale_id ON sale_services(sale_id);
CREATE INDEX idx_sale_services_service_id ON sale_services(service_id);
