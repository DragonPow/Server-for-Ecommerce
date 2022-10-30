INSERT INTO account(id, username, password, create_date, write_date) VALUES
(1,'vungocthach','123456',now(),now()),
(2,'kimtulong','12345678',now(),now());

INSERT INTO customer_info(id, account_id, name, phone, address, create_date, write_date) VALUES
(1,1,'Vũ Ngọc Thạch','038225xxxx','KTX khu B, DHQG',now(),now()),
(2,2,'12345678',null,null,now(),now());

-- Set integer to end
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"account"', 'id')), (SELECT (MAX("id") + 1) FROM "account"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"customer_info"', 'id')), (SELECT (MAX("id") + 1) FROM "customer_info"), FALSE);