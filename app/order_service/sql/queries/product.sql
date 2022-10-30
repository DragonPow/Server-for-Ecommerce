-- name: getProducts :many
SELECT *
FROM product
WHERE
    CASE WHEN array_length(@ids::int8[], 1) > 0 THEN id = ANY(@ids::int8[]) ELSE TRUE END
AND CASE WHEN @product_template_id::int8 > 0 THEN template_id = @product_template_id::int8 ELSE TRUE END;