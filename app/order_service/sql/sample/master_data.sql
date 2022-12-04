INSERT INTO category(id, "name") VALUES
(1, 'Dụng cụ nấu ăn'),
(2, 'Thiết bị công nghệ'),
(3, 'Dụng cụ văn phòng'),
(4, 'Dụng cụ trang trí nhà cửa'),
(5, 'Dụng cụ y tế');

INSERT INTO seller(id, "name", description, phone, address, logo_url, manager_id, create_uid, write_uid) VALUES
(1, 'Phong Vũ', '', '038225xxxx', 'Quận 1, thành phố Hồ Chí Minh', null, 1, 1, 1),
(2, 'Bách hóa xanh', '', '099245xxxx', 'Gò Vấp, thành phố Hồ Chí Minh', null, 2, 2, 2),
(3, 'Nhà thuốc Long Châu', '', '099956xxxx', 'Trung tâm thương mại Gigamall, Quận 1, thành phố Hồ Chí Minh', null, 1, 1, 1);

INSERT INTO uom(id, name, seller_id) VALUES
(1, 'Cái', 1),
(2, 'Bộ', 1),
(3, 'Bịch', 2),
(4, 'Chai', 2),
(5, 'Cái', 2),
(6, 'Bộ', 2),
(7, 'Tờ', 2),
(8, 'Bịch', 3),
(9, 'Viên', 3),
(10, 'Hũ', 3),
(11, 'Hôp', 3);

-- Set integer to end
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"category"', 'id')), (SELECT (MAX("id") + 1) FROM "category"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"seller"', 'id')), (SELECT (MAX("id") + 1) FROM "seller"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"uom"', 'id')), (SELECT (MAX("id") + 1) FROM "uom"), FALSE);