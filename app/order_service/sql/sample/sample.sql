INSERT INTO product_template(id,name,default_price,uom_name,inventory_quantity,create_uid,create_date,write_uid,write_date) VALUES
(1,'Laptop Acer Nitro 5',20000000,'Cái',2,1,now(),1,now()),
(2,'Chuột Gaming DareU',1500000,'Cái',2,1,now(),1,now()),
(3,'Dụng cụ vệ sinh máy tính',20000,'Bộ',2,1,now(),1,now());

INSERT INTO product(id,template_id,name,origin_price,sale_price,state,create_uid,create_date,write_uid,write_date) VALUES
(1,1,'AN5-M7882021',18000000,18000000,'available',1,now(),1,now()),
(2,1,'AN5-M7882022',25000000,25000000,'sold',1,now(),1,now()),
(3,1,'AN5-M7882023',25000000,25000000,'available',1,now(),1,now()),
(4,2,'65M1-3327',1500000,1500000,'available',1,now(),1,now()),
(5,2,'65M1-3328',1500000,1500000,'available',1,now(),1,now()),
(6,3,'',20000,20000,'wait_export',1,now(),1,now()),
(7,3,'',20000,20000,'sold',1,now(),1,now()),
(8,3,'',20000,20000,'sold',1,now(),1,now()),
(9,3,'',20000,20000,'available',1,now(),1,now());

INSERT INTO order_bill(id,customer_id,payment_method,contact_name,contact_phone,contact_address,total_price,ship_cost,state,note,create_uid,create_date,write_uid,write_date) VALUES
(1,2,'BIDV','Long','038225xxxx','KTX Khu B ĐHQG',20000,0,'done',null,2,now(),1,now()),
(2,2,'CASH','Minh Long',null,'KTX Khu A ĐHQG',25000000,0,'done',null,2,now(),1,now()),
(3,2,'CASH','Minh Long',null,'UIT',20000,10000,'wait_confirm',null,2,now(),2,now()),
(4,2,'CASH','Minh Long',null,'UIT',20000,10000,'cancel',null,2,now(),2,now());

INSERT INTO order_bill_detail(id,order_id,product_template_id,quantity,unit_price,total_price) VALUES
(1,1,3,1,20000,20000),
(2,2,1,1,25000000,25000000),
(3,2,3,1,0,0),
(4,3,3,1,20000,20000),
(5,4,3,1,20000,20000);

INSERT INTO order_shipping(id,order_id,state,shipping_name,shipping_phone,shipping_address,create_uid,create_date,write_uid,write_date) VALUES
(1,1,'done','Thạch',null,null,1,now(),1,now()),
(2,2,'done','Ngọc Thạch','038225xxxx','UIT',1,now(),1,now());

INSERT INTO order_shipping_detail(id,shipping_id,order_detail_id,product_id,quantity) VALUES
(1,1,1,7,1),
(2,2,2,2,1),
(3,2,3,8,1);

-- Set integer to end
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"product_template"', 'id')), (SELECT (MAX("id") + 1) FROM "product_template"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"product"', 'id')), (SELECT (MAX("id") + 1) FROM "product"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"order_bill"', 'id')), (SELECT (MAX("id") + 1) FROM "order_bill"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"order_bill_detail"', 'id')), (SELECT (MAX("id") + 1) FROM "order_bill_detail"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"order_shipping"', 'id')), (SELECT (MAX("id") + 1) FROM "order_shipping"), FALSE);
SELECT SETVAL((SELECT PG_GET_SERIAL_SEQUENCE('"order_shipping_detail"', 'id')), (SELECT (MAX("id") + 1) FROM "order_shipping_detail"), FALSE);