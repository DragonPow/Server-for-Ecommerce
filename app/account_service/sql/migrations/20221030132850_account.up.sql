create table if not exists account
(
    id bigserial primary key,
    username varchar not null
        constraint account_username_uniq
            unique,
    password varchar not null ,
    create_date timestamp not null default now(),
    write_date timestamp not null default now()
);

create table if not exists customer_info
(
    id bigserial primary key,
    account_id int8 not null
        constraint customer_info_account_id_fkey
            references account
            on delete cascade,
    name varchar not null ,
    phone varchar,
    address varchar,
    create_date timestamp not null default now(),
    write_date timestamp not null default now()
);