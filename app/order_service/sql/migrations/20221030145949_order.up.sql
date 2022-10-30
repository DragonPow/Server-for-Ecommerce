create table if not exists product_template
(
    id bigserial primary key,
    name varchar not null,
    default_price double precision,
    uom_name varchar not null,
    inventory_quantity double precision,
    create_uid int8,
    create_date timestamp,
    write_uid int8,
    write_date timestamp
);

create table if not exists product
(
    id bigserial primary key,
    template_id int8 not null
        constraint product_product_template_id_fkey
            references product_template
            on delete cascade,
    name varchar not null,
    origin_price double precision,
    sale_price double precision,
    state varchar not null,
    create_uid int8,
    create_date timestamp,
    write_uid int8,
    write_date timestamp
);

create table if not exists order_bill
(
    id bigserial primary key,
    customer_id int8 not null,
    payment_method varchar not null,
    contact_name varchar,
    contact_phone varchar,
    contact_address varchar,
    total_price double precision,
    ship_cost double precision,
    state varchar not null,
    note varchar,
    create_uid int8,
    create_date timestamp,
    write_uid int8,
    write_date timestamp
);

create table if not exists order_bill_detail
(
    id bigserial primary key,
    order_id int8 not null
        constraint order_bill_detail_order_id_fkey
            references order_bill
            on delete cascade,
    product_template_id int8 not null
        constraint order_bill_detail_product_template_id_fkey
            references product_template,
    quantity double precision,
    unit_price double precision,
    total_price double precision
);

