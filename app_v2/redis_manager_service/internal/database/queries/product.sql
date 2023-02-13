-- name: GetProducts :many
SELECT *
FROM product
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN id = ANY(@ids::int8[]) ELSE TRUE END;

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

-- name: GetProductAndRelations :many
SELECT p.*, c.id category_id, u.id uom_id, s.id seller_id
FROM product p
         JOIN product_template pt on pt.id = p.template_id
         JOIN category c on c.id = pt.category_id
         JOIN uom u on u.id = pt.uom_id
         JOIN seller s on s.id = pt.seller_id
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN p.id = ANY(@ids::int8[]) ELSE TRUE END;