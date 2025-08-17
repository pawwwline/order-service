-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) NOT NULL UNIQUE,
    track_number VARCHAR(50),
    entry VARCHAR(50),
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    date_created TIMESTAMP DEFAULT now(),
    date_updated TIMESTAMP DEFAULT now(),
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    shardkey VARCHAR(10),
    sm_id INT,
    oof_shard VARCHAR(10)
);

CREATE INDEX idx_orders_date_created ON orders (date_created DESC);


CREATE TABLE deliveries (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    name VARCHAR(100),
    phone VARCHAR(30),
    zip VARCHAR(20),
    city VARCHAR(50),
    address VARCHAR(255),
    region VARCHAR(50),
    email VARCHAR(100)
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    transaction VARCHAR(100),
    request_id VARCHAR(100),
    provider VARCHAR(20),
    amount INT,
    payment_dt INT,
    bank VARCHAR(20),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    chrt_id INT,
    track_number VARCHAR(100),
    price INT,
    rid VARCHAR(100),
    name VARCHAR(100),
    sale INT,
    size VARCHAR(20),
    total_price INT,
    nm_id INT,
    brand VARCHAR(100),
    status INT
);




-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
