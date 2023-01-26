create table if not exists "user"
(
    id bigserial primary key,
    "name" varchar not null,
    user_name varchar not null,
    passwrod varchar not null,
    create_date timestamp not null default now(),
    write_date timestamp not null default now()
);

create table if not exists category
(
    id bigserial primary key,
    "name" varchar not null,
    description varchar,
    create_uid int8 not null default 1,
    write_uid int8 not null default 1,
    create_date timestamp not null default now(),
    write_date timestamp not null default now()
);

create table if not exists seller
(
    id bigserial primary key,
    "name" varchar not null,
    description varchar,
    phone varchar,
    address varchar,
    logo_url varchar,
    manager_id int8 not null,
    create_uid int8 not null default 1,
    write_uid int8 not null default 1,
    create_date timestamp not null default now(),
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
    create_uid int8 not null default 1,
    write_uid int8 not null default 1,
    create_date timestamp not null default now(),
    write_date timestamp not null default now(),
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
    create_uid int8 not null default 1,
    write_uid int8 not null default 1,
    create_date timestamp not null default now(),
    write_date timestamp not null default now(),
    variants json,
    seller_id int8
        constraint product_template_seller_id_fkey
            references seller
            on delete set null,
    category_id int8
        constraint product_template_category_id_fkey
            references category
            on delete set null,
    uom_id int8
        constraint product_template_uom_id_fkey
            references uom
            on delete set null
);

create table if not exists product
(
    id bigserial primary key,
    template_id int8
        constraint product_product_template_id_fkey
            references product_template
            on delete set null,
    "name" varchar not null,
    origin_price double precision not null check ( origin_price >= 0 ) default 0,
    sale_price double precision not null check ( sale_price >= 0 ) default 0,
    "state" varchar not null default 'draft',
    variants json,
    create_uid int8 not null default 1,
    write_uid int8 not null default 1,
    create_date timestamp not null default now(),
    write_date timestamp not null default now()
);