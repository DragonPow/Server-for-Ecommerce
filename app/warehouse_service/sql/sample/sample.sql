INSERT INTO import_bill(id,code,note,contact_person_name,contact_email,contact_phone,create_uid,create_date,write_uid,write_date) VALUES
(1,'IMPORT_1',null,'Minh Long', null, null,1,now(),1,now()),
(2,'IMPORT_2',null,'Minh Long', null, null,1,now(),1,now());

INSERT INTO import_bill_detail(id,import_id,product_id,quantity) VALUES
(1,1,1,1),
(2,1,2,1),
(3,1,3,1),
(4,2,4,1),
(5,2,5,1),
(6,2,6,1),
(7,2,7,1),
(8,2,8,1),
(9,2,9,1);

-- Set integer to end
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"import_bill"', 'id')), (SELECT (MAX("id") + 1) FROM "import_bill"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"import_bill_detail"', 'id')), (SELECT (MAX("id") + 1) FROM "import_bill_detail"), FALSE);