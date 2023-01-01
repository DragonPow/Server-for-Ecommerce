package api

type GetDetailProductRequest struct {
	Id int64 `json:"id"`
}

type GetDetailProductResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    ProductDetail `json:"data"`
}

type GetListProductRequest struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"page_size,omitempty"`
}

type GetListProductResponse struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Data    GetListProductResponseItem `json:"data"`
}

type GetListProductResponseItem struct {
	TotalItem int
	Item      []ProductOverview
}

type ProductDetail struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	OriginPrice float64 `json:"origin_price"`
	SalePrice   float64 `json:"sale_price"`
	Variants    string  `json:"variants"`
	CreatedBy   string  `json:"created_by"`
	CreatedDate string  `json:"created_date"`
	UpdatedBy   string  `json:"updated_by"`
	UpdatedDate string  `json:"updated_date"`

	TemplateId          int64   `json:"template_id"`
	TemplateName        string  `json:"template_name"`
	TemplateDescription string  `json:"template_description"`
	SoldQuantity        float64 `json:"sold_quantity"`
	RemainQuantity      float64 `json:"remain_quantity"`
	Rating              float64 `json:"rating"`
	NumberRating        int32   `json:"number_rating"`

	SellerId      int64  `json:"seller_id"`
	SellerName    string `json:"seller_name"`
	SellerLogo    string `json:"seller_logo"`
	SellerAddress string `json:"seller_address"`

	CategoryId   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`

	UomId   int64  `json:"uom_id"`
	UomName string `json:"uom_name"`
}

type ProductOverview struct {
	Id             int64   `json:"id"`
	Name           string  `json:"name"`
	OriginPrice    float64 `json:"origin_price"`
	SalePrice      float64 `json:"sale_price"`
	CreatedBy      string  `json:"created_by"`
	CreatedDate    string  `json:"created_date"`
	TemplateId     int64   `json:"template_id"`
	TemplateName   string  `json:"template_name"`
	SoldQuantity   float64 `json:"sold_quantity"`
	RemainQuantity float64 `json:"remain_quantity"`
	Rating         float64 `json:"rating"`

	SellerId   int64  `json:"seller_id"`
	SellerName string `json:"seller_name"`
	SellerLogo string `json:"seller_logo"`

	CategoryId   int64  `json:"category_id"`
	CategoryName string `json:"category_name"`

	UomId   int64  `json:"uom_id"`
	UomName string `json:"uom_name"`
}
