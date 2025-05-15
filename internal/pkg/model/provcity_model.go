package model

type ProvinceResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CityResp struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}
