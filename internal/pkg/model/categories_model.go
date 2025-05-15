package model

type CategoryResp struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"nama_category"`
}

type CategoryReqCreate struct {
	CategoryName string `json:"nama_category" validate:"required"`
}

type CategoryReqUpdate struct {
	CategoryName string `json:"nama_category,omitempty"`
}
