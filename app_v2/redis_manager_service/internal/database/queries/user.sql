-- name: GetUsers :many
SELECT *
FROM "user"
WHERE CASE WHEN array_length(@ids::int8[], 1) > 0 THEN id = ANY(@ids::int8[]) ELSE TRUE END;