INSERT INTO product_template(id, name, description, default_price, remain_quantity, sold_quantity, rating,
                             number_rating, variants, seller_id, category_id, uom_id)
VALUES (1, 'Laptop Acer Nitro 5', null, 20000000, 2, 0, 4, 1, '{"color": "black", "weight":50, "uom_weight":"Cái"}', 1, 2, 1),
       (2, 'Chuột Gaming DareU', null, 1500000, 2, 0, 5, 5, '{"color": "black", "weight":50, "uom_weight":"Cái"}', 1, 2, 1),
       (3, 'Dụng cụ vệ sinh máy tính', null, 20000, 2, 0, 0, 0, '{"color": "black", "number_per_unit":3, "weight":50, "uom_weight":"Cái"}', 1, 2, 2),
       (4, 'Rau cải trắng', null, 15000, 0, 0, 4, 1, '{"weight":50, "uom_weight":"Gram"}', 1, 2, 1),
       (5, 'Cà rốt Mỹ', null, 1500000, 0, 0, 5, 5, '{"weight":50, "uom_weight":"Gram"}', 1, 2, 1),
       (6, 'Thuốc trị bách bệnh', null, 20000, 0, 0, 0, 0, '{"weight":50, "uom_weight":"Gói"}', 1, 2, 2);

INSERT INTO product(id, template_id, name, origin_price, sale_price, state,
                    variants, create_uid, write_uid)
VALUES (1, 1, 'AN5-M7882021', 18000000, 18000000, 'available', '{"warranty":"12", "warranty_time":"month", "note":null}', 1, 1),
       (2, 1, 'AN5-M7882022', 25000000, 25000000, 'sold', '{"warranty":"12", "warranty_time":"month", "note":null}', 1, 1),
       (3, 1, 'AN5-M7882023', 25000000, 25000000, 'available', '{"warranty":"24", "warranty_time":"month", "note":null}', 1, 1),
       (4, 2, '65M1-3327', 1500000, 1500000, 'available', '{"warranty":"6", "warranty_time":"month", "note":null}', 1, 1),
       (5, 2, '65M1-3328', 1500000, 1500000, 'available', '{"warranty":"7", "warranty_time":"month", "note":null}', 1, 1),
       (6, 3, '', 20000, 20000, 'wait_export', '{}', 1, 1),
       (7, 3, '', 20000, 20000, 'sold', '{}', 1, 1),
       (8, 3, '', 20000, 20000, 'sold', '{}', 1, 1),
       (9, 3, '', 20000, 20000, 'available', '{}', 1, 1);

SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"product_template"', 'id')),(SELECT (MAX("id") + 1) FROM "product_template"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"product"', 'id')), (SELECT (MAX("id") + 1) FROM "product"), FALSE);