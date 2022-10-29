create table if not exists import_bill
(
    id bigserial primary key,
    code varchar not null
        constraint import_bill_code_uniq
            unique,
    note text,
    contact_person_name varchar,
    contact_email varchar,
    contact_phone varchar,
    create_uid int8,
    create_date timestamp,
    write_uid int8,
    write_date timestamp
);

create table if not exists import_bill_detail
(
    id bigserial primary key,
    import_id int8 not null
        constraint import_bill_detail_import_id_fkey
            references import_bill
            on delete cascade,
    product_id int8 not null,
    quantity double precision not null default 0,
    create_date timestamp,
    write_date timestamp
);