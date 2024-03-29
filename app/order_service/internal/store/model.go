// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package store

import (
	"database/sql"
)

type OrderBill struct {
	ID             int64           `json:"id"`
	CustomerID     int64           `json:"customer_id"`
	PaymentMethod  string          `json:"payment_method"`
	ContactName    sql.NullString  `json:"contact_name"`
	ContactPhone   sql.NullString  `json:"contact_phone"`
	ContactAddress sql.NullString  `json:"contact_address"`
	TotalPrice     sql.NullFloat64 `json:"total_price"`
	ShipCost       sql.NullFloat64 `json:"ship_cost"`
	State          string          `json:"state"`
	Note           sql.NullString  `json:"note"`
	CreateUid      sql.NullInt64   `json:"create_uid"`
	CreateDate     sql.NullTime    `json:"create_date"`
	WriteUid       sql.NullInt64   `json:"write_uid"`
	WriteDate      sql.NullTime    `json:"write_date"`
}

type OrderBillDetail struct {
	ID                int64           `json:"id"`
	OrderID           int64           `json:"order_id"`
	ProductTemplateID int64           `json:"product_template_id"`
	Quantity          sql.NullFloat64 `json:"quantity"`
	UnitPrice         sql.NullFloat64 `json:"unit_price"`
	TotalPrice        sql.NullFloat64 `json:"total_price"`
}

type OrderShipping struct {
	ID              int64          `json:"id"`
	OrderID         int64          `json:"order_id"`
	State           string         `json:"state"`
	ShippingName    sql.NullString `json:"shipping_name"`
	ShippingPhone   sql.NullString `json:"shipping_phone"`
	ShippingAddress sql.NullString `json:"shipping_address"`
	CreateUid       sql.NullInt64  `json:"create_uid"`
	CreateDate      sql.NullTime   `json:"create_date"`
	WriteUid        sql.NullInt64  `json:"write_uid"`
	WriteDate       sql.NullTime   `json:"write_date"`
}

type OrderShippingDetail struct {
	ID            int64           `json:"id"`
	ShippingID    int64           `json:"shipping_id"`
	OrderDetailID int64           `json:"order_detail_id"`
	ProductID     int64           `json:"product_id"`
	Quantity      sql.NullFloat64 `json:"quantity"`
}

type Product struct {
	ID          int64           `json:"id"`
	TemplateID  int64           `json:"template_id"`
	Name        string          `json:"name"`
	OriginPrice sql.NullFloat64 `json:"origin_price"`
	SalePrice   sql.NullFloat64 `json:"sale_price"`
	State       string          `json:"state"`
	CreateUid   sql.NullInt64   `json:"create_uid"`
	CreateDate  sql.NullTime    `json:"create_date"`
	WriteUid    sql.NullInt64   `json:"write_uid"`
	WriteDate   sql.NullTime    `json:"write_date"`
}

type ProductTemplate struct {
	ID                int64           `json:"id"`
	Name              string          `json:"name"`
	DefaultPrice      sql.NullFloat64 `json:"default_price"`
	UomName           string          `json:"uom_name"`
	InventoryQuantity sql.NullFloat64 `json:"inventory_quantity"`
	CreateUid         sql.NullInt64   `json:"create_uid"`
	CreateDate        sql.NullTime    `json:"create_date"`
	WriteUid          sql.NullInt64   `json:"write_uid"`
	WriteDate         sql.NullTime    `json:"write_date"`
}
