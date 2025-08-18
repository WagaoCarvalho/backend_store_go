CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    supplier_id INTEGER REFERENCES suppliers(id) ON DELETE SET NULL,
    product_name VARCHAR(255) NOT NULL,
    manufacturer VARCHAR(255) NOT NULL,
    product_description TEXT,
    cost_price DECIMAL(10, 2) NOT NULL,
    sale_price DECIMAL(10, 2) NOT NULL,
    stock_quantity INTEGER NOT NULL,
    barcode VARCHAR(255) UNIQUE,
    status BOOLEAN NOT NULL DEFAULT TRUE,
    version INTEGER NOT NULL DEFAULT 1,
    allow_discount BOOLEAN NOT NULL DEFAULT FALSE,
    min_discount_percent DECIMAL(5, 2) DEFAULT 0.00 CHECK (min_discount_percent >= 0),
    max_discount_percent DECIMAL(5, 2) DEFAULT 0.00 CHECK (max_discount_percent >= 0),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_discount_range CHECK (max_discount_percent >= min_discount_percent)
);

CREATE INDEX idx_products_product_name ON products (product_name);
CREATE INDEX idx_products_manufacturer ON products (manufacturer);
CREATE INDEX idx_products_cost_price ON products (cost_price);
CREATE INDEX idx_products_sale_price ON products (sale_price);
CREATE INDEX idx_products_stock_quantity ON products (stock_quantity);
