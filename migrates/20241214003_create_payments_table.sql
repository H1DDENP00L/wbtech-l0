-- +goose Up
CREATE TABLE payments (
        id SERIAL PRIMARY KEY,
        order_uid TEXT REFERENCES orders(order_uid),
        transaction TEXT,
        request_id TEXT,
        currency TEXT,
        provider TEXT,
        amount INTEGER,
        payment_dt BIGINT,
        bank TEXT,
        delivery_cost INTEGER,
        goods_total INTEGER,
        custom_fee INTEGER
);

-- +goose Down
DROP TABLE IF EXISTS payments;

