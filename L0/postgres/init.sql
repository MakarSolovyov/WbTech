CREATE SCHEMA IF NOT EXISTS order_service
    AUTHORIZATION "goAdmin";

CREATE TABLE IF NOT EXISTS order_service.delivery
(
    delivery_id character varying COLLATE pg_catalog."default" NOT NULL,
    name character varying COLLATE pg_catalog."default",
    phone character varying COLLATE pg_catalog."default",
    zip character varying COLLATE pg_catalog."default",
    city character varying COLLATE pg_catalog."default",
    address character varying COLLATE pg_catalog."default",
    region character varying COLLATE pg_catalog."default",
    email character varying COLLATE pg_catalog."default",
    CONSTRAINT delivery_pkey PRIMARY KEY (delivery_id)
);

CREATE TABLE IF NOT EXISTS order_service.items
(
    chrt_id integer NOT NULL,
    track_number character varying COLLATE pg_catalog."default",
    price double precision,
    rid character varying COLLATE pg_catalog."default",
    name character varying COLLATE pg_catalog."default",
    sale integer,
    size character varying COLLATE pg_catalog."default",
    total_price double precision,
    nm_id integer,
    brand character varying COLLATE pg_catalog."default",
    status integer,
    CONSTRAINT items_pkey PRIMARY KEY (chrt_id)
);

CREATE TABLE IF NOT EXISTS order_service.payment
(
    transaction character varying COLLATE pg_catalog."default" NOT NULL,
    request_id character varying COLLATE pg_catalog."default",
    currency character varying COLLATE pg_catalog."default",
    provider character varying COLLATE pg_catalog."default",
    amount double precision,
    payment_dt integer,
    bank character varying COLLATE pg_catalog."default",
    delivery_cost double precision,
    goods_total double precision,
    custom_fee double precision,
    CONSTRAINT payment_pkey PRIMARY KEY (transaction)
);

CREATE TABLE IF NOT EXISTS order_service.orders
(
    order_uid character varying COLLATE pg_catalog."default" NOT NULL,
    track_number character varying COLLATE pg_catalog."default",
    entry character varying COLLATE pg_catalog."default",
    delivery character varying COLLATE pg_catalog."default",
    payment character varying COLLATE pg_catalog."default",
    locale character varying COLLATE pg_catalog."default",
    internal_signature character varying COLLATE pg_catalog."default",
    customer_id character varying COLLATE pg_catalog."default",
    delivery_service character varying COLLATE pg_catalog."default",
    shardkey character varying COLLATE pg_catalog."default",
    sm_id integer,
    date_created timestamp with time zone,
    oof_shard character varying COLLATE pg_catalog."default",
    items character varying COLLATE pg_catalog."default",
    CONSTRAINT orders_pkey PRIMARY KEY (order_uid),
    CONSTRAINT fk_delivery_delivery_delivery_id FOREIGN KEY (delivery)
        REFERENCES order_service.delivery (delivery_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT fk_payment_payment_transaction FOREIGN KEY (payment)
        REFERENCES order_service.payment (transaction) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);