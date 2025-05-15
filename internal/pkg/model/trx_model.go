package model

type TrxResp struct {
	ID            uint            `json:"id"`
	TotalPrice    int             `json:"harga_total"`
	InvoiceCode   string          `json:"kode_invoice"`
	PaymentMethod string          `json:"method_bayar"`
	Address       AddressResp     `json:"alamat_kirim"`
	TrxDetail     []TrxDetailResp `json:"detail_trx"`
}

type TrxFilter struct {
	Limit  int    `query:"limit"`
	Page   int    `query:"page"`
	Search string `query:"search"`
}

type TrxReqCreate struct {
	PaymentMethod string               `json:"method_bayar" validate:"required"`
	AddressID     uint                 `json:"alamat_kirim" validate:"required"`
	TrxDetails    []TrxDetailReqCreate `json:"detail_trx" validate:"required"`
}
