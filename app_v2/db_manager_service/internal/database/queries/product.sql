-- name: CreateUom :one
INSERT INTO seller
(name, description, phone, address, logo_url, manager_id, create_uid, create_date, write_uid, write_time) VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;