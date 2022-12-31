package api

type BaseResponse struct {
	Code    int
	Message string
}

type GetDetailProductRequest struct {
	Id int64 `json:"id"`
}

type GetDetailProductResponse struct {
	BaseResponse
	Data ProductDetail
}

type GetListProductRequest struct {
	Page     *int `json:"page, omitempty"`
	PageSize *int `json:"page_size, omitempty"`
}

type GetListProductResponse struct {
	BaseResponse
	Data GetListProductResponseItem
}

type GetListProductResponseItem struct {
	TotalItem int
	Item      []ProductOverview
}

type ProductDetail struct {
	Id          int64
	Name        string
	OriginPrice float64
	SalePrice   float64
	Variants    []byte
	CreatedBy   string
	CreatedDate string
	UpdatedBy   string
	UpdatedDate string
	Decription  string

	TemplateId          int64
	TemplateName        string
	TemplateDescription string
	SoldQuantity        float64
	RemainQuantity      float64
	Rating              float64
	NumberRating        float64

	SellerId      int64
	SellerName    string
	SellerLogo    string
	SellerAddress string

	CategoryId   int64
	CategoryName string

	UomId   int64
	UomName string
}

type ProductOverview struct {
	Id             int64
	Name           string
	OriginPrice    float64
	SalePrice      float64
	CreatedBy      string
	CreatedDate    string
	TemplateId     int64
	TemplateName   string
	SoldQuantity   float64
	RemainQuantity float64
	Rating         float64

	SellerId       int64
	SellerName     string
	SellerLogo     string

	CategoryId     int64
	CategoryName   string

	UomId          int64
	UomName        string
}
