-- name: GetUserByUserNameAndPassword :one
SELECT *
FROM account
WHERE username = @username::varchar
AND password = @password::varchar
LIMIT 1;