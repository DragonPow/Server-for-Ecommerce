-- name: getImportBills :many
SELECT *
FROM import_bill
WHERE id = ANY(@ids::int8[]);

-- name: getImportBillDetails :many
SELECT *
FROM import_bill_detail
WHERE
    CASE WHEN array_length(@ids::int8[],1) > 0 THEN id = ANY(@ids::int8[]) ELSE TRUE END
AND CASE WHEN array_length(@import_id::int8[],1) > 0 THEN import_id = ANY(@import_id::int8[]) ELSE TRUE END;