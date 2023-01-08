-- name: CreateUom :one
INSERT INTO uom(name, seller_id, create_uid, write_uid, create_date, write_date)
VALUES ($1, $2,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        now() AT TIME ZONE 'utc',
        now() AT TIME ZONE 'utc') RETURNING id;

-- name: UpdateUom :exec
UPDATE uom
SET name       = $1,
    seller_id  = $2,
    write_uid  = case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
    write_date = now() AT TIME ZONE 'utc'
WHERE id = @id::int8;
