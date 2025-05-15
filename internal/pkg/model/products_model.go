package model

type ProductResp struct {
	ID            uint               `json:"id"`
	ProductName   string             `json:"nama_produk"`
	Slug          string             `json:"slug"`
	ResellerPrice string             `json:"harga_reseler"`
	ConsumerPrice string             `json:"harga_konsumen"`
	Stock         int                `json:"stok"`
	Description   string             `json:"deskripsi"`
	Shop          ShopResp           `json:"toko"`
	Category      CategoryResp       `json:"category"`
	Images        []ProductImageResp `json:"photos"`
}

type ProductsFilter struct {
	Limit       int    `query:"limit"`
	Page        int    `query:"page"`
	ProductName string `query:"nama_produk"`
	CategoryID  uint   `query:"category_id"`
	ShopID      uint   `query:"toko_id"`
	MaxPrice    int    `query:"max_harga"`
	MinPrice    int    `query:"min_harga"`
}

type ProductReqCreate struct {
	ProductName   string `form:"nama_produk" validate:"required"`
	Slug          string `form:"slug,omitempty"`
	CategoryID    *uint  `form:"category_id,omitempty"`
	ResellerPrice string `form:"harga_reseller" validate:"required"`
	ConsumerPrice string `form:"harga_konsumen" validate:"required"`
	Stock         int    `form:"stok" validate:"required"`
	Description   string `form:"deskripsi" validate:"required"`
}

type ProductReqUpdate struct {
	ProductName   string `form:"nama_produk,omitempty"`
	Slug          string `form:"slug,omitempty"`
	CategoryID    *uint  `form:"category_id,omitempty"`
	ResellerPrice string `form:"harga_reseller,omitempty"`
	ConsumerPrice string `form:"harga_konsumen,omitempty"`
	Stock         int    `form:"stok,omitempty"`
	Description   string `form:"deskripsi,omitempty"`
}

type ProductImageResp struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	ImageURL  string `json:"url"`
}
