-- +goose Up
CREATE TABLE customers (
    id              BIGSERIAL       PRIMARY KEY,
    name           VARCHAR(255)     NOT NULL,
    date_created    TIMESTAMPTZ     NOT NULL DEFAULT clock_timestamp(),
    date_modified    TIMESTAMPTZ     NOT NULL DEFAULT clock_timestamp()
);

CREATE TABLE orders (
    id              BIGSERIAL       PRIMARY KEY,
    Amount          NUMERIC(10,2)   NOT NULL,
    item            VARCHAR(50)     NOT NULL,
    customer_id     BIGINT          NOT NULL REFERENCES customers(id),
    date_created    TIMESTAMPTZ     NOT NULL DEFAULT clock_timestamp(),
    date_modified    TIMESTAMPTZ     NOT NULL DEFAULT clock_timestamp()
);

-- +goose Down
drop table if exists orders;
drop table if exists customers;
