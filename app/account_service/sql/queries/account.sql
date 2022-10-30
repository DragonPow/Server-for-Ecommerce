-- name: GetUserByUserNameAndPassword :one
SELECT *
FROM account
WHERE username = @username::varchar
AND password = @password::varchar
LIMIT 1;

-- name: getCustomerByIds :many
SELECT *
FROM customer_info
WHERE id = ANY(@ids::int8[]);

-- name: getAccountByIds :many
SELECT id, username, create_date, write_date
FROM account
WHERE id = ANY(@ids::int8[]);