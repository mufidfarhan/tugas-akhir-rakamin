package model

type UserReqCreate struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"kata_sandi" validate:"required"`
	Name        string `json:"nama" validate:"required"`
	PhoneNumber string `json:"no_telp" validate:"required"`
	BirthDate   string `json:"tanggal_lahir"`
	Gender      string `json:"jenis_kelamin"`
	About       string `json:"tentang"`
	JobTitle    string `json:"pekerjaan"`
	ProvinceID  string `json:"id_provinsi"`
	CityID      string `json:"id_kota"`
	IsAdmin     bool   `json:"is_admin"`
}

type UserReqUpdate struct {
	Name        string `json:"nama,omitempty"`
	Password    string `json:"kata_sandi,omitempty"`
	PhoneNumber string `json:"no_telp,omitempty"`
	BirthDate   string `json:"tanggal_Lahir,omitempty"`
	JobTitle    string `json:"pekerjaan,omitempty"`
	About       string `json:"tentang,omitempty"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	ProvinceID  string `json:"id_provinsi,omitempty"`
	CityID      string `json:"id_kota,omitempty"`
}

type UserResp struct {
	Name        string       `json:"nama"`
	PhoneNumber string       `json:"no_telp"`
	BirthDate   string       `json:"tanggal_lahir"`
	About       string       `json:"tentang"`
	JobTitle    string       `json:"pekerjaan"`
	Email       string       `json:"email"`
	ProvinceID  ProvinceResp `json:"id_provinsi"`
	CityID      CityResp     `json:"id_kota"`
}
