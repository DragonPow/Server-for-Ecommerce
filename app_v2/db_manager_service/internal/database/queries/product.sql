-- name: CreateProductTemplate :one
INSERT INTO product_template(name, description, default_price, remain_quantity, sold_quantity, rating, number_rating,
                             variants, seller_id, category_id, uom_id,
                             create_uid, write_uid, create_date, write_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        now() AT TIME ZONE 'utc',
        now() AT TIME ZONE 'utc') RETURNING id;

-- name: CreateProduct :one
INSERT INTO product(template_id, name, origin_price, sale_price, state, variants,
                    create_uid, write_uid, create_date, write_date)
VALUES ($1, $2, $3, $4, $5, $6,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        now() AT TIME ZONE 'utc',
        now() AT TIME ZONE 'utc') RETURNING id;

-- name: UpdateProductTemplate :exec
UPDATE product_template
SET name            = coalesce(sqlc.narg('name'), name),
    description     = coalesce(sqlc.narg('description'), description),
    default_price   = coalesce(sqlc.narg('default_price'), default_price),
    remain_quantity = coalesce(sqlc.narg('remain_quantity'), remain_quantity),
    sold_quantity   = coalesce(sqlc.narg('sold_quantity'), sold_quantity),
    rating          = coalesce(sqlc.narg('rating'), rating),
    number_rating   = coalesce(sqlc.narg('number_rating'), number_rating),
    variants        = coalesce(sqlc.narg('variants'), variants),
    seller_id       = coalesce(sqlc.narg('seller_id'), seller_id),
    category_id     = coalesce(sqlc.narg('category_id'), category_id),
    uom_id          = coalesce(sqlc.narg('uom_id'), uom_id),
    write_uid       = case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
    write_date      = now() AT TIME ZONE 'utc'
WHERE id = @id::int8;

-- name: UpdateProduct :exec
UPDATE product
SET template_id  = coalesce(sqlc.narg('template_id'), template_id),
    name         = coalesce(sqlc.narg('name'), name),
    origin_price = coalesce(sqlc.narg('origin_price'), origin_price),
    sale_price   = coalesce(sqlc.narg('sale_price'), sale_price),
    state        = coalesce(sqlc.narg('state'), state),
    variants     = coalesce(sqlc.narg('variants'), variants),
    write_uid    = case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
    write_date   = now() AT TIME ZONE 'utc'
WHERE id = @id::int8;