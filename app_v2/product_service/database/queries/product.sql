-- name: GetProductDetails :many
SELECT p.*,
       c.id category_id, c.name category_name,
       u.id uom_id, u.name uom_name,
       s.id seller_id, s.name seller_name, s.logo_url seller_logo, s.address seller_address,
       pt.name template_name, pt.rating, pt.number_rating, pt.description template_description, pt.remain_quantity, pt.sold_quantity,
       us1."name" create_name, us2."name" write_name
FROM product p
JOIN product_template pt on pt.id = p.template_id
JOIN category c on c.id = pt.category_id
JOIN uom u on u.id = pt.uom_id
JOIN seller s on s.id = pt.seller_id
LEFT JOIN "user" us1 on us1.id = p.create_uid
LEFT JOIN "user" us2 on us2.id = p.write_uid
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN p.id = ANY(@ids::int8[]) ELSE TRUE END;

-- name: GetProducts :many
SELECT *
FROM product
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN id = ANY(@ids::int8[]) ELSE TRUE END;

-- name: GetProductAndRelations :many
SELECT p.*, c.id category_id, u.id uom_id, s.id seller_id
FROM product p
JOIN product_template pt on pt.id = p.template_id
JOIN category c on c.id = pt.category_id
JOIN uom u on u.id = pt.uom_id
JOIN seller s on s.id = pt.seller_id
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN p.id = ANY(@ids::int8[]) ELSE TRUE END;

-- name: GetProductTemplates :many
SELECT *
FROM product_template
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN id = ANY(@ids::int8[]) ELSE TRUE END;

-- name: GetCategories :many
SELECT *
FROM category
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN id = ANY(@ids::int8[]) ELSE TRUE END;

-- name: GetUoms :many
SELECT *
FROM uom
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN id = ANY(@ids::int8[]) ELSE TRUE END;

-- name: GetSellers :many
SELECT *
FROM seller
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN id = ANY(@ids::int8[]) ELSE TRUE END;

-- name: GetProductsByKeyword :many
SELECT id, COUNT(*) OVER() total
FROM product
WHERE "name" LIKE @keyword::varchar
ORDER BY id DESC
OFFSET @_offset::int8
LIMIT @_limit::int8;