-- +goose Up
CREATE TABLE items (
        id SERIAL PRIMARY KEY,
        order_uid TEXT REFERENCES orders(order_uid),
        chrt_id INTEGER,
        track_number TEXT,
        price INTEGER,
        rid TEXT,
        name TEXT,
        sale INTEGER,
        size TEXT,
        total_price INTEGER,
        nm_id INTEGER,
        brand TEXT,
        status INTEGER
);

-- +goose Down
DROP TABLE IF EXISTS items;