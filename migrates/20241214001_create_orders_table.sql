-- +goose Up
CREATE TABLE orders (
        order_uid TEXT PRIMARY KEY,
        track_number TEXT,
        entry TEXT,
        locale TEXT,
        internal_signature TEXT,
        customer_id TEXT,
        delivery_service TEXT,
        shardkey TEXT,
        sm_id INTEGER,
        date_created TEXT,
        oof_shard TEXT
);

-- +goose Down
DROP TABLE IF EXISTS orders;

