package api

import (
	"Server-for-Ecommerce/app_v2/product_service/cache"
	"Server-for-Ecommerce/app_v2/product_service/database/store"
	"Server-for-Ecommerce/app_v2/product_service/util"
)

func (p *ProductDetail) FromEntity(product store.GetProductDetailsRow) {
	*p = ProductDetail{
		Id:                  product.ID,
		Image:               product.Image,
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

func (p *ProductOverview) FromCache(product cache.Product) {
	*p = ProductOverview{
		Id:          product.ID,
		Name:        product.Name,
		OriginPrice: product.OriginPrice,
		SalePrice:   product.SalePrice,
		Image:       product.Image,
	}
}
