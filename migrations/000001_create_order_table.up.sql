CREATE TABLE orders (
    order_uid VARCHAR(255) UNIQUE,
    track_number VARCHAR(255),
    entry VARCHAR(255),
    delivery jsonb,
    payment jsonb,
    items jsonb,
    locale varchar(2),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shard_key VARCHAR(6),
    sm_id int,
    date_created timestamp,
    oof_shard VARCHAR(6)
);