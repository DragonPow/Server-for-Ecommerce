package api

import (
	"encoding/json"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/cache"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/database/store"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/util"
)

func (p *ProductDetail) FromEntity(product store.GetProductDetailsRow) {
	*p = ProductDetail{
		Id:                  product.ID,
		Name:                product.Name,
		OriginPrice:         product.OriginPrice,
		SalePrice:           product.SalePrice,
		Variants:            string(product.Variants.RawMessage),
		CreatedBy:           product.CreateName.String,
		CreatedDate:         util.ParseTimeToString(product.CreateDate),
		UpdatedBy:           product.WriteName.String,
		UpdatedDate:         util.ParseTimeToString(product.WriteDate),
		TemplateId:          product.TemplateID.Int64,
		TemplateName:        product.TemplateName,
		TemplateDescription: product.TemplateDescription.String,
		SoldQuantity:        product.SoldQuantity,
		RemainQuantity:      product.RemainQuantity,
		Rating:              product.Rating,
		NumberRating:        int32(product.NumberRating),
		SellerId:            product.SellerID,
		SellerName:          product.SellerName,
		SellerLogo:          product.SellerLogo.String,
		SellerAddress:       product.SellerAddress.String,
		CategoryId:          product.CategoryID,
		CategoryName:        product.CategoryName,
		UomId:               product.UomID,
		UomName:             product.UomName,
	}
}

func (p *ProductOverview) FromCache(product cache.Product, template cache.ProductTemplate, category cache.Category, uom cache.Uom, seller cache.Seller) {
	b, _ := json.Marshal(product.Variants)
	*p = ProductOverview{
		Id:             product.ID,
		Name:           product.Name,
		OriginPrice:    product.OriginPrice,
		SalePrice:      product.SalePrice,
		Variants:       string(b),
		TemplateId:     template.ID,
		TemplateName:   template.Name,
		SoldQuantity:   template.SoldQuantity,
		RemainQuantity: template.RemainQuantity,
		Rating:         template.Rating,
		SellerId:       seller.ID,
		SellerName:     seller.Name,
		SellerLogo:     seller.Logo,
		CategoryId:     category.ID,
		CategoryName:   category.Name,
		UomId:          uom.ID,
		UomName:        uom.Name,
	}
}
