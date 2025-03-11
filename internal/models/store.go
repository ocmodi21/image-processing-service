package models

type Store struct {
	ID       string `json:"store_id"`
	Name     string `json:"store_name"`
	AreaCode string `json:"area_code"`
}
