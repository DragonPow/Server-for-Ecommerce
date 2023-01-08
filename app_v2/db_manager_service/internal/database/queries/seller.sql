-- name: CreateSeller :one
INSERT INTO seller
(name, description, phone, address, logo_url, manager_id, create_uid, create_date, write_uid, write_date)
VALUES ($1, $2, $3, $4, $5, $6,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        now() AT TIME ZONE 'utc',
        now() AT TIME ZONE 'utc') RETURNING id;

-- name: UpdateSeller :exec
UPDATE seller
SET name        = $1,
    description = $2,
    phone       = $3,
    address     = $4,
    logo_url    = $5,
    manager_id  = $6,
    write_uid   = case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
    write_date  = now() AT TIME ZONE 'utc'
WHERE id = @id::int8;