create table if not exists category
(
    id bigserial primary key,
    name varchar not null unique
);

create table if not exists seller
(
    id bigserial primary key,
    name varchar not null,
    description varchar,
    phone varchar,
    address varchar,
    logo_url varchar,
    manager_id int8 not null,
    create_uid int8,
    create_date timestamp not null default now(),
    write_uid int8,
    write_date timestamp not null default now()
);

create table if not exists uom
(
    id bigserial primary key,
    name varchar not null,
    seller_id int8 not null
        constraint uom_seller_id_fkey
            references seller
            on delete cascade,
    constraint name_seller_id_unique unique(name, seller_id)
);

create table if not exists product_template
(
    id bigserial primary key,
    name varchar not null,
    description varchar,
    default_price double precision not null check ( default_price >= 0 ) default 0,
    remain_quantity double precision not null check ( remain_quantity >= 0 ) default 0,
    sold_quantity double precision not null check ( sold_quantity >= 0 ) default 0,
    rating double precision not null check ( rating >= 0 ) default 0,
    number_rating int8 not null check ( number_rating >= 0 ) default 0,
    create_uid int8,
    create_date timestamp not null default now(),
    write_uid int8,
    write_date timestamp not null default now(),
    variants jsonb,
    seller_id int8 not null
        constraint product_template_seller_id_fkey
            references seller
            on delete cascade,
    category_id int8 not null
        constraint product_template_category_id_fkey
            references category
            on delete cascade,
    uom_id int8 not null
        constraint product_template_uom_id_fkey
            references uom
            on delete cascade
);

create table if not exists product
(
    id bigserial primary key,
    template_id int8 not null
        constraint product_product_template_id_fkey
            references product_template
            on delete cascade,
    name varchar not null,
    origin_price double precision not null check ( origin_price >= 0 ) default 0,
    sale_price double precision not null check ( sale_price >= 0 ) default 0,
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

create table if not exists order_shipping
(
    id bigserial primary key,
    order_id int8 not null
        constraint order_shipping_order_id_fkey
            references order_bill
            on delete cascade,
    state varchar not null,
    shipping_name varchar,
    shipping_phone varchar,
    shipping_address varchar,
    create_uid int8,
    create_date timestamp,
    write_uid int8,
    write_date timestamp
);

create table if not exists order_shipping_detail
(
    id bigserial primary key,
    shipping_id int8 not null
        constraint order_shipping_detail_shipping_id_fkey
            references order_shipping
            on delete cascade,
    order_detail_id int8 not null
        constraint order_shipping_detail_order_detail_id_fkey
            references order_bill_detail
            on delete cascade,
    product_id int8 not null
        constraint order_shipping_detail_product_id_fkey
            references product,
    quantity double precision
);