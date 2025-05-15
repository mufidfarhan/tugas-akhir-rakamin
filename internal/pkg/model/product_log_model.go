package model

type ProductLogResp struct {
	ID            uint               `json:"id"`
	ProductName   string             `json:"nama_produk"`
	Slug          string             `json:"slug"`
	ResellerPrice string             `json:"harga_reseler"`
	ConsumerPrice string             `json:"harga_konsumen"`
	Description   string             `json:"deskripsi"`
	Shop          ShopInfo           `json:"toko"`
	Category      CategoryResp       `json:"category"`
	Images        []ProductImageResp `json:"photos"`
}
