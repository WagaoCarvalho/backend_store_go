CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    manufacturer VARCHAR(255) NOT NULL,
    product_description TEXT,
    cost_price DECIMAL(10, 2) NOT NULL,
    sale_price DECIMAL(10, 2) NOT NULL,
    stock_quantity INTEGER NOT NULL,
    barcode VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_products_product_name ON products (product_name);
CREATE INDEX idx_products_manufacturer ON products (manufacturer);
CREATE INDEX idx_products_cost_price ON products (cost_price);
CREATE INDEX idx_products_sale_price ON products (sale_price);
CREATE INDEX idx_products_stock_quantity ON products (stock_quantity);
