package model

type TrxDetailResp struct {
	ProductLog ProductLogResp `json:"product"`
	Shop       ShopInfo       `json:"toko"`
	Quantity   int            `json:"kuantitas"`
	TotalPrice int            `json:"harga_total"`
}

type TrxDetailReqCreate struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"kuantitas" validate:"required"`
}
