package model

type Login struct {
	PhoneNumber string `json:"no_telp" validate:"required"`
	Password    string `json:"kata_sandi" validate:"required"`
}

type LoginRes struct {
	Name        string       `json:"nama"`
	PhoneNumber string       `json:"no_telp"`
	BirthDate   string       `json:"tanggal_lahir"`
	About       string       `json:"tentang"`
	JobTitle    string       `json:"pekerjaan"`
	Email       string       `json:"email"`
	ProvinceID  ProvinceResp `json:"id_provinsi"`
	CityID      CityResp     `json:"id_kota"`
	Token       string       `json:"token"`
}
