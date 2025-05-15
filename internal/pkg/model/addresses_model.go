package model

type AddressReqCreate struct {
	AddressTitle  string `json:"judul_alamat" validate:"required"`
	RecipientName string `json:"nama_penerima" validate:"required"`
	PhoneNumber   string `json:"no_telp" validate:"required"`
	FullAddress   string `json:"detail_alamat" validate:"required"`
}

type AddressReqUpdate struct {
	AddressTitle  string `json:"judul_alamat,omitempty"`
	RecipientName string `json:"nama_penerima,omitempty"`
	PhoneNumber   string `json:"no_telp,omitempty"`
	FullAddress   string `json:"detail_alamat,omitempty"`
}

type AddressResp struct {
	ID            uint   `json:"id"`
	AddressTitle  string `json:"judul_alamat"`
	RecipientName string `json:"nama_penerima"`
	PhoneNumber   string `json:"no_telp"`
	FullAddress   string `json:"detail_alamat"`
}
