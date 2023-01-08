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
                    create_uid, create_date, write_uid, write_date)
VALUES ($1, $2, $3, $4, $5, $6,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
        now() AT TIME ZONE 'utc',
        now() AT TIME ZONE 'utc') RETURNING id;

-- name: UpdateProductTemplate :exec
UPDATE product_template
SET name            = $1,
    description     = $2,
    default_price   = $3,
    remain_quantity = $4,
    sold_quantity   = $5,
    rating          = $6,
    number_rating   = $7,
    variants        = $8,
    seller_id       = $9,
    category_id     = $10,
    uom_id          = $11,
    write_uid       = case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
    write_date      = now() AT TIME ZONE 'utc'
WHERE id = @id::int8;

-- name: UpdateProduct :exec
UPDATE product
SET template_id  = $1,
    name         = $2,
    origin_price = $3,
    sale_price   = $4,
    state        = $5,
    variants     = $6,
    write_uid    = case when @create_uid::int8 > 0 then @create_uid::int8 else 1 end,
    write_date   = now() AT TIME ZONE 'utc'
WHERE id = @id::int8;