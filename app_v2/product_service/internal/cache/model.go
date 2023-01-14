package cache

import (
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/database/store"
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
)

// Model
// ---------------------------------------------

type ModelValue interface {
	GetType() TypeCache
	GetId() int64
}

// User ...
type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (u User) GetId() int64 {
	return u.Id
}

func (u User) GetType() TypeCache {
	return TypeUser
}

// Product ...
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	OriginPrice float64   `json:"origin_price"`
	SalePrice   float64   `json:"sale_price"`
	State       string    `json:"state"`
	Variants    string    `json:"variants"`
	CreateDate  time.Time `json:"create_date"`
	WriteDate   time.Time `json:"write_date"`
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

func (p *Product) FromDb(model store.Product, categoryId, uomId, sellerId int64) {
	variants := string(model.Variants.RawMessage)
	*p = Product{
		ID:          model.ID,
		Name:        model.Name,
		OriginPrice: model.OriginPrice,
		SalePrice:   model.SalePrice,
		State:       model.State,
		Variants:    variants,
		CreateDate:  model.CreateDate,
		WriteDate:   model.WriteDate,
		CreateUid:   model.CreateUid.Int64,
		WriteUid:    model.WriteUid.Int64,
		TemplateID:  model.TemplateID.Int64,
		CategoryID:  categoryId,
		UomID:       uomId,
		SellerID:    sellerId,
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
		CreateUid:      model.CreateUid.Int64,
		WriteUid:       model.WriteUid.Int64,
		SellerID:       model.SellerID.Int64,
		CategoryID:     model.CategoryID.Int64,
		UomID:          model.UomID.Int64,
	}
	return
}

// Seller ...
type Seller struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	LogoUrl string `json:"logo_url"`
}

func (s Seller) GetId() int64 {
	return s.ID
}

func (s Seller) GetType() TypeCache {
	return TypeSeller
}

func (s *Seller) FromDb(model store.Seller) {
	*s = Seller{
		ID:      model.ID,
		Name:    model.Name,
		LogoUrl: model.LogoUrl.String,
	}
	return
}

// Uom ...
type Uom struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (u Uom) GetId() int64 {
	return u.ID
}

func (u Uom) GetType() TypeCache {
	return TypeUom
}

func (u *Uom) FromDb(model store.Uom) {
	*u = Uom{
		ID:   model.ID,
		Name: model.Name,
	}
	return
}

// Category ...
type Category struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c Category) GetId() int64 {
	return c.ID
}

func (c Category) GetType() TypeCache {
	return TypeCategory
}

func (c *Category) FromDb(model store.Category) {
	*c = Category{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description.String,
	}
	return
}
