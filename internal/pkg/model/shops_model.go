package model

type MyShopResp struct {
	ID       uint   `json:"id"`
	ShopName string `json:"nama_toko"`
	PhotoURL string `json:"url_foto"`
	UserID   uint   `json:"user_id"`
}

type ShopResp struct {
	ID       uint   `json:"id"`
	ShopName string `json:"nama_toko"`
	PhotoURL string `json:"url_foto"`
}

type ShopInfo struct {
	ShopName string `json:"nama_toko"`
	PhotoURL string `json:"url_foto"`
}

type ShopReqUpdate struct {
	ShopName string `form:"nama_toko,omitempty"`
}

type ShopsFilter struct {
	ShopName string `query:"nama"`
	Limit    int    `query:"limit"`
	Page     int    `query:"page"`
}
