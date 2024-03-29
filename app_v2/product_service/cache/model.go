package cache

import (
	"Server-for-Ecommerce/app_v2/product_service/database/store"
	"encoding/json"
	"time"
)

// Common function
// -----------------------------------------------

type TypeCache string

const (
	TypeProduct         TypeCache = "product"
	TypeUser            TypeCache = "user"
	TypeCategory        TypeCache = "category"
	TypeProductTemplate TypeCache = "product_template"
	TypeSeller          TypeCache = "seller"
	TypeUom             TypeCache = "uom"
	TypePageProduct     TypeCache = "page_product"
)

// Model
// ---------------------------------------------

type ModelValue interface {
	GetType() TypeCache
	GetId() int64
	GetVersion() string
}

// User ...
type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	UpdateDate time.Time `json:"update_date"`
}

func (u User) GetId() int64 {
	return u.ID
}

func (u User) GetType() TypeCache {
	return TypeUser
}

func (u User) GetVersion() string {
	return u.UpdateDate.Format(time.RFC3339)
}

func (u *User) FromDb(model store.User) {
	*u = User{
		ID:         model.ID,
		Name:       model.Name,
		UpdateDate: model.WriteDate,
	}
	return
}

// Product ...
type Product struct {
	ID          int64          `json:"id"`
	Image       string         `json:"image"`
	Name        string         `json:"name"`
	OriginPrice float64        `json:"origin_price"`
	SalePrice   float64        `json:"sale_price"`
	State       string         `json:"state"`
	Variants    map[string]any `json:"variants"`
	CreateDate  time.Time      `json:"create_date"`
	WriteDate   time.Time      `json:"write_date"`
	// Reference
	CreateUid  int64 `json:"create_uid"`
	WriteUid   int64 `json:"write_uid"`
	TemplateID int64 `json:"template_id"`
	CategoryID int64 `json:"category_id"`
	UomID      int64 `json:"uom_id"`
	SellerID   int64 `json:"seller_id"`
}

func (p Product) GetId() int64 {
	return p.ID
}

func (p Product) GetType() TypeCache {
	return TypeProduct
}

func (p Product) GetVersion() string {
	return p.WriteDate.Format(time.RFC3339)
}

func (p *Product) UpdateVersion(version string) error {
	t, err := time.Parse(time.RFC3339, version)
	if err != nil {
		return err
	}
	p.WriteDate = t
	return nil
}

func (p *Product) FromDb(model store.Product, categoryId, uomId, sellerId int64) {
	var variants map[string]any
	json.Unmarshal(model.Variants.RawMessage, &variants)
	*p = Product{
		ID:          model.ID,
		Image:       model.Image,
		Name:        model.Name,
		OriginPrice: model.OriginPrice,
		SalePrice:   model.SalePrice,
		State:       model.State,
		Variants:    variants,
		CreateDate:  model.CreateDate,
		WriteDate:   model.WriteDate,
		CreateUid:   model.CreateUid,
		WriteUid:    model.WriteUid,
		TemplateID:  model.TemplateID.Int64,
		CategoryID:  categoryId,
		UomID:       uomId,
		SellerID:    sellerId,
	}
	return
}

func (p *Product) FromDbV2(model store.GetProductAndRelationsRow) {
	var variants map[string]any
	json.Unmarshal(model.Variants.RawMessage, &variants)
	*p = Product{
		ID:          model.ID,
		Image:       model.Image,
		Name:        model.Name,
		OriginPrice: model.OriginPrice,
		SalePrice:   model.SalePrice,
		State:       model.State,
		Variants:    variants,
		CreateDate:  model.CreateDate,
		WriteDate:   model.WriteDate,
		CreateUid:   model.CreateUid,
		WriteUid:    model.WriteUid,
		TemplateID:  model.TemplateID.Int64,
		CategoryID:  model.CategoryID,
		UomID:       model.UomID,
		SellerID:    model.SellerID,
	}
	return
}

// ProductTemplate ...
type ProductTemplate struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	DefaultPrice   float64   `json:"default_price"`
	RemainQuantity float64   `json:"remain_quantity"`
	SoldQuantity   float64   `json:"sold_quantity"`
	Rating         float64   `json:"rating"`
	NumberRating   int64     `json:"number_rating"`
	CreateDate     time.Time `json:"create_date"`
	WriteDate      time.Time `json:"write_date"`
	Variants       string    `json:"variants"`
	// Reference
	CreateUid  int64 `json:"create_uid"`
	WriteUid   int64 `json:"write_uid"`
	SellerID   int64 `json:"seller_id"`
	CategoryID int64 `json:"category_id"`
	UomID      int64 `json:"uom_id"`
}

func (p ProductTemplate) GetId() int64 {
	return p.ID
}

func (p ProductTemplate) GetType() TypeCache {
	return TypeProductTemplate
}

func (p ProductTemplate) GetVersion() string {
	return p.WriteDate.Format(time.RFC3339)
}

func (p *ProductTemplate) FromDb(model store.ProductTemplate) {
	variants := string(model.Variants.RawMessage)
	*p = ProductTemplate{
		ID:             model.ID,
		Name:           model.Name,
		Description:    model.Description.String,
		DefaultPrice:   model.DefaultPrice,
		RemainQuantity: model.RemainQuantity,
		SoldQuantity:   model.SoldQuantity,
		Rating:         model.Rating,
		NumberRating:   model.NumberRating,
		CreateDate:     model.CreateDate,
		WriteDate:      model.WriteDate,
		Variants:       variants,
		CreateUid:      model.CreateUid,
		WriteUid:       model.WriteUid,
		SellerID:       model.SellerID.Int64,
		CategoryID:     model.CategoryID.Int64,
		UomID:          model.UomID.Int64,
	}
	return
}

// Seller ...
type Seller struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Logo      string    `json:"logo"`
	Address   string    `json:"address"`
	WriteDate time.Time `json:"write_date"`
}

func (s Seller) GetId() int64 {
	return s.ID
}

func (s Seller) GetType() TypeCache {
	return TypeSeller
}

func (s Seller) GetVersion() string {
	return s.WriteDate.Format(time.RFC3339)
}

func (s *Seller) FromDb(model store.Seller) {
	*s = Seller{
		ID:        model.ID,
		Name:      model.Name,
		Logo:      model.LogoUrl.String,
		Address:   model.Address.String,
		WriteDate: model.WriteDate,
	}
	return
}

// Uom ...
type Uom struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	WriteDate time.Time `json:"write_date"`
}

func (u Uom) GetId() int64 {
	return u.ID
}

func (u Uom) GetType() TypeCache {
	return TypeUom
}

func (u Uom) GetVersion() string {
	return u.WriteDate.Format(time.RFC3339)
}

func (u *Uom) FromDb(model store.Uom) {
	*u = Uom{
		ID:        model.ID,
		Name:      model.Name,
		WriteDate: model.WriteDate,
	}
	return
}

// Category ...
type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WriteDate   time.Time `json:"write_date"`
}

func (c Category) GetId() int64 {
	return c.ID
}

func (c Category) GetType() TypeCache {
	return TypeCategory
}

func (c Category) GetVersion() string {
	return c.WriteDate.Format(time.RFC3339)
}

func (c *Category) FromDb(model store.Category) {
	*c = Category{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description.String,
		WriteDate:   model.WriteDate,
	}
	return
}
