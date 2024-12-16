-- +goose Up
CREATE TABLE deliveries (
        id SERIAL PRIMARY KEY,
        order_uid TEXT REFERENCES orders(order_uid),
        name TEXT,
        phone TEXT,
        zip TEXT,
        city TEXT,
        address TEXT,
        region TEXT,
        email TEXT
);

-- +goose Down
DROP TABLE IF EXISTS deliveries;
