package api

import (
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/internal/database/store"
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
